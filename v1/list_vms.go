package ev

import (
	"context"

	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

func (s *EasyVirtServer) ListVMs(ctx context.Context, req *proto.ListVMsRequest) (*proto.ListVMsResponse, error) {
	err := s.checkConnected()
	if err != nil {
		return nil, err
	}

	var flags libvirt.ConnectListAllDomainsFlags
	if req.Inactive {
		flags = libvirt.CONNECT_LIST_DOMAINS_INACTIVE
	} else if req.Active {
		flags = libvirt.CONNECT_LIST_DOMAINS_ACTIVE
	} else if req.All {
		flags = 0
	}

	domains, err := s.conn.ListAllDomains(flags)
	if err != nil {
		return nil, err
	}

	vms := make([]*proto.VirtualMachine, 0, len(domains))
	for _, domain := range domains {
		vm := convertDomain(&domain)
		vms = append(vms, vm)
	}

	rsp := &proto.ListVMsResponse{Vms: vms}
	return rsp, nil
}
