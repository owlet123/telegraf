package cpustat

import (
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

type GopsuitlsCpu interface {
	Info() ([]cpu.InfoStat, error)
}

type Gc struct {}

func (gc *Gc) Info() ([]cpu.InfoStat, error) {return cpu.Info()}

type Cpustat struct {
	GopsuitlsCpu
}

func NewCpustat(csi GopsuitlsCpu) *Cpustat {
	return &Cpustat{csi}
}

func (c *Cpustat) Description() string {
	return `
# cpustat.count - number of CPUs
# cpustat.cpu.mhz - CPU's MHz
# cpustat.cpu.cores - CPU's number of cores`
}

func (c *Cpustat) SampleConfig() string {
	return ``
}

func (c *Cpustat) Gather(acc telegraf.Accumulator) error {
	get_cpustat(acc, c)

	return nil
}

func init() {
	inputs.Add("cpustat", func() telegraf.Input {
		return &Cpustat{&Gc{}}
	})
}

func get_cpustat(acc telegraf.Accumulator, c *Cpustat) {
	cpus, err := c.Info()
	if err != nil {
		acc.AddError(fmt.Errorf("Error while getting cpu info: %s\n", err))
	}

	tags := map[string]string{}

	fields := map[string]interface{}{
		"count": len(cpus),
	}

	acc.AddFields("cpustat", fields, tags, time.Now())

	for _, c := range cpus {
		tags := map[string]string{
			"model":     c.ModelName,
			"processor": fmt.Sprint(c.CPU),
		}

		fields := map[string]interface{}{
			"cpu.cores": c.Cores,
			"cpu.mhz":   c.Mhz,
		}

		acc.AddFields("cpustat", fields, tags, time.Now())
	}
}
