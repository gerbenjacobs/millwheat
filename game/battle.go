package game

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	WarriorSword WarriorType = iota
	WarriorCrossbow
	WarriorLance
)

var WarriorTypes = []WarriorType{WarriorSword, WarriorCrossbow, WarriorLance}

var WarriorCosts = map[WarriorType]ItemSetSlice{
	WarriorSword: {
		{ItemID: "sword", Quantity: 1},
		{ItemID: "wooden_shield", Quantity: 1},
		{ItemID: "wine", Quantity: 1},
		{ItemID: "bread", Quantity: 1},
	},
	WarriorCrossbow: {
		{ItemID: "crossbow", Quantity: 1},
		{ItemID: "leather_armour", Quantity: 1},
		{ItemID: "wine", Quantity: 1},
		{ItemID: "bread", Quantity: 1},
	},
	WarriorLance: {
		{ItemID: "lance", Quantity: 1},
		{ItemID: "iron_platearmour", Quantity: 1},
		{ItemID: "horse", Quantity: 1},
		{ItemID: "wine", Quantity: 1},
		{ItemID: "meat", Quantity: 1},
	},
}

type WarriorType int

type Season struct {
	ID      uuid.UUID
	Name    string
	Year    int
	Start   time.Time
	End     time.Time
	Battles []Battle
}

type Battle struct {
	ID        uuid.UUID
	Name      string
	Start     time.Time
	End       time.Time
	Attackers Army
	Defenders Army
}

type Army struct {
	ID       uuid.UUID
	Name     string
	Score    int
	Warriors []Warrior
}

type Warrior struct {
	Type     WarriorType
	Quantity int
}

func CalculateWarriorCosts(warriorType WarriorType, quantity int) (ItemSetSlice, error) {
	costs, ok := WarriorCosts[warriorType]
	if !ok {
		return nil, errors.New("warrior type not found")
	}

	var iss ItemSetSlice
	for _, i := range costs {
		item := i
		item.Quantity = i.Quantity * quantity
		iss = append(iss, item)
	}

	return iss, nil
}

func WarriorTypeFromString(w string) (WarriorType, error) {
	switch w {
	case "0":
		return WarriorSword, nil
	case "1":
		return WarriorCrossbow, nil
	case "2":
		return WarriorLance, nil
	default:
		return -1, errors.New("warrior type unknown")
	}
}

func (w Warrior) Image() string {
	switch w.Type {
	case WarriorSword:
		return "/images/items/sword.png"
	case WarriorCrossbow:
		return "/images/items/crossbow.gif"
	case WarriorLance:
		return "/images/items/lance.gif"
	default:
		logrus.WithField("warrior", w).Warn("invalid warrior image")
		return "/images/items/plank.png"
	}
}

func (w WarriorType) String() string {
	switch w {
	case WarriorSword:
		return "Infantry"
	case WarriorCrossbow:
		return "Archers"
	case WarriorLance:
		return "Cavalry"
	default:
		return "Unknown type"
	}
}
