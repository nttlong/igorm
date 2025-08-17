package handlers

import (
	"wx/internal"
)

/*
handlerInfoExtractUriParams extracts all substrings enclosed in curly braces '{}'
from the given URI string, along with their positions (based on segments split by '/').

For example:

Given URI: "abc/{word 1}/abc/dbc/{u2}/{u3}"

The function will return a slice of uriParaim structs:
[

	{Position: 1, Name: "word 1"},
	{Position: 4, Name: "u2"},
	{Position: 5, Name: "u3"},

]

Where:
- Position is the zero-based index of the segment in the URI path split by '/'.
- Name is the trimmed string inside the braces '{}'.

@return []uriParaim - a slice containing extracted parameters with their position and name.
*/
func (h *helperType) ExtractUriParams(uri string) []uriParam {
	key := uri + "/helperType/ExtractUriParams"
	ret, _ := internal.OnceCall(key, func() (*[]uriParam, error) {

		params := []uriParam{}
		segments := h.SplitUriSegments(uri)

		for i, segment := range segments {
			// Check if segment contains a URI parameter enclosed in {}
			name := h.ExtractNameInBraces(segment)
			if name != "" {
				isSlug := false

				if name[0] == '*' {
					name = name[1:]
					isSlug = true
				}
				params = append(params, uriParam{
					Position: i,
					Name:     name,
					IsSlug:   isSlug,
				})
			}
		}

		return &params, nil
	})
	return *ret
}
