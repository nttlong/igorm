module media

go 1.25.0

replace wx => ./../../packages/wx

require (
	golang.org/x/crypto v0.41.0
	wx v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
