package domain

type Apn struct {
	name string
}

func newApn(name string) *Apn {
	return &Apn{name}
}

func (a *Apn) Name() string {
	return a.name
}
