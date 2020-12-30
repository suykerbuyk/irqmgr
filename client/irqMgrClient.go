package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/suykerbuyk/irqmgr"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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
func getNewIrqAffinitySetting(irq string, affinity string) string {
	var newAffinity string
	fmt.Println("Set smp_affinity for IRQ:", irq, "currently:", affinity)
	fmt.Println("Hit enter without entering a value to cancel")
	fmt.Scanf("%s", &newAffinity)
	if newAffinity == "" {
		return ""
	}
	return newAffinity
}
func main() {
	if opts.help {
		usage()
	}
	jsonGetURI := opts.server + "/all"
	jsonSetURI := opts.server + "/set"
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
	editPath, editValue := loadJsonTree(responseData)
	if strings.Contains(editPath, "CpuSmpAffinity") {
		editPath = strings.ReplaceAll(editPath, `["IrqsServicedByCPU"][`, "")
		editPath = strings.ReplaceAll(editPath, `]["CpuSmpAffinity"]`, "")
		editValue = strings.ReplaceAll(editValue, `"`, "")
		editValue := getNewIrqAffinitySetting(editPath, editValue)
		if editValue != "" {
			var setting irqmgr.SetSmpAffinityData
			setting.Irq = editPath
			setting.Affinity = editValue
			jsonReq, err := json.Marshal(setting)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := http.Post(jsonSetURI, "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString := string(bodyBytes)
			fmt.Println(bodyString)
		}
	}
}

//affinity_hint  effective_affinity  effective_affinity_list  node  smp_affinity  smp_affinity_list  spurious
