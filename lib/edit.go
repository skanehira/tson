package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/skanehira/tson/gui"
)

// Edit use tson as a library
func Edit(b []byte) ([]byte, error) {
	// dont output log
	log.SetOutput(ioutil.Discard)

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
