package ev

import (
	"context"
	"fmt"

	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

var VOID = &proto.Void{}

func (s *EasyVirtServer) StartVM(ctx context.Context, req *proto.StartVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var domain *libvirt.Domain
	domain, err = s.findDomain(NewFindDomainParam(req.Domain))
	if err != nil {
		return nil, err
	}

	return VOID, domain.Create()
}

func (s *EasyVirtServer) StopVM(ctx context.Context, req *proto.StopVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var domain *libvirt.Domain
	domain, err = s.findDomain(NewFindDomainParam(req.Domain))
	if err != nil {
		return nil, err
	}

	if req.Force {
		err = domain.Destroy()
	} else {
		err = domain.Shutdown()
	}

	return VOID, err
}

func (s *EasyVirtServer) RestartVM(ctx context.Context, req *proto.RestartVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var domain *libvirt.Domain
	domain, err = s.findDomain(NewFindDomainParam(req.Domain))
	if err != nil {
		return nil, err
	}

	flags := libvirt.DOMAIN_REBOOT_DEFAULT

	if req.Force {
		flags |= libvirt.DOMAIN_REBOOT_SIGNAL
	}

	return VOID, domain.Reboot(flags)
}

func (s *EasyVirtServer) SuspendVM(ctx context.Context, req *proto.PauseVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var domain *libvirt.Domain
	domain, err = s.findDomain(NewFindDomainParam(req.Domain))
	if err != nil {
		return nil, err
	}

	return VOID, domain.Suspend()
}

func (s *EasyVirtServer) ResumeVM(ctx context.Context, req *proto.ResumeVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var domain *libvirt.Domain
	domain, err = s.findDomain(NewFindDomainParam(req.Domain))
	if err != nil {
		return nil, err
	}

	return VOID, domain.Resume()
}

func (s *EasyVirtServer) DeleteVM(ctx context.Context, req *proto.DeleteVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var domain *libvirt.Domain
	domain, err = s.findDomain(NewFindDomainParam(req.Domain))
	if err != nil {
		return nil, err
	}

	if req.DeleteDisks {
		return nil, fmt.Errorf("not impl")
	}

	return VOID, domain.Undefine()
}
