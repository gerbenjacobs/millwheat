package game

type ItemID string

type Item struct {
	ID          ItemID
	Name        string
	Description string
	Image       string
}

type Items map[ItemID]Item

func NewItems(items []Item) Items {
	list := Items{}
	for _, i := range items {
		list[i.ID] = i
	}

	return list
}
