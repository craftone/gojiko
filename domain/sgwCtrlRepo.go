package domain

import "net"

type sgwCtrlRepo struct {
	*spgwRepo
}

func newSgwCtrlRepo() *sgwCtrlRepo {
	return &sgwCtrlRepo{newSPgwRepo()}
}

func (s *sgwCtrlRepo) GetSgwCtrl(addr net.UDPAddr) *SgwCtrl {
	sgwCtrl := s.spgwRepo.GetCtrl(addr)
	if sgwCtrl == nil {
		return nil
	}
	return sgwCtrl.(*SgwCtrl)
}
