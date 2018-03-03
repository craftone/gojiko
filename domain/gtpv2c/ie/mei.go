package ie

// Mei struct expresses Mobile Equipment Identity (MEI) IE.
type Mei struct {
	tbcdIE
}

const (
	// MeiMinLen is minimum number of digits of a MEI
	MeiMinLen = 15
	// MeiMaxLen is maximum number of digits of a MEI
	MeiMaxLen = 16
)

func NewMei(instance byte, value string) (*Mei, error) {
	tbcdIE, err := newTbcdIE(meiNum, 0, instance, value, MeiMinLen, MeiMaxLen)
	if err != nil {
		return nil, err
	}
	return &Mei{tbcdIE}, nil
}

func (m *Mei) Marshal() []byte {
	return m.tbcdIE.marshal()
}

func unmarshalMei(h header, buf []byte) (*Mei, error) {
	if h.typeNum != meiNum {
		log.Panic("Invalid type")
	}

	s, err := unmarshalTbcd(buf)
	if err != nil {
		return nil, err
	}
	mei, err := NewMei(h.instance, s)
	if err != nil {
		return nil, err
	}
	return mei, nil
}
