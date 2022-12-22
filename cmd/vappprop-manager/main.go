package main

import (
	"context"
	"flag"
	"log"

	"spamasaurus.com/m/pkg/config"
	"spamasaurus.com/m/pkg/hypervisor"
)

type Input struct {
	VirtualMachine string
	Network        string
}

func main() {
	var input Input

	flag.StringVar(&input.VirtualMachine, "vm", "", "name of VM")
	flag.StringVar(&input.Network, "network", "", "name of network portgroup")
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Invalid configuration: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clt, err := hypervisor.NewClient(ctx, cfg.Hypervisor.Url, cfg.Hypervisor.Username, cfg.Hypervisor.Password, true)
	if err != nil {
		log.Fatalf("Login failed: %s", err)
	}

	fnd, err := hypervisor.DatacenterFinder(ctx, clt, cfg.Hypervisor.Datacenter)
	if err != nil {
		log.Fatalf("Foo indeed: %s", err)
	}

	if err := hypervisor.SetVirtualMachineProperties(ctx, fnd, input.VirtualMachine, input.Network); err != nil {
		log.Fatalf("Could not apply vApp properties: %s", err)
	}
}
