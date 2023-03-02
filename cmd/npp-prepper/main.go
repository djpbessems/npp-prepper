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
		Name         string   `short:"n" long:"name" description:"Name of datacenter" required:"true"`
		Network      string   `long:"portgroup" description:"Name of network portgroup" required:"true"`
		StartAddress string   `long:"startaddress" required:"true"`
		EndAddress   string   `long:"endaddress" required:"true"`
		Netmask      string   `long:"netmask" required:"true"`
		DnsServer    []string `long:"dnsserver" required:"true"`
		DnsDomain    string   `long:"dnsdomain" required:"true"`
		Gateway      string   `long:"gateway" required:"true"`
		Force        bool     `short:"f" long:"force"`
	} // `command:"datacenter" alias:"dc" description:"Define a Network Protocol Profile within a datacenter"`
	VirtualMachine struct {
		Name                string `short:"n" long:"name" description:"Name of virtual machine" required:"true"`
		Datacenter          string `long:"datacenter" description:"Name of datacenter" required:"true"`
		Network             string `long:"portgroup" description:"Name of network portgroup" required:"true"`
		DisableOvfTransport bool   `long:"disableovftransport" description:"Disable the OVF transport 'VMware Tools'"`
	} // `command:"virtualmachine" alias:"vm" description:"Configure a virtual machine for usage of Network Protocol Profiles"`
	GuestOS struct {
	} // `command:"guestos" alias:"os" description:"Configure guest OS network with allocated IP address"`
}

func main() {
	parser := flags.NewParser(&Global, flags.Default)
	parser.AddCommand("datacenter", "Define a Network Protocol Profile within a datacenter", "", &Commands.Datacenter)
	parser.AddCommand("dc", "", "", &Commands.Datacenter)
	parser.AddCommand("virtualmachine", "Configure a virtual machine for usage of Network Protocol Profiles", "", &Commands.VirtualMachine)
	parser.AddCommand("vm", "", "", &Commands.VirtualMachine)
	parser.AddCommand("guestos", "Configure guest OS network with allocated IP address", "", &Commands.GuestOS)
	parser.AddCommand("os", "", "", &Commands.GuestOS)

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
		if err := hypervisor.CreateNetworkProtocolProfile(ctx, clt, Commands.Datacenter.Name, Commands.Datacenter.Network, Commands.Datacenter.StartAddress, Commands.Datacenter.EndAddress, Commands.Datacenter.Netmask, Commands.Datacenter.DnsDomain, Commands.Datacenter.Gateway, Commands.Datacenter.DnsServer, Commands.Datacenter.Force); err != nil {
			log.Fatalf("[ERROR] Could not create network protocol profile: %s", err)
		}

		log.Printf("[SUCCESS] New network protocol profile created within datacenter '%s' (associated with network '%s')", Commands.Datacenter.Name, Commands.Datacenter.Network)
	case "guestos", "os":
		// TODO
	case "virtualmachine", "vm":
		if err := hypervisor.SetVirtualMachineProperties(ctx, clt, Commands.VirtualMachine.Datacenter, Commands.VirtualMachine.Name, Commands.VirtualMachine.Network, Commands.VirtualMachine.DisableOvfTransport); err != nil {
			log.Fatalf("[ERROR] Could not apply vApp properties: %s", err)
		}

		log.Printf("[SUCCESS] Network protocol profile properties added to virtual machine '%s' (configured for network '%s')", Commands.VirtualMachine.Name, Commands.VirtualMachine.Network)
	}
}
