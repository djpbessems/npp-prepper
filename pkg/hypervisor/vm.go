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

func SetVirtualMachineProperties(ctx context.Context, fnd *find.Finder, name string) error {
	vm, err := fnd.VirtualMachine(ctx, name)
	if err != nil {
		return err
	}

	vappproperties := []vAPPProperty{
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key:   "woot",
			Value: "dude",
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
				DefaultValue: fmt.Sprintf("${%v}", vappproperty.Value),
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
