package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
)

type warehouseDTO map[string]int

// defaultWarehouse returns a new map, to prevent pointer issues
func defaultWarehouse() map[game.ItemID]game.WarehouseItem {
	return map[game.ItemID]game.WarehouseItem{
		"stone":    {ItemID: "stone", Quantity: 100},
		"plank":    {ItemID: "plank", Quantity: 100},
		"log":      {ItemID: "plank", Quantity: 10},
		"wheat":    {ItemID: "wheat", Quantity: 10},
		"flour":    {ItemID: "flour", Quantity: 4},
		"iron_bar": {ItemID: "iron_bar", Quantity: 4},
	}
}

func whToDTO(wh map[game.ItemID]game.WarehouseItem) warehouseDTO {
	dto := make(warehouseDTO)
	for _, wi := range wh {
		dto[wi.ItemID.AsKey()] = wi.Quantity
	}
	return dto
}

func dtoToWH(dto warehouseDTO) map[game.ItemID]game.WarehouseItem {
	wh := make(map[game.ItemID]game.WarehouseItem)
	for k, v := range dto {
		wh[game.ItemID(k)] = game.WarehouseItem{
			ItemID:   game.ItemID(k),
			Quantity: v,
		}
	}
	return wh
}

func (t *TownRepository) createTownInDatabase(ctx context.Context, town *game.Town) error {
	tid, _ := town.ID.MarshalBinary()
	oid, _ := town.Owner.MarshalBinary()
	whBytes, err := json.Marshal(whToDTO(town.Warehouse))
	if err != nil {
		return err
	}

	stmt, err := t.db.PrepareContext(ctx, "INSERT INTO towns (id, owner, name, warehouse, createdAt, updatedAt) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, tid, oid, town.Name, whBytes, town.CreatedAt, town.UpdatedAt)
	return err

}

func (t *TownRepository) getTownFromDatabase(ctx context.Context, id uuid.UUID) (*game.Town, error) {
	tid, _ := id.MarshalBinary()
	row := t.db.QueryRowContext(ctx, "SELECT id, owner, name, warehouse, createdAt, updatedAt FROM towns WHERE id = ?", tid)

	town := new(game.Town)
	var whBytes []byte
	err := row.Scan(&town.ID, &town.Owner, &town.Name, &whBytes, &town.CreatedAt, &town.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("town with ID %q not found", id)
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning town: %v", err)
	}

	var whDTO warehouseDTO
	if err = json.Unmarshal(whBytes, &whDTO); err != nil {
		return nil, err
	}
	buildings := make(map[uuid.UUID]game.TownBuilding)
	b, err := t.getBuildingsFromDatabase(ctx, id)
	if err == nil {
		buildings = b
	}
	town.Buildings = buildings
	town.Warehouse = dtoToWH(whDTO)

	t.townCache.Set(id.String(), town, CacheDurationTown)
	return town, nil
}

func (t *TownRepository) getBuildingsFromDatabase(ctx context.Context, townID uuid.UUID) (map[uuid.UUID]game.TownBuilding, error) {
	tid, _ := townID.MarshalBinary()
	rows, err := t.db.QueryContext(ctx, "SELECT id, type, level, lastCollection, createdAt FROM buildings WHERE townId = ?", tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	buildings := make(map[uuid.UUID]game.TownBuilding)
	for rows.Next() {
		var tb game.TownBuilding
		err = rows.Scan(&tb.ID, &tb.Type, &tb.CurrentLevel, &tb.LastCollection, &tb.CreatedAt)
		if err != nil {
			return nil, err
		}
		buildings[tb.ID] = tb
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return buildings, nil
}

func (t *TownRepository) updateWarehouseInDatabase(ctx context.Context, townID uuid.UUID, wh map[game.ItemID]game.WarehouseItem) error {
	tid, _ := townID.MarshalBinary()
	whBytes, err := json.Marshal(whToDTO(wh))
	if err != nil {
		return err
	}

	query := "UPDATE towns SET warehouse = ?, updatedAt = ? WHERE id = ?"
	_, err = t.db.ExecContext(ctx, query, whBytes, time.Now().UTC(), tid)
	if err != nil {
		return err
	}

	// update town struct and save to cache
	town, _ := t.Get(ctx, townID)
	town.UpdatedAt = time.Now().UTC()
	town.Warehouse = wh
	t.townCache.Set(townID.String(), town, CacheDurationTown)
	return nil
}

func (t *TownRepository) addBuildingToDatabase(ctx context.Context, townID uuid.UUID, building game.TownBuilding) error {
	bid, _ := building.ID.MarshalBinary()
	tid, _ := townID.MarshalBinary()

	// write building to database
	stmt, err := t.db.PrepareContext(ctx, "INSERT INTO buildings (id, townId, type, level, lastCollection, createdAt) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, bid, tid, building.Type, building.CurrentLevel, building.LastCollection, building.CreatedAt)

	if err != nil {
		return err
	}

	// update town struct and save to cache
	town, _ := t.Get(ctx, townID)
	town.Buildings[building.ID] = building
	t.townCache.Set(townID.String(), town, CacheDurationTown)

	return nil
}

func (t *TownRepository) upgradeBuildingInDatabase(ctx context.Context, townID uuid.UUID, building game.TownBuilding) error {
	bid, _ := building.ID.MarshalBinary()

	query := "UPDATE buildings SET level = ?  WHERE id = ?"
	_, err := t.db.ExecContext(ctx, query, building.CurrentLevel, bid)
	if err != nil {
		return err
	}

	// update town struct and save to cache
	town, _ := t.Get(ctx, townID)
	town.Buildings[building.ID] = building
	t.townCache.Set(townID.String(), town, CacheDurationTown)

	return nil
}

func (t *TownRepository) removeBuildingInDatabase(ctx context.Context, townID uuid.UUID, buildingID uuid.UUID) error {
	bid, _ := buildingID.MarshalBinary()

	query := "DELETE FROM buildings WHERE id = ?"
	_, err := t.db.ExecContext(ctx, query, bid)
	if err != nil {
		return err
	}

	// update town struct and save to cache
	town, _ := t.Get(ctx, townID)
	delete(town.Buildings, buildingID)
	t.townCache.Set(townID.String(), town, CacheDurationTown)

	return nil
}

func (t *TownRepository) updateBuildingCollection(ctx context.Context, townID uuid.UUID, building game.TownBuilding) error {
	bid, _ := building.ID.MarshalBinary()

	query := "UPDATE buildings SET lastCollection = ?  WHERE id = ?"
	_, err := t.db.ExecContext(ctx, query, building.LastCollection, bid)
	if err != nil {
		return err
	}

	// update town struct and save to cache
	town, _ := t.Get(ctx, townID)
	town.Buildings[building.ID] = building
	t.townCache.Set(townID.String(), town, CacheDurationTown)

	return nil
}
