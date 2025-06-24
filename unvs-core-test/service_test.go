package unvscoretest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "unvs.core"
	core "unvs.core"
)

func TestServiceHashPassword(t *testing.T) {
	// TODO: Write test cases here
	core.Config.LoadConfig("./")
	txt, err := core.Factory.GetPasswordService(t.Context()).HashPassword("root", "root")
	assert.NoError(t, err)
	assert.NotEmpty(t, txt)

}
func TestServiceHashPassword_InvalidPassword(t *testing.T) {
	// TODO: Write test cases here
	core.Config.LoadConfig("./")
	svc, err := core.Factory.GetTokenService(context.Background(), "en", "default")
	assert.NoError(t, err)
	assert.NotEmpty(t, svc)

}
