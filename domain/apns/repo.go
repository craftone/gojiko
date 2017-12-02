package apns

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// Repo is a Apn's repository. It is a singleton object and
// it will called from many GoRoutines.
type Repo struct {
	apnMap map[string]*Apn
	mtx    sync.RWMutex
}

func newRepo() *Repo {
	return &Repo{
		apnMap: map[string]*Apn{},
	}
}

func (r *Repo) Post(apn *Apn) error {
	log.Infof("Post new APN : %s", apn.fullString)
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.apnMap[apn.fullString]; ok {
		return fmt.Errorf("There is already the name's APN : %s", apn.fullString)
	}
	r.apnMap[apn.fullString] = apn
	return nil
}

func (a *Repo) Find(networkID, mcc, mnc string) (*Apn, error) {
	dummyApn, err := NewApn(networkID, mcc, mnc, []net.IP{net.IPv4(127, 0, 0, 1)})
	if err != nil {
		return nil, err
	}

	a.mtx.RLock()
	defer a.mtx.RUnlock()

	apn, ok := a.apnMap[dummyApn.fullString]
	if ok {
		return apn, nil
	}
	return nil, fmt.Errorf("There is no such APN : networkID=%s MCC=%s MNC=%s", networkID, mcc, mnc)
}

func (a *Repo) Delete(name string) error {
	return errors.New("Not implemented")
}
