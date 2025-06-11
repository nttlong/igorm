module dynacallexample

go 1.24.3

require (
	dbx v0.0.0-20220315162249-d5a7a7a3d57d
	dynacall v0.0.0-20220315162249-d5a7a7a3d57d
	github.com/stretchr/testify v1.10.0
	unvs.br.auth v0.0.0-20220315162249-d5a7a7a3d57d
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/microsoft/go-mssqldb v1.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto v0.0.0-20250512202823-5a2f75b736a9 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace dynacall => ../dynacall

replace dbx => ../dbx

replace unvs.br.auth => ../businessRules/auth
