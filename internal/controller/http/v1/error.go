package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mrbelka12000/wallet_calc/internal/usecase"
)

type (
	ErrorResponse struct {
		Message string `json:"message"`
	}
)

func (c *Controller) errorResponse(w http.ResponseWriter, err error, code int) {
	if code == 0 {
		code = http.StatusInternalServerError
		var clientErr usecase.ClientError
		if errors.As(err, &clientErr) {
			code = http.StatusBadRequest
		}
	}

	c.log.With("error", err).Error("error response")
	body, err := json.Marshal(ErrorResponse{Message: err.Error()})
	if err != nil {
		c.log.With("error", err).Error("marshalling error response")
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
