package main

import (
	//	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	//	"os"
	"strings"
)

const pathInterrupts = "/proc/interrupts"

type IrqCount uint64
type IrqTally struct {
	IrqNum       uint       `json:"IrqNum"`
	CpuIrqCounts []IrqCount `json:"CpuIrqCounts"`
	Source       string     `json:"SouceOfInterrupt"`
}
type IrqTallies []IrqTally

func FetchIrqs() (*IrqTallies, error) {
	var irqTallies IrqTallies
	var cpuCnt int = 0
	buff, err := ioutil.ReadFile(pathInterrupts)
	if err != nil {
		errStr := "Could not read from " + pathInterrupts
		log.Println(errStr)
		return nil, errors.New(errStr)
	}
	lines := strings.Split(string(buff), "\n")
	for idx, line := range lines {
		if idx == 0 {
			cpuCnt = len(strings.Fields(line))
		} else {
			if cpuCnt < 1 {
				errStr := "Could not detect the number of CPUs"
				log.Println(errStr)
				return nil, errors.New(errStr)
			}
			var irqTally IrqTally
			splits := strings.Fields(line)
			if _, err := fmt.Sscanf(splits[0], "%d:", &irqTally.IrqNum); err == nil {
				for i := 1; i < (cpuCnt - 1); i++ {
					var irqCnt IrqCount
					fmt.Sscanf(splits[i], "%d", &irqCnt)
					irqTally.CpuIrqCounts = append(irqTally.CpuIrqCounts, irqCnt)
				}
				for i := cpuCnt; i < len(splits); i++ {
					if len(irqTally.Source) != 0 {
						irqTally.Source += " "
					}
					irqTally.Source += strings.TrimSpace(splits[i])
				}
				irqTallies = append(irqTallies, irqTally)
			} else {
				for _, irqTally := range irqTallies {
					var stringOfCounts string
					for _, irqCnt := range irqTally.CpuIrqCounts {
						stringOfCounts += fmt.Sprintf("%d", irqCnt) + " "
					}
				}
				break
			}
		}
	}
	//out, err := json.Marshal(IrqTallies)
	//if err != nil {
	//	log.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(string(out))
	return &irqTallies, nil
}
