package services

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/gerbenjacobs/millwheat/game"
	"github.com/gerbenjacobs/millwheat/storage"
)

type BattleSvc struct {
	storage storage.BattleStorage
}

// TODO: fix battle and army ID
var TMPCurrentBattleId = uuid.MustParse("7e19d988-201c-4cab-b1c4-29c5864d88fe")
var TMPArmyId = uuid.MustParse("70fd8839-2ea4-4e9d-9c90-5ec005448634")

func NewBattleSvc(storage storage.BattleStorage) *BattleSvc {
	return &BattleSvc{
		storage: storage,
	}
}

func (b *BattleSvc) Season(ctx context.Context) (*game.Season, error) {
	return &game.Season{
		ID:      uuid.New(),
		Name:    "Spring",
		Year:    1,
		Start:   time.Now().UTC(),
		End:     time.Now().UTC().Add(480 * time.Hour),
		Battles: nil,
	}, nil
}

func (b *BattleSvc) LastBattle(ctx context.Context) (*game.Battle, error) {
	return &game.Battle{
		ID:    uuid.MustParse("34461712-9fb6-4e02-b456-01d8f35341c8"),
		Name:  "Battle of Aria",
		Start: time.Now().UTC(),
		End:   time.Now().UTC().Add(10 * time.Hour),
		Attackers: game.Army{
			Name:  "Alyria",
			Score: 1023,
		},
		Defenders: game.Army{
			Name:  "Herkoonni",
			Score: 957,
		},
	}, nil
}

func (b *BattleSvc) UpcomingBattle(ctx context.Context) (*game.Battle, error) {
	return &game.Battle{
		ID:    uuid.MustParse("7268d57c-359e-47f2-91c8-ccac7c4d3ccf"),
		Name:  "Battle of Wulraven",
		Start: time.Now().UTC().Add(240 * time.Hour),
	}, nil
}

func (b *BattleSvc) AddWarrior(ctx context.Context, battleId, armyId, townId uuid.UUID, warriorType game.WarriorType, quantity int) error {
	return b.storage.AddWarrior(ctx, battleId, armyId, townId, warriorType, quantity)
}

func (b *BattleSvc) MyWarriors(ctx context.Context) ([]game.Warrior, error) {
	return b.storage.CurrentWarriors(ctx, TMPCurrentBattleId, TMPArmyId, TownFromContext(ctx))
}
