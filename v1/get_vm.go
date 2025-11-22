package ev

import (
	"context"

	"github.com/fanyang89/easyvirt/proto"
)

func (s *EasyVirtServer) GetVM(ctx context.Context, req *proto.GetVMRequest) (*proto.GetVMResponse, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	domain, err := s.findDomain(req.Domain)
	if err != nil {
		return nil, err
	}

	vm := convertDomain(domain)
	return &proto.GetVMResponse{Vm: vm}, nil
}
