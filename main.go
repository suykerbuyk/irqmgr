package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	</ul>`
	fmt.Println("Endpoint Hit: homepage")
	fmt.Fprintf(w, staticHome)
}
func handleRequest() {
	fmt.Println("Entering handle request")
	irqMgrRouter := mux.NewRouter().StrictSlash(true)
	irqMgrRouter.HandleFunc("/", homePage)
	irqMgrRouter.HandleFunc("/all", serveIrqTallies)
	log.Fatal(http.ListenAndServe(":10000", irqMgrRouter))
}

func serveIrqTallies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: serveIrqTallies")
	irqTallies, err := FetchIrqs()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(irqTallies)
}

func main() {
	irqTallies, err := FetchIrqs()
	out, err := json.Marshal(irqTallies)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	handleRequest()
}

//affinity_hint  effective_affinity  effective_affinity_list  node  smp_affinity  smp_affinity_list  spurious
