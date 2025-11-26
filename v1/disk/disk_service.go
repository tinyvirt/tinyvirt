package disk

import (
	"context"

	"github.com/google/uuid"
	"github.com/tinyvirt/tinyvirt/ent"
	"google.golang.org/protobuf/proto"
)

var VOID = &Void{}

type Service struct {
	UnimplementedDiskManagerServer
	m *Manager
}

func NewService(db *ent.Client) *Service {
	return &Service{
		m: &Manager{
			db: db,
		},
	}
}

func (s *Service) CreateDisk(ctx context.Context, req *CreateDiskRequest) (*CreateDiskResponse, error) {
	d := proto.Clone(req.Disk).(*Disk)
	id, err := s.m.CreateDisk(ctx, d.DiskName, d.Format, d.Size, d.Description)
	if err != nil {
		return nil, err
	}
	d.DiskId = id.String()
	return &CreateDiskResponse{CreatedDisk: d}, nil
}

func (s *Service) DeleteDisk(ctx context.Context, req *DeleteDiskRequest) (*Void, error) {
	id, err := uuid.Parse(req.DiskId)
	if err != nil {
		return nil, err
	}
	return VOID, s.m.DeleteDisk(ctx, id)
}

func (s *Service) UpdateDisk(ctx context.Context, req *UpdateDiskRequest) (*Void, error) {
	d := req.Disk
	id, err := uuid.Parse(d.DiskId)
	if err != nil {
		return nil, err
	}

	return VOID, s.m.UpdateDisk(ctx, &ent.Disk{
		ID:          id,
		Name:        d.DiskName,
		Format:      d.Format.String(),
		Description: d.Description,
		SizeGB:      d.Size,
	})
}

func (s *Service) ListDisk(ctx context.Context, req *ListDiskRequest) (*ListDiskResponse, error) {
	disks, err := s.m.ListDisks(ctx)
	if err != nil {
		return nil, err
	}

	rsp := &ListDiskResponse{
		Disk: make([]*Disk, 0),
	}

	for _, d := range disks {
		rsp.Disk = append(rsp.Disk, &Disk{
			DiskId:      d.ID.String(),
			DiskName:    d.Name,
			Description: d.Description,
			Format:      DiskFormat(DiskFormat_value[d.Format]),
			Size:        d.SizeGB,
		})
	}

	return rsp, nil
}

func (s *Service) GetDisk(ctx context.Context, req *GetDiskRequest) (*GetDiskResponse, error) {
	id, err := uuid.Parse(req.DiskId)
	if err != nil {
		return nil, err
	}

	d, err := s.m.GetDisk(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetDiskResponse{
		Disk: &Disk{
			DiskId:      d.ID.String(),
			DiskName:    d.Name,
			Description: d.Description,
			Format:      DiskFormat(DiskFormat_value[d.Format]),
			Size:        d.SizeGB,
		},
	}, nil
}
