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
