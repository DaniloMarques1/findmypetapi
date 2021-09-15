package util

import (
	"encoding/json"
	"net/http"

	"github.com/danilomarques1/findmypetapi/dto"
)

func RespondJson(w http.ResponseWriter, status int, body interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func HandleError(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case *ApiError:
		RespondJson(w, v.Code, dto.ErrorDto{Message: v.Message})
	default:
		RespondJson(w, http.StatusInternalServerError,
			dto.ErrorDto{Message: "Unnexpected error"})
	}
}
