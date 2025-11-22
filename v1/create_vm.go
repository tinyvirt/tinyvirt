package ev

import (
	"context"
	"fmt"

	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

func (s *EasyVirtServer) createVM(domainXML string) (*libvirt.Domain, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	dom, err := s.conn.DomainDefineXML(domainXML)
	if err != nil {
		return nil, fmt.Errorf("create vm, %w", err)
	}
	return dom, nil
}

func (s *EasyVirtServer) CreateVM(ctx context.Context, req *proto.CreateVMRequest) (*proto.Void, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
