package ev

import (
	"fmt"

	"github.com/fanyang89/easyvirt/proto"
	"github.com/libvirt/libvirt-go"
)

type EasyVirtServer struct {
	proto.UnimplementedEasyVirtServer

	conn *libvirt.Connect
}

func NewEasyVirtServer() (*EasyVirtServer, error) {
	conn, err := dial("")
	if err != nil {
		return nil, err
	}

	return &EasyVirtServer{
		conn: conn,
	}, nil
}

func dial(uri string) (*libvirt.Connect, error) {
	if uri == "" {
		uri = "qemu:///system"
	}

	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, fmt.Errorf("connect failed, %w", err)
	}
	return conn, nil
}

func (s *EasyVirtServer) Close() error {
	if s.conn != nil {
		_, err := s.conn.Close()
		if err != nil {
			return fmt.Errorf("close conn, %w", err)
		}
	}
	return nil
}

func (s *EasyVirtServer) checkConnected() error {
	if s.conn == nil {
		return fmt.Errorf("not connected")
	}
	return nil
}

type FindDomainParam struct {
	ID   *uint32
	Name *string
	UUID *string
}

func newPtrFromValue[T any](value T) *T {
	return &value
}

func NewFindDomainParam(finder *proto.DomainID) FindDomainParam {
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
