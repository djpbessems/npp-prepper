package utils

import "github.com/vmware/govmomi/vim25/types"

func BoolAddr(b bool) *bool {
	boolVar := b
	return &boolVar
}

func MoRefAddr(mo types.ManagedObjectReference) *types.ManagedObjectReference {
	morefVar := mo
	return &morefVar
}
