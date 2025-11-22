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

func (s *EasyVirtServer) findDomain(p *proto.DomainID) (domain *libvirt.Domain, err error) {
	err = s.checkConnected()
	if err != nil {
		return
	}

	if p.Id != 0 {
		domain, err = s.conn.LookupDomainById(uint32(p.Id))
		if err == nil {
			return
		}
	}

	if p.Uuid != "" {
		domain, err = s.conn.LookupDomainByUUIDString(p.Uuid)
		if err == nil {
			return
		}
	}

	if p.Name != "" {
		domain, err = s.conn.LookupDomainByName(p.Name)
		if err == nil {
			return
		}
	}

	return nil, fmt.Errorf("not found")
}
