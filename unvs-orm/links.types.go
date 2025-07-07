package orm

import "unvs-orm/internal"

// type Base = internal.Base
type Base struct {
	internal.Base
}

func init() {
	internal.OnRequestBaseFn = func(base *internal.Base) interface{} {
		return &Base{*base}
	}
}
