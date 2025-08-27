package repo

import (
	"testing"
	"vdb"
	"xauth/config"

	"github.com/stretchr/testify/assert"
)

func TestCreateDefautUser(t *testing.T) {
	cfg, err := config.NewConfig("./../config/config.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	db, err := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	userRepo := NewUserRepoSql(db)

	assert.NotNil(t, userRepo)
	err = userRepo.CreateDefautUser()
	assert.NoError(t, err)

}
func BenchmarkCreateDefautUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cfg, err := config.NewConfig("./../config/config.yaml")
		assert.NoError(b, err)
		if err != nil {
			b.Fatal(err)
		}
		assert.NotNil(b, cfg)
		db, err := vdb.Open(cfg.Database.Driver, cfg.Database.Dsn)
		if err != nil {
			b.Fatal(err)
		}
		defer db.Close()
		assert.NoError(b, err)
		assert.NotNil(b, db)
		userRepo := NewUserRepoSql(db)

		assert.NotNil(b, userRepo)
		err = userRepo.CreateDefautUser()
		assert.NoError(b, err)
	}
}
