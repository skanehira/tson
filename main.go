package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/skanehira/tson/gui"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	enableLog = flag.Bool("log", false, "enable log")
	url       = flag.String("url", "", "get json from url")
)

func printError(err error) int {
	fmt.Fprintln(os.Stderr, err)
	return 1
}

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
			os.Exit(printError(err))
		}

		log.SetOutput(logWriter)
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

func run() int {
	if *url != "" {
		resp, err := http.Get(*url)
		if err != nil {
			return printError(err)
		}

		i, err := gui.UnMarshalJSON(resp.Body)
		if err != nil {
			return printError(err)
		}

		if err := gui.New().Run(i); err != nil {
			return printError(err)
		}
		return 0
	}

	if !terminal.IsTerminal(0) {
		i, err := gui.UnMarshalJSON(os.Stdin)
		if err != nil {
			return printError(err)
		}

		if err := gui.New().Run(i); err != nil {
			log.Println(err)
			return 1
		}
	}
	return 0
}

func main() {
	os.Exit(run())
}
