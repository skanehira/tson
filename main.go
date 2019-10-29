package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/skanehira/tson/gui"
	"golang.org/x/crypto/ssh/terminal"
)

func run() int {
	if !terminal.IsTerminal(0) {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if len(b) == 0 {
			fmt.Println("json is empty")
			return 0
		}

		var t interface{}
		if err := json.Unmarshal(b, &t); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		if err := gui.New().Run(t); err != nil {
			log.Println(err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
