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

func (t *TownRepository) getTownFromDatabase(ctx context.Context, id uuid.UUID) (*game.Town, error) {
	tid, _ := id.MarshalBinary()
	row := t.db.QueryRowContext(ctx, "SELECT id, owner, name, warehouse, createdAt, updatedAt FROM towns WHERE id = ?", tid)

	town := new(game.Town)
	warehouse := make(map[game.ItemID]game.WarehouseItem)
	var whBytes []byte
	err := row.Scan(&town.ID, &town.Owner, &town.Name, &whBytes, &town.CreatedAt, &town.UpdatedAt)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("town with ID %q not found", id)
	case err != nil:
		return nil, fmt.Errorf("unknown error while scanning town: %v", err)
	}

	if err = json.Unmarshal(whBytes, &warehouse); err != nil {
		return nil, err
	}
	buildings := make(map[uuid.UUID]game.TownBuilding)
	b, err := t.getBuildingsFromDatabase(ctx, id)
	if err == nil {
		buildings = b
	}
	town.Buildings = buildings

	t.townCache.Set(id.String(), town, CacheDurationTown)
	t.warehouseCache.Set(id.String(), warehouse, CacheDurationWarehouse)
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
	whBytes, err := json.Marshal(wh)
	if err != nil {
		return err
	}

	query := "UPDATE towns SET warehouse = ?, updatedAt = ? WHERE id = ?"
	_, err = t.db.ExecContext(ctx, query, whBytes, time.Now().UTC(), tid)
	if err != nil {
		return err
	}

	t.warehouseCache.Set(townID.String(), wh, CacheDurationWarehouse)

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
