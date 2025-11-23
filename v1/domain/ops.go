package domain

import (
	"context"
	"fmt"

	"libvirt.org/go/libvirt"
)

var VOID = &Void{}

func (s *EasyVirtService) StartVM(ctx context.Context, req *StartVMRequest) (*Void, error) {
	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	return VOID, domain.Create()
}

func (s *EasyVirtService) StopVM(ctx context.Context, req *StopVMRequest) (*Void, error) {
	domain, err := s.findDomain(req.Domain)
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

func (s *EasyVirtService) RestartVM(ctx context.Context, req *RestartVMRequest) (*Void, error) {
	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	flags := libvirt.DOMAIN_REBOOT_DEFAULT

	if req.Force {
		flags |= libvirt.DOMAIN_REBOOT_SIGNAL
	}

	return VOID, domain.Reboot(flags)
}

func (s *EasyVirtService) SuspendVM(ctx context.Context, req *PauseVMRequest) (*Void, error) {
	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	return VOID, domain.Suspend()
}

func (s *EasyVirtService) ResumeVM(ctx context.Context, req *ResumeVMRequest) (*Void, error) {
	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	return VOID, domain.Resume()
}

func (s *EasyVirtService) DeleteVM(ctx context.Context, req *DeleteVMRequest) (*Void, error) {
	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	if req.DeleteDisks {
		return nil, fmt.Errorf("not impl")
	}

	return VOID, domain.Undefine()
}
