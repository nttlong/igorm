package caller

import (
	"github.com/labstack/echo/v4"
)

type ExtractError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e ExtractError) Error() string {
	return e.Message
}

type ExtractInfo struct {
	Feature string `json:"feature"`
	Tenant  string `json:"tenant"`
	Lan     string `json:"lan"`
	Action  string `json:"action"`
	Module  string `json:"module"`
}

func ExtractRequireQueryStrings(c echo.Context) (*ExtractInfo, error) {
	feature := c.Request().URL.Query().Get("feature")
	if feature == "" {
		return nil, ExtractError{
			Code:    "INVALID_REQUEST_BODY",
			Message: "query string 'feature' is required",
		}

	}

	tenantName := c.Request().URL.Query().Get("tenant")
	if tenantName == "" {
		return nil, ExtractError{
			Code:    "INVALID_REQUEST_BODY",
			Message: "query string 'tenant' is required",
		}

	}
	lan := c.Request().URL.Query().Get("lan")
	if lan == "" {
		return nil, ExtractError{
			Code:    "INVALID_REQUEST_BODY",
			Message: "query string 'lan' is required",
		}

	}
	action := c.Request().URL.Query().Get("action")
	if action == "" {
		return nil, ExtractError{
			Code:    "INVALID_REQUEST_BODY",
			Message: "query string 'action' is required",
		}

	}
	module := c.Request().URL.Query().Get("module")
	if module == "" {
		return nil, ExtractError{
			Code:    "INVALID_REQUEST_BODY",
			Message: "query string'module' is required",
		}

	}
	return &ExtractInfo{
		Feature: feature,
		Tenant:  tenantName,
		Lan:     lan,
		Action:  action,
		Module:  module,
	}, nil
}
