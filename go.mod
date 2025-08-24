module github.com/GoogleCloudPlatform/golang-samples/run/helloworld

go 1.25.0

replace wx => ./packages/wx

require wx v0.0.0-00010101000000-000000000000

require github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
