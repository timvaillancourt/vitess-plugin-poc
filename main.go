package main

import (
	"flag"
	"fmt"
	"plugin"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

var pluginPath string

func init() {
	flag.StringVar(&pluginPath, "plugin-path", "", "path to durability policy plugin")
	flag.Parse()
}

func main() {
	if pluginPath == "" {
		panic("must define -plugin-path")
	}
	fmt.Printf("plugin path: %s\n", pluginPath)

	p, err := plugin.Open(pluginPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("plugin: %v\n", p)

	crossCell, err := p.Lookup("DurabilityCrossCell")
	if err != nil {
		panic(err)
	}
	RegisterDurability("cross_cell", func() Durabler {
		return crossCell.(Durabler)
	})

	dur, err := GetDurabilityPolicy("cross_cell")
	if err != nil {
		panic(err)
	}

	fmt.Println(PromotionRule(dur, &topodatapb.Tablet{}))
}
