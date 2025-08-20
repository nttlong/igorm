package handlers

import (
	"reflect"
	"sort"
)

type RouteItem struct {
	Info HandlerInfo
}
type RouteTypes struct {
	Data    map[string]RouteItem
	UriList []string
}

var Routes *RouteTypes = &RouteTypes{
	Data:    map[string]RouteItem{},
	UriList: []string{},
}

func (r *RouteTypes) Add(baseUri string, ins ...any) error {
	insTypes := make([]reflect.Type, len(ins))
	for i, x := range ins {
		if t, ok := x.(reflect.Type); ok {
			insTypes[i] = t
		} else {
			insTypes[i] = reflect.TypeOf(x)
		}
	}
	for _, typ := range insTypes {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct {
			continue
		}
		ptrType := reflect.PointerTo(typ)
		for i := 0; i < ptrType.NumMethod(); i++ {
			method := ptrType.Method(i)
			info, err := Helper.GetHandlerInfo(method)
			if err != nil {
				return err
			}
			if info == nil {
				continue
			}
			if info.IsAbsUri {
				r.Data[info.UriHandler] = RouteItem{
					Info: *info,
				}
				r.UriList = append(r.UriList, info.UriHandler)
			} else {
				r.Data[baseUri+"/"+info.UriHandler] = RouteItem{
					Info: *info,
				}
				r.UriList = append(r.UriList, baseUri+"/"+info.UriHandler)
			}

		}

	}
	//sort r.UriList by len form large to small
	sort.Slice(r.UriList, func(i, j int) bool {
		return len(r.UriList[i]) > len(r.UriList[j])
	})
	return nil
}
