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
			Value: fmt.Sprintf("${searchPath:%s}", network),
		},
		{
			Key:   "guestinfo.dns.servers",
			Value: fmt.Sprintf("${dns:%s}", network),
		},
		{
			Key:   "guestinfo.interface.0.ip.0.address",
			Value: fmt.Sprintf("${autoIp:%s}", network),
		},
		{
			Key:   "guestinfo.interface.0.ip.0.netmask",
			Value: fmt.Sprintf("${netmask:%s}", network),
		},
		{
			Key:   "guestinfo.interface.0.route.0.gateway",
			Value: fmt.Sprintf("${gateway:%s}", network),
		},
	}

	vappconfig := &types.VmConfigSpec{
		// OvfEnvironmentTransport: []string{"com.vmware.guestinfo"},
	}

	for i, vappproperty := range vappproperties {
		vappconfig.Property = append(vappconfig.Property, types.VAppPropertySpec{
			ArrayUpdateSpec: types.ArrayUpdateSpec{
				Operation: types.ArrayUpdateOperationAdd,
			},
			Info: &types.VAppPropertyInfo{
				Key:          int32(i),
				Id:           vappproperty.Key,
				DefaultValue: vappproperty.Value,
				Type:         "expression",
			},
		})
	}

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
