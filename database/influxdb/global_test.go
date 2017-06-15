package influxdb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mcasviper/yanic/data"
	"github.com/mcasviper/yanic/runtime"
)

func TestGlobalStats(t *testing.T) {
	stats := runtime.NewGlobalStats(createTestNodes())

	assert := assert.New(t)
	fields := GlobalStatsFields(stats)

	// check fields
	assert.EqualValues(3, fields["nodes"])
}

func createTestNodes() *runtime.Nodes {
	nodes := runtime.NewNodes(&runtime.Config{})

	nodeData := &data.ResponseData{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 23,
			},
		},
		NodeInfo: &data.NodeInfo{
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	}
	nodeData.NodeInfo.Software.Firmware.Release = "2016.1.6+entenhausen1"
	nodes.Update("abcdef012345", nodeData)

	nodes.Update("112233445566", &data.ResponseData{
		Statistics: &data.Statistics{
			Clients: data.Clients{
				Total: 2,
			},
		},
		NodeInfo: &data.NodeInfo{
			Hardware: data.Hardware{
				Model: "TP-Link 841",
			},
		},
	})

	nodes.Update("0xdeadbeef0x", &data.ResponseData{
		NodeInfo: &data.NodeInfo{
			VPN: true,
			Hardware: data.Hardware{
				Model: "Xeon Multi-Core",
			},
		},
	})

	return nodes
}
