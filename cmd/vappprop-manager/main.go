package main

import (
	"context"
	"log"
	"os"

	"github.com/jessevdk/go-flags"

	"spamasaurus.com/m/pkg/hypervisor"
)

var Global struct {
	FQDN     string `short:"s" long:"server" description:"FQDN of the vCenter appliance" required:"true"`
	Username string `short:"u" long:"username" description:"Username to authenticate with" required:"true"`
	Password string `short:"p" long:"password" description:"Password to authenticate with" required:"true"`
}
var Commands struct {
	Datacenter struct {
		Name    string `short:"n" long:"name" description:"Name of datacenter" required:"true"`
		Network string `short:"p" long:"portgroup" description:"Name of network portgroup" required:"true"`
	} // `command:"datacenter" alias:"dc" description:"Define a Network Protocol Profile within a datacenter"`
	VirtualMachine struct {
		Datacenter string `short:"d" long:"datacenter" description:"Name of datacenter" required:"true"`
		Name       string `short:"n" long:"name" description:"Name of virtual machine" required:"true"`
		Network    string `short:"p" long:"portgroup" description:"Name of network portgroup" required:"true"`
	} // `command:"virtualmachine" alias:"vm" description:"Configure a virtual machine for usage of Network Protocol Profiles"`
}

func main() {
	parser := flags.NewParser(&Global, flags.Default)
	parser.AddCommand("datacenter", "Define a Network Protocol Profile within a datacenter", "", &Commands.Datacenter)
	parser.AddCommand("dc", "", "", &Commands.Datacenter)
	parser.AddCommand("virtualmachine", "Configure a virtual machine for usage of Network Protocol Profiles", "", &Commands.VirtualMachine)
	parser.AddCommand("vm", "", "", &Commands.VirtualMachine)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clt, err := hypervisor.NewClient(ctx, Global.FQDN, Global.Username, Global.Password, true)
	if err != nil {
		log.Fatalf("[ERROR] Login failed: %s", err)
	}

	switch parser.Active.Name {
	case "datacenter", "dc":
		if err := hypervisor.CreateNetworkProtocolProfile(ctx, clt, Commands.Datacenter.Name, Commands.Datacenter.Network); err != nil {
			log.Fatalf("[ERROR] Could not create network protocol profile: %s", err)
		}

		log.Printf("[SUCCESS] Hurray! We did it!")
	case "virtualmachine", "vm":
		if err := hypervisor.SetVirtualMachineProperties(ctx, clt, Commands.VirtualMachine.Datacenter, Commands.VirtualMachine.Name, Commands.VirtualMachine.Network); err != nil {
			log.Fatalf("[ERROR] Could not apply vApp properties: %s", err)
		}

		log.Printf("[SUCCESS] Network protocol profile properties added to virtual machine '%s' and configured for network '%s'", Commands.VirtualMachine.Name, Commands.VirtualMachine.Network)
	}
}
