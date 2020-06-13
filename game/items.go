package game

type ItemID string

type Item struct {
	ID          ItemID
	Name        string
	Description string
	Image       string
}

type Items map[ItemID]Item

// WarehouseItem represents an instance of an item in a warehouse
type WarehouseItem struct {
	ItemID   ItemID
	Quantity int
}
