package net

type IcmpHdr struct {
	Type     IcmpType
	Code     uint8
	Checksum uint16
	Content  uint32
}

type IcmpType uint8

// ICMP Type Numbers
// https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types

var icmpTypes = map[IcmpType]string{
	0: "Echo Reply",
	// 1-2: Unassigned
	3: "Destination Unreachable",
	4: "Source Quench (Deprecated)",
	5: "Redirect",
	6: "Alternate Host Address (Deprecated)",
	// 7: Unassigned
	8:  "Echo",
	9:  "Router Advertisement",
	10: "Router Solicitation",
	11: "Time Exceeded",
	12: "Parameter Problem",
	13: "Timestamp",
	14: "Timestamp Reply",
	15: "Information Request (Deprecated)",
	16: "Information Reply (Deprecated)",
	17: "Address Mask Request (Deprecated)",
	18: "Address Mask Reply (Deprecated)",
	19: "Reserved (for Security)",
	// 20-29: Reserved (for Robustness Experiment)
	30: "Traceroute (Deprecated)",
	31: "Datagram Conversion Error (Deprecated)",
	32: "Mobile Host Redirect (Deprecated)",
	33: "IPv6 Where-Are-You (Deprecated)",
	34: "IPv6 I-Am-Here (Deprecated)",
	35: "Mobile Registration Request (Deprecated)",
	36: "Mobile Registration Reply (Deprecated)",
	37: "Domain Name Request (Deprecated)",
	38: "Domain Name Reply (Deprecated)",
	39: "SKIP (Deprecated)",
	40: "Photuris",
	41: "ICMP messages utilized by experimental mobility protocols such as Seamoby",
	42: "Extended Echo Request",
	43: "Extended Echo Reply",
	// 44-252: Unassigned
	253: "RFC3692-style Experiment 1",
	254: "RFC3692-style Experiment 2",
	// 255: Reserved
}
