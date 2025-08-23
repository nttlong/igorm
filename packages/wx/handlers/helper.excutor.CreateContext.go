package handlers

import (
	"net/http"
	"reflect"
	"strings"
	wxErr "wx/errors"
	"wx/internal"
)

func (reqExec *RequestExecutor) GetParamFieldOfHandlerContext(typ reflect.Type, fieldName string) (reflect.StructField, bool) {
	key := typ.String() + "/" + fieldName
	ret, err := internal.OnceCall(key, func() (*reflect.StructField, error) {
		field, ok := typ.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, fieldName)
		})
		if !ok {
			return nil, nil
		}
		return &field, nil
	})
	if err != nil {
		return reflect.StructField{}, false
	}
	return *ret, true

}
func (reqExec *RequestExecutor) CreateHandlerContext(info HandlerInfo, r *http.Request, w http.ResponseWriter) (*reflect.Value, error) {
	var placeHolders [][]string
	if info.IsRegexHandler {
		placeHolders = info.RegexUriFind.FindAllStringSubmatch(r.URL.Path, -1)
		if len(placeHolders) == 0 {
			return nil, wxErr.NewRegexUriNotMatchError("regex uri not match")
		}
	}
	ret := reflect.New(info.TypeOfArgsElem)
	ctx := Handler{
		Req: r,
		Res: w,
	}
	ctxField := ret.Elem().FieldByIndex(info.FieldIndex)
	if ctxField.IsValid() {
		if ctxField.Kind() == reflect.Ptr {
			ctxField.Set(reflect.ValueOf(&ctx))
		} else {
			ctxField.Set(reflect.ValueOf(ctx))
		}
	}
	if len(placeHolders) == 1 {

		for i, x := range info.UriParams {
			field := ret.Elem().FieldByIndex(x.FieldIndex)
			if field.IsValid() {
				if field.Kind() == reflect.String {
					field.SetString(placeHolders[0][i+1])
				} else if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.String {
					field.Set(reflect.ValueOf(&placeHolders[0][i+1]))
				}
			}
		}
	}
	if info.IsQueryUri {
		url, err := r.URL.Parse(r.URL.Path)
		if err != nil {
			return nil, err
		}

		query := url.Query()
		for k, x := range query {
			field, ok := reqExec.GetParamFieldOfHandlerContext(info.TypeOfArgsElem, k)
			if ok {
				fieldSet := ret.Elem().FieldByIndex(field.Index)
				if fieldSet.IsValid() {
					if fieldSet.Kind() == reflect.String {
						fieldSet.SetString(x[0])
					} else if fieldSet.Kind() == reflect.Ptr {
						if fieldSet.Type().Elem().Kind() == reflect.String {
							fieldSet.Set(reflect.ValueOf(&x[0]))
						}
					} else if fieldSet.Kind() == reflect.Slice {
						if fieldSet.Type().Elem().Kind() == reflect.String {
							fieldSet.Set(reflect.ValueOf(x))
						} else if fieldSet.Type().Elem().Kind() == reflect.Ptr {
							if fieldSet.Type().Elem().Elem().Kind() == reflect.String {
								vals := make([]*string, len(x))
								for i, v := range x {
									vals[i] = &v
								}

								fieldSet.Set(reflect.ValueOf(vals))
							}
						}
					}

				}
			}

		}

	}

	return &ret, nil
}
