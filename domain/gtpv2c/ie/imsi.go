package ie

type Imsi struct {
	tbcdIE
}

func NewImsi(instance byte, value string) (*Imsi, error) {
	tbcdIE, err := newTbcdIE(imsiNum, 0, instance, value, 6, 15)
	if err != nil {
		return nil, err
	}
	return &Imsi{tbcdIE}, nil
}

func (i *Imsi) Marshal() []byte {
	return i.tbcdIE.marshal()
}

func unmarshalImsi(h header, buf []byte) (*Imsi, error) {
	if h.typeNum != imsiNum {
		log.Panic("Invalid type")
	}

	s, err := unmarshalTbcd(buf)
	if err != nil {
		return nil, err
	}
	imsi, err := NewImsi(h.instance, s)
	if err != nil {
		return nil, err
	}
	return imsi, nil
}
