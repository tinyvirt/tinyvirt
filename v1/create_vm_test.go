package ev

import (
	"testing"

	"github.com/fanyang89/easyvirt/proto"
	"libvirt.org/go/libvirtxml"
)

func createVmXml(r *proto.VM) *libvirtxml.Domain {
	return &libvirtxml.Domain{
		Type:   "kvm",
		Name:   r.DomainId.Name,
		UUID:   "", // libvirt auto-fill on create
		VCPU:   &libvirtxml.DomainVCPU{Value: uint(r.VcpuCount)},
		Memory: &libvirtxml.DomainMemory{Value: r.MemoryMiB, Unit: "MiB"},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm",
				Arch: "x86_64",
			},
			BootDevices: []libvirtxml.DomainBootDevice{
				{Dev: "hd"},
			},
		},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
			APIC: &libvirtxml.DomainFeatureAPIC{},
		},
		Devices: &libvirtxml.DomainDeviceList{
			Emulator: "/usr/bin/qemu-system-x86_64",
			Disks: []libvirtxml.DomainDisk{
				{
					Device: "disk",
					Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2"},
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{
							File: "/var/lib/libvirt/images/demo.qcow2",
						},
					},
					Target: &libvirtxml.DomainDiskTarget{Dev: "vda", Bus: "virtio"},
				},
			},
			Interfaces: []libvirtxml.DomainInterface{
				{
					Source: &libvirtxml.DomainInterfaceSource{Bridge: &libvirtxml.DomainInterfaceSourceBridge{Bridge: "br0"}},
					Model:  &libvirtxml.DomainInterfaceModel{Type: "virtio"},
				},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{
						Port:     -1,
						AutoPort: "yes",
						Listen:   "0.0.0.0",
					},
				},
			},
		},
	}
}

func TestCreateVMDomainXML(t *testing.T) {

}
