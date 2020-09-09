package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/gerbenjacobs/millwheat/game"
	gamedata "github.com/gerbenjacobs/millwheat/game/data"
)

type GameSvc struct {
	townSvc   TownService
	prodSvc   ProductionService
	battleSvc BattleService

	// game data
	Items     game.Items
	Buildings game.Buildings
}

func NewGameSvc(townSvc TownService, prodSvc ProductionService, battleSvc BattleService, items game.Items, buildings game.Buildings) *GameSvc {
	return &GameSvc{
		townSvc:   townSvc,
		prodSvc:   prodSvc,
		battleSvc: battleSvc,
		Items:     items,
		Buildings: buildings}
}

func (g *GameSvc) Produce(ctx context.Context, buildingID uuid.UUID, set game.ItemSet) error {
	// get building
	townBuilding, building, err := g.getBuilding(ctx, buildingID)
	if err != nil {
		return err
	}

	// check and produce item
	if !building.CanDealWith(set.ItemID) {
		return errors.New("building can't handle this product")
	}
	productionResult, err := building.CreateProduct(set.ItemID, set.Quantity, townBuilding.CurrentLevel)
	if err != nil {
		return err
	}

	// extract consumption items from warehouse
	if err := g.townSvc.TakeFromWarehouse(ctx, productionResult.Consumption); err != nil {
		return err
	}

	// queue job
	job := &game.InputJob{
		Type: game.JobTypeProduct,
		ProductJob: &game.ProductJob{
			BuildingID:  townBuilding.ID,
			Production:  productionResult.Production,
			Consumption: productionResult.Consumption,
		},
		Duration: time.Duration(productionResult.Hours) * time.Hour,
	}
	if err := g.prodSvc.CreateJob(ctx, job); err != nil {
		// return back items
		_ = g.townSvc.GiveToWarehouse(ctx, productionResult.Consumption)
		return err
	}

	logrus.
		WithField("town", TownFromContext(ctx)).
		Debugf("producing %s in %s", job.ProductJob.Production, building.Name)
	return nil
}

func (g *GameSvc) Collect(ctx context.Context, buildingID uuid.UUID) error {
	// get building
	townBuilding, building, err := g.getBuilding(ctx, buildingID)
	if err != nil {
		return err
	}

	// take and store production
	cp, err := townBuilding.GetCurrentProduction(*building)
	if err != nil {
		return err
	}
	if err = g.townSvc.GiveToWarehouse(ctx, []game.ItemSet{*cp}); err != nil {
		return err
	}

	logrus.
		WithField("town", TownFromContext(ctx)).
		Debugf("collecting %s in %s", cp, building.Name)

	// update database
	return g.townSvc.BuildingCollected(ctx, buildingID)
}

func (g *GameSvc) AddBuilding(ctx context.Context, buildingType game.BuildingType) error {
	return g.upgradeBuilding(ctx, nil, buildingType, 1)
}

func (g *GameSvc) UpgradeBuilding(ctx context.Context, buildingID uuid.UUID) error {
	// get building
	townBuilding, _, err := g.getBuilding(ctx, buildingID)
	if err != nil {
		return err
	}
	return g.upgradeBuilding(ctx, &buildingID, townBuilding.Type, townBuilding.CurrentLevel+1)
}

func (g *GameSvc) DemolishBuilding(ctx context.Context, buildingID uuid.UUID) error {
	townBuilding, _, err := g.getBuilding(ctx, buildingID)
	if err != nil {
		return err
	}

	// demolish building
	if err = g.townSvc.RemoveBuilding(ctx, buildingID); err != nil {
		return err
	}

	// give recovered resources to warehouse
	b, ok := gamedata.Buildings[townBuilding.Type]
	pr, err := game.RecoverBuilding(b, townBuilding.CurrentLevel)
	if !ok || err != nil {
		return errors.New("failed to recover building materials")
	}

	if err := g.townSvc.GiveToWarehouse(ctx, pr.Consumption); err != nil {
		return err
	}

	logrus.
		WithField("town", TownFromContext(ctx)).
		Debugf("demolished a %s", b.Name)

	return nil
}

func (g *GameSvc) CancelJob(ctx context.Context, jobID uuid.UUID) error {
	// Collect returnable resources
	resources, err := g.prodSvc.RevertJobResources(ctx, jobID)
	if err != nil {
		return err
	}

	// Cancel job && reshuffle
	if err = g.prodSvc.CancelJob(ctx, jobID); err != nil {
		return err
	}
	g.prodSvc.ReshuffleQueue(ctx)

	// Apply resources to warehouse
	if err = g.townSvc.GiveToWarehouse(ctx, resources); err != nil {
		return err
	}

	logrus.
		WithField("town", TownFromContext(ctx)).
		WithField("jobID", jobID).
		Debugf("canceled job, returning %s", resources)

	return nil
}

func (g *GameSvc) CreateWarriors(ctx context.Context, warriorType game.WarriorType, quantity int) error {
	costs, err := game.CalculateWarriorCosts(warriorType, quantity)
	if err != nil {
		return err
	}

	// extract consumption items from warehouse
	if err := g.townSvc.TakeFromWarehouse(ctx, costs); err != nil {
		return err
	}

	if err := g.battleSvc.AddWarrior(ctx, TMPCurrentBattleId, TMPArmyId, TownFromContext(ctx), warriorType, quantity); err != nil {
		// return back items
		_ = g.townSvc.GiveToWarehouse(ctx, costs)
	}

	return err
}

func (g *GameSvc) getBuilding(ctx context.Context, buildingID uuid.UUID) (*game.TownBuilding, *game.Building, error) {
	// get town
	town, err := g.townSvc.Town(ctx, TownFromContext(ctx))
	if err != nil {
		return nil, nil, err
	}

	// get buildings
	var townBuilding *game.TownBuilding
	var building *game.Building
	for _, tb := range town.Buildings {
		if buildingID == tb.ID {
			townB := tb
			b := g.Buildings[tb.Type]

			townBuilding = &townB
			building = &b
			break
		}
	}
	if townBuilding == nil || building == nil {
		return nil, nil, errors.New("building not found")
	}

	return townBuilding, building, nil
}

func (g *GameSvc) upgradeBuilding(ctx context.Context, buildingID *uuid.UUID, buildingType game.BuildingType, level int) error {
	building, ok := gamedata.Buildings[buildingType]
	if !ok {
		return errors.New("building type doesn't exist")
	}
	// get production requirements for building
	productionResult, err := game.CreateBuilding(building, level)
	if err != nil {
		return err
	}
	// check if items are in warehouse
	if !g.townSvc.ItemsInWarehouse(ctx, productionResult.Consumption) {
		// how to use this for messaging?
		//_ = handler.storeAndSaveFlash(r, w, "error|You don't have the required products; "+productionResult.Consumption.String())
		return errors.New("missing required items")
	}

	// extract consumption items from warehouse
	if err := g.townSvc.TakeFromWarehouse(ctx, productionResult.Consumption); err != nil {
		return err
	}

	// set or create building id
	bID := uuid.New()
	if buildingID != nil {
		bID = *buildingID
	}

	// queue building job
	if err := g.prodSvc.CreateJob(ctx, &game.InputJob{
		Type: game.JobTypeBuilding,
		BuildingJob: &game.BuildingJob{
			ID:    bID,
			Type:  buildingType,
			Level: level,
		},
		Duration: 3600 * time.Second, // TODO: fix time
	}); err != nil {
		// return back items
		_ = g.townSvc.GiveToWarehouse(ctx, productionResult.Consumption)
		return err
	}

	logrus.
		WithField("town", TownFromContext(ctx)).
		Debugf("upgrading %s to level %d", building.Name, level)
	return nil
}
