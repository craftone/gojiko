package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApnRepo_CreateFind(t *testing.T) {
	// Create first APN
	apn, err := TheApnRepo().Create("example.com")
	assert.Equal(t, "example.com", apn.Name())
	assert.NoError(t, err)

	// Create second APN
	apn2, err := TheApnRepo().Create("example2.com")
	assert.Equal(t, "example2.com", apn2.Name())
	assert.NoError(t, err)

	// Create dupulication error
	apnDup, err := TheApnRepo().Create("EXAMPLE.COM")
	assert.Equal(t, apn, apnDup)
	assert.Error(t, err)

	// Find success
	apnFinded := TheApnRepo().Find("example.com")
	assert.Equal(t, apn, apnFinded)

	// Find unsuccess
	apnFinded = TheApnRepo().Find("noexist.example.com")
	assert.Nil(t, apnFinded)

}
