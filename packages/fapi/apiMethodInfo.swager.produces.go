package fapi

import (
	"mime"
	"strings"
)

func (m *apiMethodInfo) GetResponseContentType() string {
	if m.ResponseContentType != "" {
		return m.ResponseContentType

	}
	switch m.httpMethod {
	case "GET":
		if strings.Contains(m.tags.Url, ".") {
			items := strings.Split(m.tags.Url, ".")
			extFile := items[len(items)-1]
			m.ResponseContentType = mime.TypeByExtension("." + extFile)
			return m.ResponseContentType

		} else {
			m.ResponseContentType = "application/octet-stream"
			return m.ResponseContentType

		}
	case "POST":
		m.ResponseContentType = "application/json"
		return m.ResponseContentType
	default:
		{
			m.ResponseContentType = "application/json"
			return m.ResponseContentType
		}
	}
}
func (m *apiMethodInfo) GetRequestContentType() string {
	if m.requestContentType != "" {
		return m.requestContentType
	} else if m.hasUploadFile {
		m.requestContentType = "application/x-www-form-urlencoded"
		return m.requestContentType
	} else if m.httpMethod == "post" {
		m.requestContentType = "application/json"
		return m.requestContentType
	} else if m.httpMethod == "get" {
		m.requestContentType = "*"
		return m.requestContentType
	} else {
		m.requestContentType = "application/json"
		return m.requestContentType
	}
}
