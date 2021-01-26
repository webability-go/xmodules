package tools

import (
	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"
)

func BuildConfigSet(config *xconfig.XConfig) xcore.XDataset {
	data := xcore.XDataset{}
	for id := range config.Parameters {
		d, _ := config.Get(id)
		if val, ok := d.(*xconfig.XConfig); ok {
			data[id] = BuildConfigSet(val)
		} else {
			data[id] = d
		}
	}
	return data
}
