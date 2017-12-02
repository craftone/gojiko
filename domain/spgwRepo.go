package domain

import (
	"fmt"
	"net"
	"sync"
)

type spgwRepo struct {
	spgwCtrls map[string]SPgwIf
	mtx4Ctrl  sync.RWMutex
	spgwData  map[string]SPgwIf
	mtx4Data  sync.RWMutex
}

func newSPgwRepo() *spgwRepo {
	log.Info("Initialize S/PGW Repository")
	repo := &spgwRepo{
		spgwCtrls: make(map[string]SPgwIf),
		mtx4Ctrl:  sync.RWMutex{},
		spgwData:  make(map[string]SPgwIf),
		mtx4Data:  sync.RWMutex{},
	}

	return repo
}

func (s *spgwRepo) GetCtrl(addr net.UDPAddr) SPgwIf {
	s.mtx4Ctrl.RLock()
	defer s.mtx4Ctrl.RUnlock()
	if spgw, ok := s.spgwCtrls[addr.String()]; ok {
		return spgw
	}
	return nil
}

func (s *spgwRepo) AddCtrl(spgwCtrl SPgwIf) error {
	spgwData := spgwCtrl.getPair()
	if spgwData == nil {
		return fmt.Errorf("S/P-GW Ctrl must have a paired S/P-GW Data")
	}

	s.mtx4Ctrl.Lock()
	defer s.mtx4Ctrl.Unlock()
	ctrlAddr := spgwCtrl.UDPAddr()
	ctrlAddrKey := ctrlAddr.String()
	if _, ok := s.spgwCtrls[ctrlAddrKey]; ok {
		return fmt.Errorf("There is a S/P-GW Ctrl that has same UDP address")
	}

	s.mtx4Data.Lock()
	defer s.mtx4Data.Unlock()
	dataAddr := spgwData.UDPAddr()
	dataAddrKey := dataAddr.String()
	if _, ok := s.spgwData[dataAddrKey]; ok {
		return fmt.Errorf("There is a S/P-GW Data that has same UDP address")
	}

	s.spgwCtrls[ctrlAddrKey] = spgwCtrl
	s.spgwData[dataAddrKey] = spgwData
	return nil
}
