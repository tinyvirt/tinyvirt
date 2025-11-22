package ev

import (
	"fmt"

	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

type FindDomainParam struct {
	ID   *uint32
	Name *string
	UUID *string
}

func newPtrFromValue[T any](value T) *T {
	return &value
}

func NewFindDomainParam(finder *proto.FindDomain) FindDomainParam {
	return FindDomainParam{
		ID:   newPtrFromValue(finder.Id),
		Name: newPtrFromValue(finder.Name),
		UUID: newPtrFromValue(finder.Uuid),
	}
}

func (s *EasyVirtServer) findDomain(p FindDomainParam) (domain *libvirt.Domain, err error) {
	err = s.checkConnected()
	if err != nil {
		return
	}

	if p.ID != nil {
		domain, err = s.conn.LookupDomainById(*p.ID)
		if err == nil {
			return
		}
	}

	if p.UUID != nil {
		domain, err = s.conn.LookupDomainByUUIDString(*p.UUID)
		if err == nil {
			return
		}
	}

	if p.Name != nil {
		domain, err = s.conn.LookupDomainByName(*p.Name)
		if err == nil {
			return
		}
	}

	return nil, fmt.Errorf("not found")
}
