package pco

// IPAllocViaNasSignalling represents a PCO's container
// "IP address allocation via NAS signalling".
type IPAllocViaNasSignalling struct {
	header
}

func NewIPAllocViaNasSignalling() *IPAllocViaNasSignalling {
	return &IPAllocViaNasSignalling{
		header{ipAllocViaNasSigNum, 0},
	}
}

func (i *IPAllocViaNasSignalling) marshal() []byte {
	return i.header.marshal([]byte{})
}
