package main

import (
	"encoding/json"
	"fmt"
	"github.com/suykerbuyk/irqmgr"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	response, err := http.Get("http://localhost:10000/all")
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var irqTallies irqmgr.IrqTallies
	json.Unmarshal(responseData, &irqTallies)
	fmt.Println(irqTallies)
}

//affinity_hint  effective_affinity  effective_affinity_list  node  smp_affinity  smp_affinity_list  spurious
