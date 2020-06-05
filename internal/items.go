package internal

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	app "github.com/gerbenjacobs/millwheat"
)

func MustReadItemsForWarehouse(filePath string) app.Warehouse {
	var i []app.Item
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("failed to read " + filePath + ": " + err.Error())
	}
	if err := yaml.Unmarshal(b, &i); err != nil {
		panic("failed to unmarshal" + filePath + ": " + err.Error())
	}

	return app.NewWarehouse(i)
}
