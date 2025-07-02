package v1

import (
	"encoding/json"
	"net/http"

	"github.com/mrbelka12000/wallet_calc/internal/converter"
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/dto"
)

func (c *Controller) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.UserCU
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.errorResponse(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	dmUser := converter.ConvertFromDTOUserCUToDomainModel(req)
	if err := c.userUsecase.Register(r.Context(), dmUser); err != nil {
		c.errorResponse(w, err, 0)
		return
	}
}

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.errorResponse(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := c.userUsecase.Login(r.Context(), domainmodel.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.errorResponse(w, err, 0)
		return
	}

	_ = user
}

func (c *Controller) Profile(w http.ResponseWriter, r *http.Request) {
	// TODO get user id from context
}
