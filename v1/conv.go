package ev

import (
	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

func convertDomain(d *libvirt.Domain) *proto.VirtualMachine {
	vm := &proto.VirtualMachine{}
	return vm
}
