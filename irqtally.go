package irqmgr

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const pathInterrupts = "/proc/interrupts"

type IrqCount uint64

/* IrqCpuAffinity - Struct to capture CPU affinity for HW IRQs
   see: https://github.com/torvalds/linux/blob/bebc6082da0a9f5d47a1ea2edc099bf671058bd4/include/linux/irq.h#L138
*/
type IrqCpuAffinity struct {
	// Integer value of the hardware interrupt, not sequential!
	NumericInterruptValue uint `json:"NumericInterruptValue"`
	//This will need to become a bit mask of arbitrary length (len == core count)
	SmpAffinity string `json:"SmpAffinity"`
	/* This needs to be codified to reflect a (possibly) stepped range, that is portable across the wire.
	   What would be cool is if we could find a way to imply sockets vs hardware threads. l2, and l3 are socket based  */
	SmpAffinityList string `json:"SmpAffinityList"`
	// Not entirely sure how this is leveraged by anything except irqbalance
	AffinityHint string `json:"AffinityHint"`
	/* To be researched */
	EffectiveAffinity string `json:"EffectiveAffinity"`
	/* To be researched */
	EffectiveAffinityList string `json:"EffectiveAffinityList"`
}

type SetSmpAffinityData struct {
	Irq      string `json:"Irq"`
	Affinity string `json:"Affinity"`
}

type IrqsServicedTally struct {
	NumericInterruptValue   uint       `json:"NumericInterruptValue"`
	InterruptsServicedByCPU []IrqCount `json:"InterruptsServicedByCPU"`
	CpuSmpAffinity          string     `json:"CpuSmpAffinity"`
	SourceOfHwInterrupt     string     `json:"SourceOfHwInterrupt"`
}
type IrqTallies struct {
	TotalCpuCount     int                 `json:"TotalCpuCount"`
	TotalNumericIRQs  int                 `json:"TotalIrqCount"`
	IrqsServicedByCPU []IrqsServicedTally `json:"IrqsServicedByCPU"`
}

var irqTalliesPrevious *IrqTallies = nil
var irqTalliesCurrent *IrqTallies = nil
var irqTalliesDelta *IrqTallies = nil

// Stub of what needs to be added
func calcIrqTalliesDeltas() error {
	if irqTalliesPrevious == nil {
		irqTallies, err := FetchIrqs()
		if err != nil {
			return err
		}
		// Need to return an error code that indicates we don't yet have the data.
		irqTalliesPrevious = irqTallies
		return nil
	}
	// calculate the delta, update the global structs.
	return nil
}

func ReadIrqCpuAffinity(irq uint) string {
	var path string = "/proc/irq/" + strconv.Itoa(int(irq)) + "/smp_affinity"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	ret := string(data[:])
	ret = strings.Replace(ret, "\n", "", -1)
	return ret
}
func SetIrqAffinity(irq, value string) {
	path := "/proc/irq/" + irq + "/smp_affinity"
	//path := irq + "_smp_affinity.test"
	err := ioutil.WriteFile(path, []byte(value), 0)
	if err != nil {
		log.Fatal(err)
	}
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
	// Read each line the buffer, until we fail to convert a column 1 (irq number) to an int.
	for idx, line := range lines {
		if idx == 0 {
			// First line is the header, with a column header for each CPU, count them.
			irqTallies.TotalCpuCount = len(strings.Fields(line))
		} else {
			if irqTallies.TotalCpuCount < 1 {
				errStr := "Could not detect the number of CPUs"
				log.Println(errStr)
				return nil, errors.New(errStr)
			}
			var irqTally IrqsServicedTally
			splits := strings.Fields(line)
			if _, err := fmt.Sscanf(splits[0], "%d:", &irqTally.NumericInterruptValue); err == nil {
				// Parse out the numeric totals for IRQs seen per CPU.
				for i := 1; i < (irqTallies.TotalCpuCount - 1); i++ {
					var irqCnt IrqCount
					fmt.Sscanf(splits[i], "%d", &irqCnt)
					irqTally.InterruptsServicedByCPU = append(irqTally.InterruptsServicedByCPU, irqCnt)
				}
				irqTally.CpuSmpAffinity = ReadIrqCpuAffinity(irqTally.NumericInterruptValue)
				// Combine the strings following per CPU Irq counts into an SourceOfHwInterrupt string
				for i := irqTallies.TotalCpuCount; i < len(splits); i++ {
					if len(irqTally.SourceOfHwInterrupt) != 0 {
						irqTally.SourceOfHwInterrupt += " "
					}
					irqTally.SourceOfHwInterrupt += strings.TrimSpace(splits[i])
				}
				irqTallies.IrqsServicedByCPU = append(irqTallies.IrqsServicedByCPU, irqTally)
			} else {
				// How many numbered IRQs did we inventory?
				irqTallies.TotalNumericIRQs = len(irqTallies.IrqsServicedByCPU)
				//We've hit a non-numeric (numbered IRQ), bail out of top for/loop.
				break
			}
		}
	}
	return &irqTallies, nil
}
