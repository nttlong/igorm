module unvs.br.auth

go 1.24.3

require (
	dbmodels v0.0.0-00010101000000-000000000000
	dbx v0.0.0-20210811154550-d7f2e956a1a9
	dynacall v0.0.0-20210615155502-d5a7c7a7d97a
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	golang.org/x/crypto v0.38.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/labstack/echo/v4 v4.13.4 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/microsoft/go-mssqldb v1.8.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xwb1989/sqlparser v0.0.0-20180606152119-120387863bf2 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto v0.0.0-20250512202823-5a2f75b736a9 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace dynacall => ../../dynacall

replace dbmodels => ../../dbmodels

replace dbx => ../../dbx
