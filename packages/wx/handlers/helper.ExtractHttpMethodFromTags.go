package handlers

import (
	"strings"
	"wx/internal"
)

func (h *helperType) ExtractHttpMethodFromTags(tags []string) string {
	key := strings.Join(tags, "**") + "/helperType/ExtractHttpMethodFromTags"
	ret, _ := internal.OnceCall(key, func() (*string, error) {
		ret := ""
		for i := len(tags) - 1; i >= 0; i-- {
			tag := tags[i]
			if tag == "" {
				continue
			}
			items := strings.Split(tag, ";")
			for _, item := range items {
				if strings.HasPrefix(item, "method:") {
					ret = strings.ToUpper(item[7:])

				}
			}
		}

		return &ret, nil
	})
	return *ret
}
