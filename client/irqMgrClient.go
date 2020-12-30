package main

import (
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/suykerbuyk/irqmgr"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type CmdLineOpts struct {
	help   bool
	server string
}

var (
	opts        CmdLineOpts
	irqmgrFlags = flag.NewFlagSet("irqmgr", flag.ExitOnError)
)

func init() {
	opts.help = false
	irqmgrFlags.StringVarP(&opts.server, "server", "s", "http://127.0.0.1:10000", "Server URL")
	irqmgrFlags.BoolVarP(&opts.help, "help", "h", false, "Show Help")
	irqmgrFlags.Usage = usage
	if err := irqmgrFlags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]...\n", os.Args[0])
	irqmgrFlags.PrintDefaults()
	os.Exit(0)
}
func main() {
	if opts.help {
		usage()
	}
	jsonGetURI := opts.server + "/all"
	response, err := http.Get(jsonGetURI)
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var irqTallies irqmgr.IrqTallies
	json.Unmarshal(responseData, &irqTallies)
	loadJsonTree(responseData)
	fmt.Println("Done")
}

//affinity_hint  effective_affinity  effective_affinity_list  node  smp_affinity  smp_affinity_list  spurious
