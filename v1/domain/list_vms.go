package domain

import (
	"context"

	"libvirt.org/go/libvirt"
)

func (s *EasyVirtService) ListVMs(ctx context.Context, req *ListVMsRequest) (*ListVMsResponse, error) {
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

	vms := make([]*VM, 0, len(domains))
	for _, domain := range domains {
		vm := convertDomain(&domain)
		vms = append(vms, vm)
	}

	rsp := &ListVMsResponse{Vms: vms}
	return rsp, nil
}
