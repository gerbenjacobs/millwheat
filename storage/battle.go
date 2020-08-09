package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
)

type BattleRepo struct {
	db *sql.DB
}

func NewBattleRepo(db *sql.DB) *BattleRepo {
	return &BattleRepo{db: db}
}

func (b *BattleRepo) AddWarrior(ctx context.Context, battleId, armyId, townId uuid.UUID, warriorType game.WarriorType, quantity int) error {
	battle, _ := battleId.MarshalBinary()
	army, _ := armyId.MarshalBinary()
	town, _ := townId.MarshalBinary()

	// write warrior to database
	stmt, err := b.db.PrepareContext(ctx, "INSERT INTO warriors (battleId, armyId, townId, warriorType, quantity) VALUES(?, ?, ?, ?, ?)  ON DUPLICATE KEY UPDATE quantity = quantity + ?")
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, battle, army, town, warriorType, quantity, quantity)

	return err
}

func (b *BattleRepo) WarriorsFromTown(ctx context.Context, townId, battleId uuid.UUID) ([]game.Warrior, error) {
	panic("implement me")
}

func (b *BattleRepo) AllWarriorsForBattle(ctx context.Context, battleId uuid.UUID) ([]game.Army, error) {
	panic("implement me")
}

func (b *BattleRepo) CurrentWarriors(ctx context.Context, battleId, armyId, townId uuid.UUID) ([]game.Warrior, error) {
	battle, _ := battleId.MarshalBinary()
	army, _ := armyId.MarshalBinary()
	town, _ := townId.MarshalBinary()

	rows, err := b.db.QueryContext(ctx, "SELECT warriorType, quantity FROM warriors WHERE battleId = ? AND armyId = ? AND townId = ? ORDER BY warriorType", battle, army, town)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warriors []game.Warrior
	for rows.Next() {
		var w game.Warrior
		err = rows.Scan(&w.Type, &w.Quantity)
		if err != nil {
			return nil, err
		}
		warriors = append(warriors, w)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return warriors, nil
}
