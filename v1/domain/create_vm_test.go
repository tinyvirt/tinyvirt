package domain

import (
	"fmt"
	"testing"

	"github.com/fanyang89/easyvirt/v1"
	"github.com/stretchr/testify/assert"
	"libvirt.org/go/libvirtxml"
)

func domainFromPb(r *VM, c *ev.Config) *libvirtxml.Domain {
	bootDevices := make([]libvirtxml.DomainBootDevice, 0)
	for _, b := range r.BootDevices {
		bootDevices = append(bootDevices, libvirtxml.DomainBootDevice{Dev: b})
	}

	disks := make([]libvirtxml.DomainDisk, 0)
	for _, d := range r.DiskList {
		var disk libvirtxml.DomainDisk

		switch d.DeviceType {
		case DiskDeviceType_DISK_DEVICE_TYPE_CDROM:
			disk.Device = "cdrom"
		case DiskDeviceType_DISK_DEVICE_TYPE_DISK:
			disk.Device = "disk"
		default:
			panic(fmt.Sprintf("unexpected DiskDeviceType: %#v", d.DeviceType))
		}

		disk.Driver = &libvirtxml.DomainDiskDriver{Type: d.DriverType, Name: d.DeviceName}
		disk.Source = &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: d.FilePath}}
		disk.Target = &libvirtxml.DomainDiskTarget{Dev: d.DeviceName, Bus: d.Bus}

		if d.ReadOnly {
			disk.ReadOnly = &libvirtxml.DomainDiskReadOnly{}
		}

		disks = append(disks, disk)
	}

	return &libvirtxml.Domain{
		Type: "kvm",
		Name: r.DomainId.Name,
		UUID: "", // libvirt auto-fill on create
		VCPU: &libvirtxml.DomainVCPU{
			Value: uint(r.VcpuCount),
		},
		Memory: &libvirtxml.DomainMemory{
			Value: uint(r.Memory),
			Unit:  "MiB",
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm",
				Arch: "x86_64",
			},
			BootDevices: bootDevices,
		},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
			APIC: &libvirtxml.DomainFeatureAPIC{},
		},
		Devices: &libvirtxml.DomainDeviceList{
			Emulator: c.QemuPath,
			Disks:    disks,
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{
						AutoPort: "yes",
						Listen:   "0.0.0.0",
					},
				},
			},
		},
	}
}

func TestCreateVMDomainXML(t *testing.T) {
	vm := &VM{
		DomainId: &DomainID{
			Name: "test-vm",
		},
		VcpuCount:   4,
		Memory:      4096, // 4GiB
		BootDevices: []string{"hd", "cdrom", "network"},
		DiskList: []*DiskDescriptor{
			{
				DeviceType: DiskDeviceType_DISK_DEVICE_TYPE_DISK,
				DriverType: "qcow2",
				FilePath:   "tmp/sys.qcow2",
				DeviceName: "qemu",
				Bus:        "sata",
			},
			{
				DeviceType: DiskDeviceType_DISK_DEVICE_TYPE_CDROM,
				DriverType: "raw",
				ReadOnly:   true,
			},
		},
	}

	cfg := ev.NewConfig()
	d := domainFromPb(vm, cfg)
	b, err := d.Marshal()
	assert.NoError(t, err)
	_ = b
}
