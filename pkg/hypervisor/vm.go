package hypervisor

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
)

type vAPPProperty struct {
	Key   string
	Value string
}

func SetVirtualMachineProperties(ctx context.Context, fnd *find.Finder, virtualmachine string, network string) error {
	vm, err := fnd.VirtualMachine(ctx, virtualmachine)
	if err != nil {
		return err
	}

	vappproperties := []vAPPProperty{
		{
			Key:   "guestinfo.dns.domains",
			Value: "${searchPath:%s}",
		},
		{
			Key:   "guestinfo.dns.servers",
			Value: "${dns:%s}",
		},
		{
			Key:   "guestinfo.interface.0.ip.0.address",
			Value: "${autoIp:%s}",
		},
		{
			Key:   "guestinfo.interface.0.ip.0.netmask",
			Value: "${netmask:%s}",
		},
		{
			Key:   "guestinfo.interface.0.route.0.gateway",
			Value: "${gateway:%s}",
		},
	}

	vappconfig := &types.VmConfigSpec{
		OvfEnvironmentTransport: []string{"com.vmware.guestInfo"},
	}

	for i, vappproperty := range vappproperties {
		vappconfig.Property = append(vappconfig.Property, types.VAppPropertySpec{
			ArrayUpdateSpec: types.ArrayUpdateSpec{
				Operation: types.ArrayUpdateOperationAdd,
			},
			Info: &types.VAppPropertyInfo{
				Key:          int32(i),
				Id:           vappproperty.Key,
				DefaultValue: fmt.Sprintf(vappproperty.Value, network),
				Type:         "expression",
			},
		})
	}

	// spew.Dump(vappconfig)

	task, err := vm.Reconfigure(ctx, types.VirtualMachineConfigSpec{
		VAppConfig: vappconfig,
	})
	if err != nil {
		return err
	}
	if err := task.Wait(ctx); err != nil {
		return err
	}

	return nil
}
