package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
)

type GameSvc struct {
	townSvc TownService
	prodSvc ProductionService

	// game data
	Items     game.Items
	Buildings game.Buildings
}

func NewGameSvc(townSvc TownService, prodSvc ProductionService, items game.Items, buildings game.Buildings) *GameSvc {
	return &GameSvc{townSvc: townSvc, prodSvc: prodSvc, Items: items, Buildings: buildings}
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
		Duration: time.Duration(productionResult.Hours) * time.Minute,
	}
	if err := g.prodSvc.CreateJob(ctx, job); err != nil {
		// return back items
		_ = g.townSvc.GiveToWarehouse(ctx, productionResult.Consumption)
		return err
	}

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

	// update database
	return g.townSvc.BuildingCollected(ctx, buildingID)
}

func (g *GameSvc) AddBuilding(ctx context.Context, buildingType game.BuildingType) error {
	panic("implement me")
}

func (g *GameSvc) UpgradeBuilding(ctx context.Context, buildingType game.BuildingType) error {
	panic("implement me")
}

func (g *GameSvc) DemolishBuilding(ctx context.Context, buildingID uuid.UUID) error {
	panic("implement me")
}

func (g *GameSvc) CancelJob(ctx context.Context, jobID uuid.UUID) error {
	panic("implement me")
}
