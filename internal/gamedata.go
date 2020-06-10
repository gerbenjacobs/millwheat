package internal

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/gerbenjacobs/millwheat/game"
)

func MustReadItems(filePath string) game.Items {
	var i []game.Item
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic("failed to read " + filePath + ": " + err.Error())
	}
	if err := yaml.Unmarshal(b, &i); err != nil {
		panic("failed to unmarshal" + filePath + ": " + err.Error())
	}

	return game.NewItems(i)
}
