package hypervisor

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

func CreateNetworkProtocolProfile(ctx context.Context, clt *vim25.Client, datacenter, network string) error {
	dc, err := find.NewFinder(clt, true).Datacenter(ctx, datacenter)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine datacenter: %s", err)
	}

	object.NewCommon(clt, dc.Reference())
	request := &types.CreateIpPool{
		This: types.ManagedObjectReference{
			Type:  "",
			Value: "",
		},
		Pool: types.IpPool{
			Name: "Test",
		},
	}
	response, err := methods.CreateIpPool(ctx, clt.RoundTripper, request)
	if err != nil {
		log.Fatalf("Bar: %s", err)
	}
	spew.Dump(response)

	// task := object.NewTask(clt)
	// err = task.Wait(ctx)
	// if err != nil {
	// 	log.Fatalf("Foo: %s", err)
	// }

	return nil
}
