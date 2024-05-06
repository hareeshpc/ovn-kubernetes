//go:build linux
// +build linux

package util

import (
	"fmt"
	"net"

	"github.com/k8snetworkplumbingwg/govdpa/pkg/kvdpa"
	nadapi "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	"github.com/k8snetworkplumbingwg/sriovnet"
	"k8s.io/klog/v2"
)

type SriovnetOps interface {
	GetNetDevicesFromPci(pciAddress string) ([]string, error)
	GetNetDevicesFromAux(auxDev string) ([]string, error)
	GetPciFromNetDevice(name string) (string, error)
	GetUplinkRepresentor(vfPciAddress string) (string, error)
	GetUplinkRepresentorFromAux(auxDev string) (string, error)
	GetVfIndexByPciAddress(vfPciAddress string) (int, error)
	GetPfIndexByVfPciAddress(vfPciAddress string) (int, error)
	GetSfIndexByAuxDev(auxDev string) (int, error)
	GetVfRepresentor(uplink string, vfIndex int) (string, error)
	GetSfRepresentor(uplink string, sfIndex int) (string, error)
	GetPfPciFromVfPci(vfPciAddress string) (string, error)
	GetPfPciFromAux(auxDev string) (string, error)
	GetVfRepresentorDPU(pfID, vfIndex string) (string, error)
	GetRepresentorPeerMacAddress(netdev string) (net.HardwareAddr, error)
	GetRepresentorPortFlavour(netdev string) (sriovnet.PortFlavour, error)
	GetPortIndexFromRepresentor(name string) (int, error)
	GenDPDKPortParameters(netdev string) ([]string, error)
}

type defaultSriovnetOps struct {
}

var sriovnetOps SriovnetOps = &defaultSriovnetOps{}

// SetSriovnetOpsInst method would be used by unit tests in other packages
func SetSriovnetOpsInst(mockInst SriovnetOps) {
	sriovnetOps = mockInst
}

// GetSriovnetOps will be invoked by functions in other packages that would need access to the sriovnet library methods.
func GetSriovnetOps() SriovnetOps {
	return sriovnetOps
}

func (defaultSriovnetOps) GetNetDevicesFromPci(pciAddress string) ([]string, error) {
	return sriovnet.GetNetDevicesFromPci(pciAddress)
}

func (defaultSriovnetOps) GetNetDevicesFromAux(auxDev string) ([]string, error) {
	return sriovnet.GetNetDevicesFromAux(auxDev)
}

func (defaultSriovnetOps) GetPciFromNetDevice(name string) (string, error) {
	return sriovnet.GetPciFromNetDevice(name)
}

func (defaultSriovnetOps) GetUplinkRepresentor(vfPciAddress string) (string, error) {
	return sriovnet.GetUplinkRepresentor(vfPciAddress)
}

func (defaultSriovnetOps) GetUplinkRepresentorFromAux(auxDev string) (string, error) {
	return sriovnet.GetUplinkRepresentorFromAux(auxDev)
}

func (defaultSriovnetOps) GetVfIndexByPciAddress(vfPciAddress string) (int, error) {
	return sriovnet.GetVfIndexByPciAddress(vfPciAddress)
}

func (defaultSriovnetOps) GetPfIndexByVfPciAddress(vfPciAddress string) (int, error) {
	return sriovnet.GetPfIndexByVfPciAddress(vfPciAddress)
}

func (defaultSriovnetOps) GetSfIndexByAuxDev(auxDev string) (int, error) {
	return sriovnet.GetSfIndexByAuxDev(auxDev)
}

func (defaultSriovnetOps) GetVfRepresentor(uplink string, vfIndex int) (string, error) {
	return sriovnet.GetVfRepresentor(uplink, vfIndex)
}

func (defaultSriovnetOps) GetSfRepresentor(uplink string, sfIndex int) (string, error) {
	return sriovnet.GetSfRepresentor(uplink, sfIndex)
}

func (defaultSriovnetOps) GetPfPciFromVfPci(vfPciAddress string) (string, error) {
	return sriovnet.GetPfPciFromVfPci(vfPciAddress)
}

func (defaultSriovnetOps) GetPfPciFromAux(auxDev string) (string, error) {
	return sriovnet.GetPfPciFromAux(auxDev)
}

func (defaultSriovnetOps) GetVfRepresentorDPU(pfID, vfIndex string) (string, error) {
	return sriovnet.GetVfRepresentorDPU(pfID, vfIndex)
}

func (defaultSriovnetOps) GetRepresentorPeerMacAddress(netdev string) (net.HardwareAddr, error) {
	return sriovnet.GetRepresentorPeerMacAddress(netdev)
}

func (defaultSriovnetOps) GetRepresentorPortFlavour(netdev string) (sriovnet.PortFlavour, error) {
	return sriovnet.GetRepresentorPortFlavour(netdev)
}

func (defaultSriovnetOps) GetPortIndexFromRepresentor(name string) (int, error) {
	return sriovnet.GetPortIndexFromRepresentor(name)
}

func (defaultSriovnetOps) GenDPDKPortParameters(netdev string) ([]string, error) {
	return GenDPDKPortParameters(netdev)
}

// GenDPDKPortParameters generates the parameters required by
// DPDK to use the netdevice which name was provided.
// Only physical port and PF, VF, SF representors are supported.
func GenDPDKPortParameters(netdev string) ([]string, error) {
	var pci string
	var portID int

	flavor, err := GetSriovnetOps().GetRepresentorPortFlavour(netdev)
	if err != nil {
		return nil, fmt.Errorf("failure to find port %v type: %v", netdev, err)
	}

	switch flavor {
	case sriovnet.PORT_FLAVOUR_PHYSICAL:
	case sriovnet.PORT_FLAVOUR_PCI_PF:
	case sriovnet.PORT_FLAVOUR_PCI_VF:
	case sriovnet.PORT_FLAVOUR_PCI_SF:
	default:
		return nil, fmt.Errorf("unknown port %s type %v", netdev, flavor)
	}

	pci, err = GetSriovnetOps().GetPciFromNetDevice(netdev)
	if err != nil {
		return nil, fmt.Errorf("failed to get PCI address of port %v: %v", netdev, err)
	}

	if flavor == sriovnet.PORT_FLAVOUR_PCI_VF ||
		flavor == sriovnet.PORT_FLAVOUR_PCI_SF {
		portID, err = GetSriovnetOps().GetPortIndexFromRepresentor(netdev)
		if err != nil {
			return nil, fmt.Errorf("failed to get port %s index: %v", netdev, err)
		}
	}

	typeToDevarg := map[sriovnet.PortFlavour]string{
		sriovnet.PORT_FLAVOUR_PHYSICAL: fmt.Sprintf("%v", pci),
		// the HPF representor port ID is always -1 in DPDK.
		sriovnet.PORT_FLAVOUR_PCI_PF: fmt.Sprintf("%v,representor=[-1]", pci),
		sriovnet.PORT_FLAVOUR_PCI_VF: fmt.Sprintf("%v,representor=vf[%v]", pci, portID),
		sriovnet.PORT_FLAVOUR_PCI_SF: fmt.Sprintf("%v,representor=sf[%v]", pci, portID),
	}

	return []string{
		"type=dpdk",
		fmt.Sprintf("options:dpdk-devargs=\"%v,dv_xmeta_en=4,dv_flow_en=2\"", typeToDevarg[flavor]),
	}, nil
}

// GetFunctionRepresentorName returns representor name for passed device ID. Supported devices are Virtual Function
// or Scalable Function
func GetFunctionRepresentorName(deviceID string) (string, error) {
	var rep, uplink string
	var err error
	var index int

	if IsPCIDeviceName(deviceID) { // PCI device
		uplink, err = GetSriovnetOps().GetUplinkRepresentor(deviceID)
		if err != nil {
			return "", err
		}
		index, err = GetSriovnetOps().GetVfIndexByPciAddress(deviceID)
		if err != nil {
			return "", err
		}
		rep, err = GetSriovnetOps().GetVfRepresentor(uplink, index)
	} else if IsAuxDeviceName(deviceID) { // Auxiliary device
		uplink, err = GetSriovnetOps().GetUplinkRepresentorFromAux(deviceID)
		if err != nil {
			return "", err
		}
		index, err = GetSriovnetOps().GetSfIndexByAuxDev(deviceID)
		if err != nil {
			return "", err
		}
		rep, err = GetSriovnetOps().GetSfRepresentor(uplink, index)
	} else {
		return "", fmt.Errorf("cannot determine device type for id '%s'", deviceID)
	}
	if err != nil {
		return "", err
	}
	return rep, nil
}

// GetNetdevNameFromDeviceId returns the netdevice name from the passed device ID.
func GetNetdevNameFromDeviceId(deviceId string, deviceInfo nadapi.DeviceInfo) (string, error) {
	var netdevices []string
	var err error

	if IsPCIDeviceName(deviceId) {
		if deviceInfo.Vdpa != nil {
			if deviceInfo.Vdpa.Driver == "vhost" {
				klog.V(2).Info("deviceInfo.Vdpa.Driver is vhost, returning empty netdev")
				return "", nil
			}
		}

		// If a virtio/vDPA device exists, it takes preference over the vendor device, steering-wize
		var vdpaDevice kvdpa.VdpaDevice
		vdpaDevice, err = GetVdpaOps().GetVdpaDeviceByPci(deviceId)
		if err == nil && vdpaDevice != nil && vdpaDevice.Driver() == kvdpa.VirtioVdpaDriver {
			klog.V(2).Infof("deviceInfo.Vdpa.Driver is virtio, returning netdev %s", vdpaDevice.VirtioNet().NetDev())
			return vdpaDevice.VirtioNet().NetDev(), nil
		}
		if err != nil {
			klog.Warningf("Error when searching for the virtio/vdpa netdev: %v", err)
		}

		netdevices, err = GetSriovnetOps().GetNetDevicesFromPci(deviceId)
	} else { // Auxiliary network device
		netdevices, err = GetSriovnetOps().GetNetDevicesFromAux(deviceId)
	}
	if err != nil {
		return "", err
	}

	// Make sure we have 1 netdevice per pci address
	numNetDevices := len(netdevices)
	if numNetDevices != 1 {
		return "", fmt.Errorf("failed to get one netdevice interface (count %d) per Device ID %s", numNetDevices, deviceId)
	}
	return netdevices[0], nil
}
