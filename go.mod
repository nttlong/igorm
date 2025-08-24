module github.com/GoogleCloudPlatform/golang-samples/run/helloworld

go 1.25.0

replace wx => ./packages/wx

replace xauth => ./packages/xauth

replace xconfig => ./packages/xconfig

replace vdb => ./packages/vdb

require (
	wx v0.0.0-00010101000000-000000000000
	xauth v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/microsoft/go-mssqldb v1.9.2 // indirect
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	vdb v0.0.0-00010101000000-000000000000 // indirect
	xconfig v0.0.0-00010101000000-000000000000 // indirect
)
