package gtpv2c

import "github.com/craftone/gojiko/gtpv2c/ie"

type CreateSessionResponse struct {
	header
	cause          *ie.Cause
	senderFteid    *ie.Fteid
	paa            *ie.Paa
	apnRestriction *ie.ApnRestriction
	apnAmbr        *ie.Ambr
	pco            *ie.PcoMsToNetwork
	bearerContextC *ie.BearerContextCreatedWithinCSRes
	recovery       *ie.Recovery
}
