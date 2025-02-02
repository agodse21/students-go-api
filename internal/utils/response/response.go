package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "ok"
	StatusError = "error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	res := Response{
		Status: StatusError,
		Error:  err.Error(),
	}
	return res
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("Fild %s is required", err.Field()))

		default:
			errMsg = append(errMsg, fmt.Sprintf("Fild %s is invalid", err.Field()))

		}
	}

	res := Response{
		Status: StatusError,
		Error:  strings.Join(errMsg, ", "),
	}
	return res
}
