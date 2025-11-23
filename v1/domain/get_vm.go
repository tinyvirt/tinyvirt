package domain

import (
	"context"

	"libvirt.org/go/libvirt"
)

func convertDomain(d *libvirt.Domain) *VM {
	return nil
}

func (s *EasyVirtService) GetVM(ctx context.Context, req *GetVMRequest) (*GetVMResponse, error) {
	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	vm := convertDomain(domain)
	return &GetVMResponse{Vm: vm}, nil
}
