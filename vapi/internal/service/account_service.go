package service

import (
	"vcache"
	"vdb"
)

type AccountService struct {
	db    *vdb.TenantDB
	cache *vcache.Cache
}
