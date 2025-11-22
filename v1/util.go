package ev

import (
	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

func convertLibvirtState(state libvirt.DomainState) proto.VMState {
	switch state {
	case libvirt.DOMAIN_NOSTATE:
		return proto.VMState_VM_STATE_NOSTATE
	case libvirt.DOMAIN_RUNNING:
		return proto.VMState_VM_STATE_RUNNING
	case libvirt.DOMAIN_BLOCKED:
		return proto.VMState_VM_STATE_BLOCKED
	case libvirt.DOMAIN_PAUSED:
		return proto.VMState_VM_STATE_PAUSED
	case libvirt.DOMAIN_SHUTDOWN:
		return proto.VMState_VM_STATE_SHUTDOWN
	case libvirt.DOMAIN_SHUTOFF:
		return proto.VMState_VM_STATE_SHUTOFF
	case libvirt.DOMAIN_CRASHED:
		return proto.VMState_VM_STATE_CRASHED
	case libvirt.DOMAIN_PMSUSPENDED:
		return proto.VMState_VM_STATE_PMSUSPENDED
	default:
		return proto.VMState_VM_STATE_UNDEFINED
	}
}
