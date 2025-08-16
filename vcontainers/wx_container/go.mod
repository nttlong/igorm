module wx_container

go 1.24.5

replace vdi => ./../../packages/vdi

replace wx => ./../../packages/wx

require (
	vdi v0.0.0-00010101000000-000000000000
	wx v0.0.0-00010101000000-000000000000
)

require github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
