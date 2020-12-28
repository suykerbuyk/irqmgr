package main

import (
	"encoding/json"
	"fmt"
	//	"io/ioutil"
	"log"
	"os"
	//	"strings"
)

func main() {
	irqTallies, err := FetchIrqs()
	out, err := json.Marshal(irqTallies)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
