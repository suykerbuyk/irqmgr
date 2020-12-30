package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/suykerbuyk/irqmgr"
	"log"
	"net/http"
	//	"os"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	staticHome := `<h1 style="text-align: center;">IrqMgr</h1>
	<h3 style="text-align: center;">John Suykerbuyk</h3>
	<h2>From here you can:</h2>
	<ul>
	<li><a href="/all">View Every All IRQ Info</a></li>
	<li><a href="/just">View just the IRQ Tallies</a></li>
	</ul>`
	fmt.Println("Endpoint Hit: homepage")
	fmt.Fprintf(w, staticHome)
}
func serveAllIrqTallies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: serveIrqTallies")
	irqTallies, err := irqmgr.FetchIrqs()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(irqTallies)
}
func setIrqAffinity(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: setIrqAffinity")
	var affinityReq irqmgr.SetSmpAffinityData
	err := json.NewDecoder(r.Body).Decode(&affinityReq)
	if err != nil {
		log.Fatal(err)
	}
	irqmgr.SetIrqAffinity(affinityReq.Irq, affinityReq.Affinity)
}

func serveJustIrqTallies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: serveIrqTallies")
	irqTallies, err := irqmgr.FetchIrqs()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(irqTallies.IrqsServicedByCPU)
}

func handleRequest() {
	fmt.Println("Entering handle request")
	irqMgrRouter := mux.NewRouter().StrictSlash(true)
	irqMgrRouter.HandleFunc("/", homePage)
	irqMgrRouter.HandleFunc("/all", serveAllIrqTallies)
	irqMgrRouter.HandleFunc("/just", serveJustIrqTallies)
	irqMgrRouter.HandleFunc("/set", setIrqAffinity).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", irqMgrRouter))
}

func main() {
	irqTallies, err := irqmgr.FetchIrqs()
	out, err := json.Marshal(irqTallies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	handleRequest()
}

//affinity_hint  effective_affinity  effective_affinity_list  node  smp_affinity  smp_affinity_list  spurious
