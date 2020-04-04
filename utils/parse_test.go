package utils_test

import (
	"dinarchy/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStringNormal(t *testing.T) {
	str := "name,message,* * * * *"
	cr, err := utils.ParseCreate(str)
	assert.Equal(t, "name", cr.Name)
	assert.Equal(t, "* * * * *", cr.Cronstring)
	assert.Equal(t, "message", cr.Message)
	assert.Nil(t, err)
}

func TestParseStringWithSpaces(t *testing.T) {
	str := "name, message , * * * * *"
	cr, err := utils.ParseCreate(str)
	assert.Equal(t, "name", cr.Name)
	assert.Equal(t, "* * * * *", cr.Cronstring)
	assert.Equal(t, "message", cr.Message)
	assert.Nil(t, err)
}

func TestParseStringWithNameSpaces(t *testing.T) {
	str := "* * * * *, job name, message "
	_, err := utils.ParseCreate(str)
	assert.NotNil(t, err)
}
