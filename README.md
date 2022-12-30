# NPP-Prepper
A simple static binary that allows the creation of vCenter's Network Protocol Profiles (within vCenter's API referred to as IpPools).

Existing tooling like `govc` or `PowerCLI` can either not be used to configure IpPools, or adds unreasonably large overhead (container image `vmware/powerclicore` is nearly 1GB large). In fact, VMware's documentation *only* mentions how to use this feature through the use of vCenter's webinterface, which is entirely pointless in nowadays' world of automation.

`npp-prepper` was born out of necessity; one of my projects involves bootstrapping OVA's in greenfield environments, with the explicit design choice to not rely on external processes, like DHCP.

### Usage
```bash
Usage:
  npp-prepper [OPTIONS] <command>

Application Options:
  -s, --server=   FQDN of the vCenter appliance
  -u, --username= Username to authenticate with
  -p, --password= Password to authenticate with

Help Options:
  -h, --help      Show this help message

Available commands:
  datacenter      Define a Network Protocol Profile within a datacenter
  dc
  guestos         Configure guest OS network with allocated IP address
  os
  virtualmachine  Configure a virtual machine for usage of Network Protocol Profiles
  vm
```

#### Subcommand `datacenter` (`dc`) usage
```bash
Usage:
  npp-prepper [OPTIONS] dc [dc-OPTIONS]

Application Options:
  -s, --server=           FQDN of the vCenter appliance
  -u, --username=         Username to authenticate with
  -p, --password=         Password to authenticate with

Help Options:
  -h, --help              Show this help message

[dc command options]
      -n, --name=         Name of datacenter
      -p, --portgroup=    Name of network portgroup
          --startaddress=
          --endaddress=
          --netmask=
          --dnsserver=
          --dnsdomain=
          --gateway=
      -f, --force
```

#### Subcommand `guestos` (`os`) usage
// TODO

#### Subcommand `virtualmachine` (`vm`) usage
```bash
Usage:
  npp-prepper [OPTIONS] vm [vm-OPTIONS]

Application Options:
  -s, --server=         FQDN of the vCenter appliance
  -u, --username=       Username to authenticate with
  -p, --password=       Password to authenticate with

Help Options:
  -h, --help            Show this help message

[vm command options]
      -d, --datacenter= Name of datacenter
      -n, --name=       Name of virtual machine
      -p, --portgroup=  Name of network portgroup
```

### Future plans
- Create a container image and publish it
- Add support for various Linux network implementations to automatically configure the guest OS of a VM.
