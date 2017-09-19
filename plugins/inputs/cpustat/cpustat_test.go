package cpustat

import (
	"fmt"
	"github.com/influxdata/telegraf/testutil"
	"github.com/shirou/gopsutil/cpu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
func Test_cpustat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cpustat integration tests.")
	}

	var acc testutil.Accumulator

	c := &Cpustat{}

	err := acc.GatherError(c.Gather)
	require.NoError(t, err)

	assert.True(t, acc.HasIntField("cpustat", "count"))

	assert.True(t, acc.HasMeasurement("cpustat"))

	assert.True(t, acc.HasInt32Field("cpustat", "cpu.cores"))
	assert.True(t, acc.HasFloatField("cpustat", "cpu.mhz"))
}
*/

type MockC struct {
	mock.Mock
}

func (m *MockC) Info() ([]cpu.InfoStat, error) {
	ret := m.Called()

	r0 := ret.Get(0).([]cpu.InfoStat)
	r1 := ret.Error(1)

	return r0, r1
}

func Test_Cpustat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cpustat integration tests.")
	}

	var mc MockC
	defer mc.AssertExpectations(t)

	ci := cpu.InfoStat{
		Cores:     int32(2),
		Mhz:       float64(2),
		ModelName: "modelName",
		CPU:       int32(1),
	}
	mc.On("Info").Return([]cpu.InfoStat{ci}, nil)

	var acc testutil.Accumulator

	mcsi := NewCpustat(&mc)
	err := mcsi.Gather(&acc)
	require.NoError(t, err)

	assert.True(t, acc.HasMeasurement("cpustat"))
	assert.True(t, acc.HasIntField("cpustat", "count"))
	assert.True(t, acc.HasInt32Field("cpustat", "cpu.cores"))
	assert.True(t, acc.HasFloatField("cpustat", "cpu.mhz"))

	tags := map[string]string{
		"model":     "modelName",
		"processor": fmt.Sprint(1),
	}
	fields := map[string]interface{}{
		"cpu.cores": int32(2),
		"cpu.mhz":   float64(2),
	}
	acc.AssertContainsTaggedFields(t, "cpustat", fields, tags)

	tags_count := map[string]string{}
	field_count := map[string]interface{}{
		"count": 1,
	}
	acc.AssertContainsTaggedFields(t, "cpustat", field_count, tags_count)
}
