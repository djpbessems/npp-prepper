package main

import (
	"context"
	"flag"
	"log"

	"spamasaurus.com/m/pkg/config"
	"spamasaurus.com/m/pkg/hypervisor"
)

func main() {
	var virtualmachine string

	flag.StringVar(&virtualmachine, "vm", "", "name of VM")
	flag.Parse()

	if virtualmachine == "" {
		log.Fatalf("Name of vm is required")
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
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

	if err := hypervisor.SetVirtualMachineProperties(ctx, fnd, virtualmachine); err != nil {
		log.Fatalf("Could not apply vApp properties: %s", err)
	}
}
