package main

import (
	"context"
	"flag"
	"log"

	"spamasaurus.com/m/pkg/hypervisor"
)

type Input struct {
	FQDN     string
	Username string
	Password string

	Datacenter     string
	VirtualMachine string
	Network        string
}

func main() {
	var input Input

	flag.StringVar(&input.FQDN, "server", "", "FQDN of the vCenter appliance")
	flag.StringVar(&input.Username, "username", "", "Username to authenticate with")
	flag.StringVar(&input.Password, "password", "", "Password to authenticate with")
	flag.StringVar(&input.Datacenter, "dc", "", "Name of datacenter")
	flag.StringVar(&input.VirtualMachine, "vm", "", "Name of VM")
	flag.StringVar(&input.Network, "network", "", "Name of network portgroup")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clt, err := hypervisor.NewClient(ctx, input.FQDN, input.Username, input.Password, true)
	if err != nil {
		log.Fatalf("Login failed: %s", err)
	}

	fnd, err := hypervisor.DatacenterFinder(ctx, clt, input.Datacenter)
	if err != nil {
		log.Fatalf("Foo indeed: %s", err)
	}

	if err := hypervisor.SetVirtualMachineProperties(ctx, fnd, input.VirtualMachine, input.Network); err != nil {
		log.Fatalf("Could not apply vApp properties: %s", err)
	}
}
