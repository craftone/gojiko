package domain

import "net"

type sgwCtrlRepo struct {
	*spgwRepo
}

func newSgwCtrlRepo() *sgwCtrlRepo {
	return &sgwCtrlRepo{newSPgwRepo()}
}

func (s *sgwCtrlRepo) getCtrl(addr net.UDPAddr) *SgwCtrl {
	sgwCtrl := s.spgwRepo.GetCtrl(addr)
	return sgwCtrl.(*SgwCtrl)
}
