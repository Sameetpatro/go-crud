package response

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/go-playground/validator/v10"
)

type Response struct{
	Status string `json:"status"`
	Error string `json:"error"`
}

const (
	StatusSuccess = "SUCCESS"
	StatusError = "ERROR"
)


func WriteJson(res http.ResponseWriter, statusCode int, data any) error{
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)

	return json.NewEncoder(res).Encode(data)
}

func GeneralError(err error) Response{
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}

func ValidationError(err validator.ValidationErrors) Response{
	var errMsgs []string
	for _, err := range err{
		switch err.ActualTag(){
		case "required":
			errMsgs = append(errMsgs, err.Field() + " is required")
		default:
			errMsgs = append(errMsgs, err.Field() + " is not valid")
		}
	}
	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs, ","),
	}
}