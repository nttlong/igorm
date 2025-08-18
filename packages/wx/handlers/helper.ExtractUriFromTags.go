package handlers

import (
	"strings"
	"wx/internal"
)

func (h *helperType) ExtractUriFromTags(tags []string) string {
	key := strings.Join(tags, "**") + "/helperType/ExtractUriFromTags"
	ret, _ := internal.OnceCall(key, func() (*string, error) {
		ret := ""
		for i := len(tags) - 1; i >= 0; i-- {
			tag := tags[i]
			if tag == "" {
				continue
			}
			items := strings.Split(tag, ";")

			for _, item := range items {
				uriVal := ""
				if strings.HasPrefix(item, "uri:") {
					uriVal = item[4:]

				} else if item != "" && !strings.Contains(item, ":") {
					uriVal = item
				}
				if uriVal != "" {
					if strings.Contains(ret, "@") {
						ret = strings.Replace(ret, "@", uriVal, 1)
					} else {
						ret += "/" + uriVal
					}
				}
			}

		}
		ret = strings.TrimPrefix(strings.TrimSuffix(ret, "/"), "/")
		return &ret, nil
	})
	return *ret

}
