package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
    Status string
    Error string
}

const (
    StatusOK = "ok"
    StatusError = "Error"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{})  error {

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status) 

    return json.NewEncoder(w).Encode(data)

}

func ValidatonError(errs validator.ValidationErrors) Response {
    var errMsg []string

    for _, err := range errs {
        switch err.ActualTag() {
            case "required":
                errMsg = append(errMsg, fmt.Sprintf("field %s is required field", err.Field()))
            default:
                errMsg = append(errMsg, fmt.Sprintf("field %s is invalid ", err.Field()))
        }
    }

    return Response{
        Status: StatusError,
        Error: strings.Join(errMsg, ","),
    }
}