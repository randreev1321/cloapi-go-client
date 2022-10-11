package request_tools

import (
	"net/http"
)

func IsError(statusCode int) bool {
	switch statusCode {
	case http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusNotFound, http.StatusNoContent:
		return false
	default:
		return true
	}
}
