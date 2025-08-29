module hello

go 1.25.0

replace wx => ./../../packages/wx

require (
	github.com/go-chi/chi v1.5.5
	github.com/stretchr/testify v1.10.0
	wx v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
