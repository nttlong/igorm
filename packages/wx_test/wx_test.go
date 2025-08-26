package wxtest

import (
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type Config interface {
}
type DbConfig struct {
	wx.Provider[DbConfig, Config]
}
type Repo interface {
	GetName() string
	IsWeb() bool
	GetConfig() Config
}
type RepoSql struct {
	wx.Provider[RepoSql, Repo]
	config Config
}

func (r *RepoSql) GetName() string {
	return "RepoSql"
}
func (r *RepoSql) IsWeb() bool {
	return r.Context != nil
}
func (r *RepoSql) GetConfig() Config {
	return r.config
}

type MyApp struct {
	wx.Inject[MyApp]
	config Config
	repo   Repo
}

func TestContainer(t *testing.T) {
	(&MyApp{}).Register(func(svc *MyApp) (*MyApp, error) {

		svc.config = &DbConfig{}
		svc.repo = &RepoSql{}
		return svc, nil
	})
	app, err := (&MyApp{}).New()
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, app.repo.GetName(), "RepoSql")

}
func BenchmarkTestContainer(b *testing.B) {
	(&MyApp{}).Register(func(svc *MyApp) (*MyApp, error) {

		svc.config = &DbConfig{}
		svc.repo = &RepoSql{}
		return svc, nil
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		app, err := (&MyApp{}).New()
		assert.NoError(b, err)
		assert.NotNil(b, app)
		assert.Equal(b, app.repo.GetName(), "RepoSql")
	}
}

func init() {
	(&DbConfig{}).Register(func(config *DbConfig) (Config, error) {
		return config, nil
	})
	(&RepoSql{}).Register(func(repo *RepoSql) (Repo, error) {
		var err error
		repo.config, err = (&DbConfig{}).New()
		if err != nil {
			return nil, err
		}
		return repo, nil
	})
}
func TestNew(t *testing.T) {

	repo, err := (&RepoSql{}).New()
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	repo1, err := (&RepoSql{}).NewOnce()
	assert.NoError(t, err)
	assert.NotNil(t, repo1)
	repo2, err := (&RepoSql{}).NewOnce()
	assert.NoError(t, err)
	assert.NotNil(t, repo2)

	assert.Same(t, repo1, repo2)

}
func TestNewOnce(t *testing.T) {

	repo1, err := (&RepoSql{}).NewOnce()
	assert.NoError(t, err)
	assert.NotNil(t, repo1)
	repo2, err := (&RepoSql{}).NewOnce()
	assert.NoError(t, err)
	assert.NotNil(t, repo2)

	assert.Same(t, repo1, repo2)

}
func BenchmarkXxx(b *testing.B) {
	for i := 0; i < b.N; i++ {

		repo, err := (&RepoSql{}).New()
		assert.NoError(b, err)
		assert.NotNil(b, repo)
		assert.False(b, repo.IsWeb())

	}
}
