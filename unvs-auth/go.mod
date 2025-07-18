module unvs-auth

go 1.24.3

require (
	dbv v0.0.0-20211015132956-d9d5a7a2a1a9
	github.com/golang-jwt/jwt/v5 v5.2.3
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.38.0
	vdi v0.0.0-20220223132316-b832511892a9
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/microsoft/go-mssqldb v1.9.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace dbv => ../dbv

replace vdi => ../vdi
