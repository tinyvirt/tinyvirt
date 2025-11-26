package domain

import (
	"fmt"

	"github.com/tinyvirt/tinyvirt/ent"
	"libvirt.org/go/libvirt"
)

type EasyVirtService struct {
	UnimplementedTinyVirtServer

	conn *libvirt.Connect
	db   *ent.Client
}

func NewVirtService(uri string, db *ent.Client) (*EasyVirtService, error) {
	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, fmt.Errorf("connect failed, %w", err)
	}

	return &EasyVirtService{conn: conn, db: db}, nil
}

func (s *EasyVirtService) Close() error {
	if s.conn != nil {
		_, err := s.conn.Close()
		if err != nil {
			return fmt.Errorf("close conn, %w", err)
		}
	}
	return nil
}

func (s *EasyVirtService) findDomain(p *DomainID) (domain *libvirt.Domain, err error) {
	if p.Id != 0 {
		domain, err = s.conn.LookupDomainById(p.Id)
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
