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
	IrqSource    string     `json:"SouceOfInterrupt"`
}
type IrqTallies struct {
	CpuCount   int        `json:"TotalCpuCount"`
	IrqCount   int        `json:"TotalIrqCount"`
	IrqsPerCpu []IrqTally `json:"IrqsPerCpu"`
}

func FetchIrqs() (*IrqTallies, error) {
	var irqTallies IrqTallies
	buff, err := ioutil.ReadFile(pathInterrupts)
	if err != nil {
		errStr := "Could not read from " + pathInterrupts
		log.Println(errStr)
		return nil, errors.New(errStr)
	}
	lines := strings.Split(string(buff), "\n")
	// Read each line the buffer, until we fail to convert a column 1 irq number to an int.
	for idx, line := range lines {
		if idx == 0 {
			// First line is the header, with a column header for each CPU, count them.
			irqTallies.CpuCount = len(strings.Fields(line))
		} else {
			if irqTallies.CpuCount < 1 {
				errStr := "Could not detect the number of CPUs"
				log.Println(errStr)
				return nil, errors.New(errStr)
			}
			var irqTally IrqTally
			splits := strings.Fields(line)
			if _, err := fmt.Sscanf(splits[0], "%d:", &irqTally.IrqNum); err == nil {
				// Parse out the numeric totals for IRQs seen per CPU.
				for i := 1; i < (irqTallies.CpuCount - 1); i++ {
					var irqCnt IrqCount
					fmt.Sscanf(splits[i], "%d", &irqCnt)
					irqTally.CpuIrqCounts = append(irqTally.CpuIrqCounts, irqCnt)
				}
				// Combine the strings following per CPU Irq counts into an IrqSource string
				for i := irqTallies.CpuCount; i < len(splits); i++ {
					if len(irqTally.IrqSource) != 0 {
						irqTally.IrqSource += " "
					}
					irqTally.IrqSource += strings.TrimSpace(splits[i])
				}
				irqTallies.IrqsPerCpu = append(irqTallies.IrqsPerCpu, irqTally)
			} else {
				// How many numbered IRQs did we inventory?
				irqTallies.IrqCount = len(irqTallies.IrqsPerCpu)
				//We've hit a non-numeric (numbered IRQ), bail out of top for/loop.
				break
			}
		}
	}
	return &irqTallies, nil
}
