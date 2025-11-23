package domain

import (
	"context"
	"fmt"

	"libvirt.org/go/libvirt"
)

func (s *EasyVirtService) createVM(domainXML string) (*libvirt.Domain, error) {
	dom, err := s.conn.DomainDefineXML(domainXML)
	if err != nil {
		return nil, fmt.Errorf("create vm, %w", err)
	}
	return dom, nil
}

func (s *EasyVirtService) CreateVM(ctx context.Context, req *CreateVMRequest) (*Void, error) {
	return nil, nil
}
