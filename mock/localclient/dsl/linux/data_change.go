package linux

import (
	"go.ligato.io/vpp-agent/v3/clientv2/linux"
	"go.ligato.io/vpp-agent/v3/clientv2/vpp"

	"github.com/americanbinary/vpp/mock/localclient/dsl"

	"go.ligato.io/vpp-agent/v3/pkg/models"

	"go.ligato.io/vpp-agent/v3/proto/ligato/linux/interfaces"
	linux_iptables "go.ligato.io/vpp-agent/v3/proto/ligato/linux/iptables"
	"go.ligato.io/vpp-agent/v3/proto/ligato/linux/l3"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/abf"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/acl"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/ipsec"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l2"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/l3"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/nat"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/punt"
	"go.ligato.io/vpp-agent/v3/proto/ligato/vpp/stn"
)

// MockDataChangeDSL is mock for DataChangeDSL.
type MockDataChangeDSL struct {
	dsl.CommonMockDSL
}

// NewMockDataChangeDSL is a constructor for MockDataChangeDSL.
func NewMockDataChangeDSL(commitFunc dsl.CommitFunc) *MockDataChangeDSL {
	return &MockDataChangeDSL{CommonMockDSL: dsl.NewCommonMockDSL(commitFunc)}
}

// MockPutDSL is a mock for PutDSL.
type MockPutDSL struct {
	parent *MockDataChangeDSL
}

// MockDeleteDSL is a mock for DeleteDSL.
type MockDeleteDSL struct {
	parent *MockDataChangeDSL
}

// Put initiates a chained sequence of data change DSL statements declaring
// new or changing existing configurable objects.
func (d *MockDataChangeDSL) Put() linuxclient.PutDSL {
	return &MockPutDSL{d}
}

// Delete initiates a chained sequence of data change DSL statements
// removing existing configurable objects.
func (d *MockDataChangeDSL) Delete() linuxclient.DeleteDSL {
	return &MockDeleteDSL{d}
}

// Send commits the transaction into the mock DB.
func (d *MockDataChangeDSL) Send() vppclient.Reply {
	err := d.CommitFunc(d.Values)
	return &dsl.Reply{Err: err}
}

// VppInterface adds a mock request to create or update VPP network interface.
func (d *MockPutDSL) VppInterface(val *vpp_interfaces.Interface) linuxclient.PutDSL {
	key := vpp_interfaces.InterfaceKey(val.Name)
	d.parent.Values[key] = val
	return d
}

// ABF adds a request to create or update VPP ACL-based forwarding.
func (d *MockPutDSL) ABF(val *vpp_abf.ABF) linuxclient.PutDSL {
	key := vpp_abf.Key(val.Index)
	d.parent.Values[key] = val
	return d
}

// BD adds a mock request to create or update VPP Bridge Domain.
func (d *MockPutDSL) BD(val *vpp_l2.BridgeDomain) linuxclient.PutDSL {
	key := vpp_l2.BridgeDomainKey(val.Name)
	d.parent.Values[key] = val
	return d
}

// BDFIB adds a mock request to create or update VPP L2 Forwarding Information
// Base.
func (d *MockPutDSL) BDFIB(val *vpp_l2.FIBEntry) linuxclient.PutDSL {
	key := vpp_l2.FIBKey(val.BridgeDomain, val.PhysAddress)
	d.parent.Values[key] = val
	return d
}

// XConnect adds a mock request to create or update VPP Cross Connect.
func (d *MockPutDSL) XConnect(val *vpp_l2.XConnectPair) linuxclient.PutDSL {
	key := vpp_l2.XConnectKey(val.ReceiveInterface)
	d.parent.Values[key] = val
	return d
}

// StaticRoute adds a mock request to create or update VPP L3 Static Route.
func (d *MockPutDSL) StaticRoute(val *vpp_l3.Route) linuxclient.PutDSL {
	key := vpp_l3.RouteKey(val.OutgoingInterface, val.VrfId, val.DstNetwork, val.NextHopAddr)
	d.parent.Values[key] = val
	return d
}

// Span adds a mock VPP span to the Put request.
func (d *MockPutDSL) Span(val *vpp_interfaces.Span) linuxclient.PutDSL {
	key := vpp_interfaces.SpanKey(val.InterfaceFrom, val.InterfaceTo)
	d.parent.Values[key] = val
	return d
}

// ACL adds a mock request to create or update VPP Access Control List.
func (d *MockPutDSL) ACL(val *vpp_acl.ACL) linuxclient.PutDSL {
	key := vpp_acl.Key(val.Name)
	d.parent.Values[key] = val
	return d
}

// Arp adds a request to create or update VPP L3 ARP.
func (d *MockPutDSL) Arp(val *vpp_l3.ARPEntry) linuxclient.PutDSL {
	key := vpp_l3.ArpEntryKey(val.Interface, val.IpAddress)
	d.parent.Values[key] = val
	return d
}

// ProxyArp adds a request to create or update VPP L3 proxy ARP.
func (d *MockPutDSL) ProxyArp(val *vpp_l3.ProxyARP) linuxclient.PutDSL {
	key := vpp_l3.ProxyARPKey()
	d.parent.Values[key] = val
	return d
}

// IPScanNeighbor adds L3 IP Scan Neighbor to the RESYNC request.
func (d *MockPutDSL) IPScanNeighbor(val *vpp_l3.IPScanNeighbor) linuxclient.PutDSL {
	key := vpp_l3.IPScanNeighborKey()
	d.parent.Values[key] = val
	return d
}

// StnRule adds a request to create or update VPP Stn rule.
func (d *MockPutDSL) StnRule(val *vpp_stn.Rule) linuxclient.PutDSL {
	key := vpp_stn.Key(val.Interface, val.IpAddress)
	d.parent.Values[key] = val
	return d
}

// NAT44Global adds a request to set global configuration for NAT44
func (d *MockPutDSL) NAT44Global(val *vpp_nat.Nat44Global) linuxclient.PutDSL {
	key := vpp_nat.GlobalNAT44Key()
	d.parent.Values[key] = val
	return d
}

// DNAT44 adds a request to create a new DNAT configuration
func (d *MockPutDSL) DNAT44(val *vpp_nat.DNat44) linuxclient.PutDSL {
	key := vpp_nat.DNAT44Key(val.Label)
	d.parent.Values[key] = val
	return d
}

// LinuxInterface adds a mock request to create or update Linux network interface.
func (d *MockPutDSL) LinuxInterface(val *linux_interfaces.Interface) linuxclient.PutDSL {
	key := linux_interfaces.InterfaceKey(val.Name)
	d.parent.Values[key] = val
	return d
}

func (d *MockPutDSL) LinuxArpEntry(val *linux_l3.ARPEntry) linuxclient.PutDSL {
	key := linux_l3.ArpKey(val.Interface, val.IpAddress)
	d.parent.Values[key] = val
	return d
}

func (d *MockPutDSL) LinuxRoute(val *linux_l3.Route) linuxclient.PutDSL {
	key := linux_l3.RouteKey(val.DstNetwork, val.OutgoingInterface)
	d.parent.Values[key] = val
	return d
}

// IPSecSA adds request to create a new Security Association
func (d *MockPutDSL) IPSecSA(val *vpp_ipsec.SecurityAssociation) linuxclient.PutDSL {
	key := vpp_ipsec.SAKey(val.Index)
	d.parent.Values[key] = val
	return d
}

// IPSecSPD adds request to create a new Security Policy Database
func (d *MockPutDSL) IPSecSPD(val *vpp_ipsec.SecurityPolicyDatabase) linuxclient.PutDSL {
	key := vpp_ipsec.SPDKey(val.Index)
	d.parent.Values[key] = val
	return d
}

// IPSecTunnelProtection adds request to create a new IPSec tunnel protection
func (d *MockPutDSL) IPSecTunnelProtection(val *vpp_ipsec.TunnelProtection) linuxclient.PutDSL {
	key := models.Key(val)
	d.parent.Values[key] = val
	return d
}

// PuntIPRedirect adds request to create or update rule to punt L3 traffic via interface.
func (d *MockPutDSL) PuntIPRedirect(val *vpp_punt.IPRedirect) linuxclient.PutDSL {
	key := vpp_punt.IPRedirectKey(val.L3Protocol, val.TxInterface)
	d.parent.Values[key] = val
	return d
}

// PuntToHost adds request to create or update rule to punt L4 traffic to a host.
func (d *MockPutDSL) PuntToHost(val *vpp_punt.ToHost) linuxclient.PutDSL {
	key := vpp_punt.ToHostKey(val.L3Protocol, val.L4Protocol, val.Port)
	d.parent.Values[key] = val
	return d
}

// PuntException adds request to delete a rule used to add exception in punting of traffic to a host.
func (d *MockPutDSL) PuntException(val *vpp_punt.Exception) linuxclient.PutDSL {
	key := vpp_punt.ExceptionKey(val.Reason)
	d.parent.Values[key] = val
	return d
}

// VrfTable adds a request to create or update VPP VRF table.
func (d *MockPutDSL) VrfTable(val *vpp_l3.VrfTable) linuxclient.PutDSL {
	key := vpp_l3.VrfTableKey(val.Id, val.Protocol)
	d.parent.Values[key] = val
	return d
}

// NAT44AddressPool adds a request to create or update NAT44 address pool.
func (d *MockPutDSL) NAT44AddressPool(pool *vpp_nat.Nat44AddressPool) linuxclient.PutDSL {
	key := vpp_nat.Nat44AddressPoolKey(pool.VrfId, pool.FirstIp, pool.LastIp)
	d.parent.Values[key] = pool
	return d
}

// NAT44Interface adds a request to create or update NAT44 interface configuration.
func (d *MockPutDSL) NAT44Interface(natif *vpp_nat.Nat44Interface) linuxclient.PutDSL {
	key := vpp_nat.Nat44InterfaceKey(natif.Name)
	d.parent.Values[key] = natif
	return d
}

// IptablesRuleChain adds request to create or update iptables rule chain.
func (d *MockPutDSL) IptablesRuleChain(val *linux_iptables.RuleChain) linuxclient.PutDSL {
	key := linux_iptables.RuleChainKey(val.Name)
	d.parent.Values[key] = val
	return d
}

// Delete changes the DSL mode to allow removal of an existing configuration.
func (d *MockPutDSL) Delete() linuxclient.DeleteDSL {
	return &MockDeleteDSL{d.parent}
}

// Send commits the transaction into the mock DB.
func (d *MockPutDSL) Send() vppclient.Reply {
	return d.parent.Send()
}

// VppInterface adds a mock request to delete an existing VPP network interface.
func (d *MockDeleteDSL) VppInterface(interfaceName string) linuxclient.DeleteDSL {
	key := vpp_interfaces.InterfaceKey(interfaceName)
	d.parent.Values[key] = nil
	return d
}

// IptablesRuleChain adds request to delete iptables rule chain.
func (d *MockDeleteDSL) IptablesRuleChain(name string) linuxclient.DeleteDSL {
	key := linux_iptables.RuleChainKey(name)
	d.parent.Values[key] = nil
	return d
}

// NAT44AddressPool adds a request to delete NAT44 address pool.
func (d *MockDeleteDSL) NAT44AddressPool(pool *vpp_nat.Nat44AddressPool) linuxclient.DeleteDSL {
	key := vpp_nat.Nat44AddressPoolKey(pool.VrfId, pool.FirstIp, pool.LastIp)
	d.parent.Values[key] = nil
	return d
}

// NAT44Interface adds a request to delete NAT44 interface configuration.
func (d *MockDeleteDSL) NAT44Interface(natIf *vpp_nat.Nat44Interface) linuxclient.DeleteDSL {
	key := vpp_nat.Nat44InterfaceKey(natIf.Name)
	d.parent.Values[key] = nil
	return d
}

// ABF adds a request to delete and existing VPP Access Control List.
func (d *MockDeleteDSL) ABF(abfIndex uint32) linuxclient.DeleteDSL {
	key := vpp_abf.Key(abfIndex)
	d.parent.Values[key] = nil
	return d
}

// BD adds a mock request to delete an existing VPP Bridge Domain.
func (d *MockDeleteDSL) BD(bdName string) linuxclient.DeleteDSL {
	key := vpp_l2.BridgeDomainKey(bdName)
	d.parent.Values[key] = nil
	return d
}

// BDFIB adds a mock request to delete an existing VPP L2 Forwarding Information
// Base.
func (d *MockDeleteDSL) BDFIB(bdName string, mac string) linuxclient.DeleteDSL {
	key := vpp_l2.FIBKey(bdName, mac)
	d.parent.Values[key] = nil
	return d
}

// XConnect adds a mock request to delete an existing VPP Cross Connect.
func (d *MockDeleteDSL) XConnect(rxIfName string) linuxclient.DeleteDSL {
	key := vpp_l2.XConnectKey(rxIfName)
	d.parent.Values[key] = nil
	return d
}

// StaticRoute adds a mock request to delete an existing VPP L3 Static Route..
func (d *MockDeleteDSL) StaticRoute(iface string, vrf uint32, dstAddr string, nextHopAddr string) linuxclient.DeleteDSL {
	key := vpp_l3.RouteKey(iface, vrf, dstAddr, nextHopAddr)
	d.parent.Values[key] = nil
	return d
}

// Span adds a mock VPP span to the Delete request.
func (d *MockDeleteDSL) Span(span *vpp_interfaces.Span) linuxclient.DeleteDSL {
	key := vpp_interfaces.SpanKey(span.InterfaceFrom, span.InterfaceTo)
	d.parent.Values[key] = nil
	return d
}

// ACL adds a mock request to delete an existing VPP Access Control List.
func (d *MockDeleteDSL) ACL(aclName string) linuxclient.DeleteDSL {
	key := vpp_acl.Key(aclName)
	d.parent.Values[key] = nil
	return d
}

// Arp adds a request to delete an existing VPP L3 ARP entry.
func (d *MockDeleteDSL) Arp(ifaceName string, ipAddr string) linuxclient.DeleteDSL {
	key := vpp_l3.ArpEntryKey(ifaceName, ipAddr)
	d.parent.Values[key] = nil
	return d
}

// ProxyArp adds a request to delete an existing VPP L3 proxy ARP.
func (d *MockDeleteDSL) ProxyArp() linuxclient.DeleteDSL {
	key := vpp_l3.ProxyARPKey()
	d.parent.Values[key] = nil
	return d
}

// IPScanNeighbor adds a request to delete an existing VPP L3 IP Scan Neighbor.
func (d *MockDeleteDSL) IPScanNeighbor() linuxclient.DeleteDSL {
	key := vpp_l3.IPScanNeighborKey()
	d.parent.Values[key] = nil
	return d
}

// StnRule adds request to delete Stn rule.
func (d *MockDeleteDSL) StnRule(ifName, ipAddr string) linuxclient.DeleteDSL {
	key := vpp_stn.Key(ifName, ipAddr)
	d.parent.Values[key] = nil
	return d
}

// NAT44Global adds a request to remove global configuration for NAT44
func (d *MockDeleteDSL) NAT44Global() linuxclient.DeleteDSL {
	key := vpp_nat.GlobalNAT44Key()
	d.parent.Values[key] = nil
	return d
}

// DNAT44 adds a request to delete a DNAT configuration identified by label
func (d *MockDeleteDSL) DNAT44(label string) linuxclient.DeleteDSL {
	key := vpp_nat.DNAT44Key(label)
	d.parent.Values[key] = nil
	return d
}

// LinuxInterface adds a mock request to delete an existing Linux network
// interface.
func (d *MockDeleteDSL) LinuxInterface(ifName string) linuxclient.DeleteDSL {
	key := linux_interfaces.InterfaceKey(ifName)
	d.parent.Values[key] = nil
	return d
}

func (d *MockDeleteDSL) LinuxArpEntry(ifaceName string, ipAddr string) linuxclient.DeleteDSL {
	key := linux_l3.ArpKey(ifaceName, ipAddr)
	d.parent.Values[key] = nil
	return d
}

func (d *MockDeleteDSL) LinuxRoute(dstAddr, outIfaceName string) linuxclient.DeleteDSL {
	key := linux_l3.RouteKey(dstAddr, outIfaceName)
	d.parent.Values[key] = nil
	return d
}

// IPSecSA adds request to delete a Security Association
func (d *MockDeleteDSL) IPSecSA(saIndex uint32) linuxclient.DeleteDSL {
	key := vpp_ipsec.SAKey(saIndex)
	d.parent.Values[key] = nil
	return d
}

// IPSecSPD adds request to delete a Security Policy Database
func (d *MockDeleteDSL) IPSecSPD(spdIndex uint32) linuxclient.DeleteDSL {
	key := vpp_ipsec.SPDKey(spdIndex)
	d.parent.Values[key] = nil
	return d
}

// IPSecTunnelProtection adds request to delete an IPSec tunnel protection from an interface
func (d *MockDeleteDSL) IPSecTunnelProtection(tp *vpp_ipsec.TunnelProtection) linuxclient.DeleteDSL {
	key := models.Key(tp)
	d.parent.Values[key] = nil
	return d
}

// PuntIPRedirect adds request to delete a rule used to punt L3 traffic via interface.
func (d *MockDeleteDSL) PuntIPRedirect(l3Proto vpp_punt.L3Protocol, txInterface string) linuxclient.DeleteDSL {
	key := vpp_punt.IPRedirectKey(l3Proto, txInterface)
	d.parent.Values[key] = nil
	return d
}

// PuntToHost adds request to delete a rule used to punt L4 traffic to a host.
func (d *MockDeleteDSL) PuntToHost(l3Proto vpp_punt.L3Protocol, l4Proto vpp_punt.L4Protocol, port uint32) linuxclient.DeleteDSL {
	key := vpp_punt.ToHostKey(l3Proto, l4Proto, port)
	d.parent.Values[key] = nil
	return d
}

// PuntException adds request to delete a rule used to add exception in punting of traffic to a host.
func (d *MockDeleteDSL) PuntException(reason string) linuxclient.DeleteDSL {
	key := vpp_punt.ExceptionKey(reason)
	d.parent.Values[key] = nil
	return d
}

// VrfTable adds a request to delete existing VPP VRF table.
func (d *MockDeleteDSL) VrfTable(id uint32, proto vpp_l3.VrfTable_Protocol) linuxclient.DeleteDSL {
	key := vpp_l3.VrfTableKey(id, proto)
	d.parent.Values[key] = nil
	return d
}

// Put changes the DSL mode to allow configuration editing.
func (d *MockDeleteDSL) Put() linuxclient.PutDSL {
	return &MockPutDSL{d.parent}
}

// Send commits the transaction into the mock DB.
func (d *MockDeleteDSL) Send() vppclient.Reply {
	return d.parent.Send()
}
