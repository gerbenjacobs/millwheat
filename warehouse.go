package millwheat

type ItemID string

type Item struct {
	ID          ItemID
	Name        string
	Description string
	Image       string
}

type Warehouse map[ItemID]Item

func NewWarehouse(items []Item) Warehouse {
	w := Warehouse{}
	for _, i := range items {
		w[i.ID] = i
	}

	return w
}
