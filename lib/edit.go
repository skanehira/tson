package lib

import (
	"encoding/json"
	"log"

	"github.com/skanehira/tson/gui"
)

// Edit use tson as a library
func Edit(b []byte) ([]byte, error) {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		log.Println(err)
		return nil, err
	}

	g := gui.New()
	if err := g.Run(i); err != nil {
		return nil, err
	}

	return json.Marshal(g.MakeJSON(g.Tree.GetRoot()))
}
