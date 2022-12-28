package hypervisor

import (
	"context"
	"fmt"
	"log"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type vAPPProperty struct {
	Key   string
	Value string
}

func SetVirtualMachineProperties(ctx context.Context, clt *vim25.Client, datacenter, virtualmachine, network string) error {
	finder := find.NewFinder(clt, true)

	dc, err := finder.Datacenter(ctx, datacenter)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine datacenter: %s", err)
	}

	finder.SetDatacenter(dc)
	vm, err := finder.VirtualMachine(ctx, virtualmachine)
	if err != nil {
		log.Fatalf("[ERROR] Unable to determine virtual machine: %s", err)
	}

	var moref mo.VirtualMachine
	vm.Properties(ctx, vm.Reference(), []string{"config.vAppConfig"}, &moref)

	vappconfig := &types.VmConfigSpec{
		OvfEnvironmentTransport: []string{"com.vmware.guestInfo"},
	}

	currentvappproperties := moref.Config.VAppConfig.GetVmConfigInfo().Property
	for _, vappproperty := range currentvappproperties {
		vappconfig.Property = append(vappconfig.Property, types.VAppPropertySpec{
			ArrayUpdateSpec: types.ArrayUpdateSpec{
				Operation: types.ArrayUpdateOperationAdd,
			},
			Info: &types.VAppPropertyInfo{
				Key:              vappproperty.Key,
				ClassId:          vappproperty.ClassId,
				InstanceId:       vappproperty.InstanceId,
				Id:               vappproperty.Id,
				Category:         vappproperty.Category,
				Label:            vappproperty.Label,
				Type:             vappproperty.Type,
				TypeReference:    vappproperty.TypeReference,
				UserConfigurable: vappproperty.UserConfigurable,
				DefaultValue:     vappproperty.DefaultValue,
				Value:            vappproperty.Value,
				Description:      vappproperty.Description,
			},
		})
	}

	newvappproperties := []vAPPProperty{
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

	for i, vappproperty := range newvappproperties {
		vappconfig.Property = append(vappconfig.Property, types.VAppPropertySpec{
			ArrayUpdateSpec: types.ArrayUpdateSpec{
				Operation: types.ArrayUpdateOperationAdd,
			},
			Info: &types.VAppPropertyInfo{
				Key:          int32(i + len(currentvappproperties)),
				Id:           vappproperty.Key,
				DefaultValue: fmt.Sprintf(vappproperty.Value, network),
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
