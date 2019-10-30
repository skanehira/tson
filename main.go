package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/skanehira/tson/gui"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	enableLog = flag.Bool("log", false, "enable log")
)

func init() {
	flag.Parse()
	if *enableLog {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		logWriter, err := os.OpenFile(filepath.Join(home, "tson.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		log.SetOutput(logWriter)
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

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
