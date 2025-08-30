package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordArgon(t *testing.T) {
	pwd := NewAuthServiceArgon()
	hash, err := pwd.HashPassword("admin@123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}
func TestPasswordArgonVerifyPassword(t *testing.T) {
	pwd := NewAuthServiceArgon()
	hash, err := pwd.HashPassword("admin@123456")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
	ok, err := pwd.VerifyPassword(hash, "admin@123456")
	assert.NoError(t, err)
	assert.True(t, ok)

}
