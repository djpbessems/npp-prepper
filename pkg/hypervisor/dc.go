package hypervisor

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

func CreateNetworkProtocolProfile(ctx context.Context, clt *vim25.Client, datacenter, network string) error {
	finder := find.NewFinder(clt, true)
	dc, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine datacenter: %s", err)
	}
	finder.SetDatacenter(dc)
	// nw, err := finder.Network(ctx, network)
	// if err != nil {
	// 	log.Fatalf("[ERROR] Unable to determine network: %s", err)
	// }

	request := &types.CreateIpPool{
		This: *clt.ServiceContent.IpPoolManager,
		Dc:   dc.Reference(),
		Pool: types.IpPool{
			Name:          "ippool-" + (uuid.New().String()),
			DnsDomain:     "meta.k8s.cluster",
			DnsSearchPath: "meta.k8s.cluster",
			// NetworkAssociation: []types.IpPoolAssociation{{
			// 	// This generates the error 'cannot use nw.Reference() (value of type types.ManagedObjectReference) as type *types.ManagedObjectReference in struct literal'
			// 	Network:     nw.Reference(),
			// 	NetworkName: network,
			// },
			// },
		},
	}

	if _, err := methods.CreateIpPool(ctx, clt.RoundTripper, request); err != nil {
		log.Fatalf("[ERROR]: Failed creating new network protocol profile (for network '%s' within datacenter '%s'): %s", network, datacenter, err)
	}

	return nil
}
