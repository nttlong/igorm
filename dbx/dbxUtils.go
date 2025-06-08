package dbx

import "sync"

type dbxEntityCacheType struct {
	uk sync.Map // cache for unique key fields map[table_name+"_"+uk_name] = []string{uk_fields}
}

func (d *dbxEntityCacheType) set_uk(uk_name string, uk_fields []string) {
	d.uk.Store(uk_name, uk_fields)
}
func (d *dbxEntityCacheType) get_uk(uk_name string) []string {
	if uk, ok := d.uk.Load(uk_name); ok {
		return uk.([]string)
	}
	return nil
}

var dbxEntityCache dbxEntityCacheType

func init() {
	dbxEntityCache = dbxEntityCacheType{
		uk: sync.Map{},
	}
}
