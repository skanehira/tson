package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func run() int {
	if !terminal.IsTerminal(0) {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Println(err)
			return 1
		}
		var t interface{}
		if err := json.Unmarshal(b, &t); err != nil {
			log.Println(err)
			return 1
		}

		fmt.Println(t)
	}

	return 0
}

func main() {
	os.Exit(run())
}
