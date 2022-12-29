package hypervisor

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	netrange "github.com/netdata/go.d.plugin/pkg/iprange"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"spamasaurus.com/m/pkg/utils"
)

func CreateNetworkProtocolProfile(ctx context.Context, clt *vim25.Client, datacenter, network, startaddress, endaddress, netmask, dnsdomain, gateway string, dnsserver []string, force bool) error {
	finder := find.NewFinder(clt, true)
	dc, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine datacenter: %s", err)
	}
	finder.SetDatacenter(dc)
	nw, err := finder.Network(ctx, network)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine network: %s", err)
	}

	var networksummary mo.Network
	pc := property.DefaultCollector(clt)
	pc.Retrieve(ctx, []types.ManagedObjectReference{nw.Reference()}, []string{"summary"}, &networksummary)
	if networksummary.Summary.GetNetworkSummary().IpPoolId != nil {
		if force == true {
			request := &types.DestroyIpPool{
				This:  *clt.ServiceContent.IpPoolManager,
				Dc:    dc.Reference(),
				Id:    *networksummary.Summary.GetNetworkSummary().IpPoolId,
				Force: true,
			}
			if _, err := methods.DestroyIpPool(ctx, clt.RoundTripper, request); err != nil {
				log.Fatalf("[ERROR] Could not remove existing network protocol profile '%s'",
					networksummary.Summary.GetNetworkSummary().IpPoolName)
			}
		} else {
			log.Fatalf("[ERROR] Network '%s' already has existing network protocol profile '%s' associated; use the --force flag to replace it",
				network,
				networksummary.Summary.GetNetworkSummary().IpPoolName)
		}
	}

	iprange, err := netrange.ParseRange(fmt.Sprintf("%s-%s", startaddress, endaddress))
	if err != nil {
		log.Fatalf("[ERROR] Invalid IP range: %s", err)
	}
	ipnetwork := net.ParseIP(startaddress).Mask(net.IPMask(net.ParseIP(netmask).To4()))

	request := &types.CreateIpPool{
		This: *clt.ServiceContent.IpPoolManager,
		Dc:   dc.Reference(),
		Pool: types.IpPool{
			Name:          "ippool-" + uuid.New().String()[0:5],
			DnsDomain:     dnsdomain,
			DnsSearchPath: dnsdomain,
			NetworkAssociation: []types.IpPoolAssociation{{
				Network:     utils.MoRefAddr(nw.Reference()),
				NetworkName: network,
			}},
			Ipv4Config: &types.IpPoolIpPoolConfigInfo{
				SubnetAddress:       ipnetwork.String(),
				Netmask:             netmask,
				Gateway:             gateway,
				Range:               fmt.Sprintf("%s#%d", startaddress, iprange.Size()),
				Dns:                 dnsserver,
				DhcpServerAvailable: utils.BoolAddr(false),
				IpPoolEnabled:       utils.BoolAddr(true),
			},
		},
	}

	if _, err := methods.CreateIpPool(ctx, clt.RoundTripper, request); err != nil {
		log.Fatalf("[ERROR]: Failed creating new network protocol profile (for network '%s' within datacenter '%s'): %s",
			network,
			datacenter,
			err)
	}

	return nil
}
