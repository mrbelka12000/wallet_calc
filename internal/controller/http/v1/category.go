package v1

import (
	"encoding/json"
	"net/http"

	"github.com/mrbelka12000/wallet_calc/internal/converter"
)

func (c *Controller) CategoryList(w http.ResponseWriter, r *http.Request) {

	categories, err := c.categoryUsecase.List(r.Context())
	if err != nil {
		c.errorResponse(w, err, 0)
		return
	}

	resp := converter.ConvertFromDMCategoriesToDTO(categories)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		c.errorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
