package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/nomadcoin/explorer"
	"github.com/nomadcoin/rest"
)

func usage() {
	fmt.Printf("welcome to nomadcoin\n\n")
	fmt.Printf("please use the following commands:\n\n")
	fmt.Printf("explorer: 	The start the HTML Explorer\n")
	fmt.Printf("rest: 		The start the REST API(recommended)\n")
	runtime.Goexit()
}

func Start() {

	if len(os.Args) < 2 {
		usage()
	}

	port := flag.Int("port", 4000, "Set of the server port")

	mode := flag.String("mode", "rest", "")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	}
}
