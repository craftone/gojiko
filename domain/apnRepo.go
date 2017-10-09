package domain

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// ApnRepo is a Apn's repository. It is a singleton object and
// it will called from many GoRoutines.
type ApnRepo struct {
	apnMap map[string]*Apn
	mtx    sync.RWMutex
}

// TheApnRepo returns the global APN repository in this program.
func TheApnRepo() *ApnRepo {
	return theApnRepo
}

// singleton object
var theApnRepo = newApnRepo()

func newApnRepo() *ApnRepo {
	return &ApnRepo{
		apnMap: map[string]*Apn{},
	}
}

func (a *ApnRepo) Create(name string) (*Apn, error) {
	rname := strings.ToLower(name)

	a.mtx.Lock()
	defer a.mtx.Unlock()

	if apn, ok := a.apnMap[rname]; ok {
		return apn, fmt.Errorf("There is already the name's APN : %s", name)
	}

	apn := newApn(rname)
	a.apnMap[rname] = apn
	return apn, nil
}

func (a *ApnRepo) Find(name string) *Apn {
	rname := strings.ToLower(name)

	a.mtx.RLock()
	defer a.mtx.RUnlock()

	apn, ok := a.apnMap[rname]
	if ok {
		return apn
	}
	return nil
}

func (a *ApnRepo) Delete(name string) error {
	return errors.New("Not implemented")
}
