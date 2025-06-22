package dbx

import "fmt"

type Updater struct {
}

func (t *DBXTenant) Update(entity interface{}) *Updater {

	panic(fmt.Sprintf("Update method is not supported for tenant db %s", `igorm\dbx\tenantdb.updator.go`))

}
