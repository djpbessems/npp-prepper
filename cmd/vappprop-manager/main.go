package main

import (
	"context"
	"log"
	"os"

	"github.com/jessevdk/go-flags"

	"spamasaurus.com/m/pkg/hypervisor"
)

func main() {
	var opts struct {
		FQDN     string `short:"s" long:"server" description:"FQDN of the vCenter appliance" required:"true"`
		Username string `short:"u" long:"username" description:"Username to authenticate with" required:"true"`
		Password string `short:"p" long:"password" description:"Password to authenticate with" required:"true"`

		Datacenter     string `short:"d" long:"datacenter" description:"Name of datacenter" required:"true"`
		VirtualMachine string `short:"v" long:"virtualmachine" description:"Name of virtual machine" required:"true"`
		Network        string `short:"n" long:"network" description:"Name of network portgroup" required:"true"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		if e, ok := err.(*flags.Error); ok {
			if e.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clt, err := hypervisor.NewClient(ctx, opts.FQDN, opts.Username, opts.Password, true)
	if err != nil {
		log.Fatalf("[ERROR] Login failed: %s", err)
	}

	fnd, err := hypervisor.DatacenterFinder(ctx, clt, opts.Datacenter)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine datacenter: %s", err)
	}

	if err := hypervisor.SetVirtualMachineProperties(ctx, fnd, opts.VirtualMachine, opts.Network); err != nil {
		log.Fatalf("[ERROR] Could not apply vApp properties: %s", err)
	}
}
