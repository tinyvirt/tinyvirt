package disk

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/fanyang89/easyvirt/ent"
	"github.com/google/uuid"
)

type Manager struct {
	baseDir     string
	qemuImgPath string
	db          *ent.Client
}

func (s *Manager) getQemuDiskPath() string {
	if s.qemuImgPath != "" {
		return s.qemuImgPath
	}
	return "/usr/bin/qemu-img"
}

func (s *Manager) CreateDisk(ctx context.Context,
	diskName string, diskFormat DiskFormat, sizeGiB uint32, description string) (*uuid.UUID, error) {

	format := diskFormat.String()

	createdDisk, err := s.db.Disk.Create().
		SetName(diskName).
		SetFormat(format).
		SetDescription(description).
		SetSizeGB(sizeGiB).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// disk file name: <id>.<ext>
	diskFileName := createdDisk.ID.String() + "." + format
	diskFilePath := filepath.Join(s.baseDir, diskFileName)

	_, err = exec.LookPath(s.getQemuDiskPath())
	if err != nil {
		return nil, err
	}

	// qemu-img create Win2k.img 20G
	c := exec.CommandContext(ctx, s.getQemuDiskPath(), "create", "-f",
		format, diskFilePath, fmt.Sprintf("%dG", sizeGiB))
	b, err := c.CombinedOutput()
	slog.Debug("qemu-img create", "output", string(b), "err", err)
	if err != nil {
		return nil, err
	}

	return &createdDisk.ID, nil
}

func (s *Manager) DeleteDisk(ctx context.Context, id uuid.UUID) error {
	d, err := s.db.Disk.Get(ctx, id)
	if err != nil {
		return err
	}

	diskFileName := d.ID.String() + "." + d.Format
	diskFilePath := filepath.Join(s.baseDir, diskFileName)
	slog.Info("Deleting disk", "id", id, "path", diskFilePath)

	err1 := os.Remove(diskFilePath)
	err2 := s.db.Disk.DeleteOneID(id).Exec(ctx)
	return errors.CombineErrors(err1, err2)
}

func (s *Manager) GetDisk(ctx context.Context, id uuid.UUID) (*ent.Disk, error) {
	return s.db.Disk.Get(ctx, id)
}

func (s *Manager) ListDisks(ctx context.Context) ([]*ent.Disk, error) {
	return s.db.Disk.Query().All(ctx)
}

func (s *Manager) UpdateDisk(ctx context.Context, disk *ent.Disk) error {
	d, err := s.db.Disk.Get(ctx, disk.ID)
	if err != nil {
		return err
	}

	u := s.db.Disk.Update()

	if d.Name != disk.Name {
		u.SetName(disk.Name)
	}

	if d.Description != disk.Description {
		u.SetDescription(disk.Description)
	}

	var err1 error
	if d.SizeGB != disk.SizeGB {
		err1 = s.resizeDisk(ctx, d.ID.String()+"."+d.Format, disk.SizeGB)
		if err1 == nil {
			u.SetSizeGB(disk.SizeGB)
		}
	}

	err2 := u.Exec(ctx)
	return errors.CombineErrors(err1, err2)
}

func (s *Manager) resizeDisk(ctx context.Context, diskFilePath string, sizeGB uint32) error {
	_, err := exec.LookPath(s.getQemuDiskPath())
	if err != nil {
		return err
	}

	// qemu-img resize disk.img 30G
	c := exec.CommandContext(ctx, s.getQemuDiskPath(), "resize", diskFilePath, fmt.Sprintf("%dG", sizeGB))
	b, err := c.CombinedOutput()
	slog.Debug("qemu-img resize", "output", string(b), "err", err)
	return err
}
