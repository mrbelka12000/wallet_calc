package validator

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateStruct(obj interface{}) error {
	log.Println(fmt.Sprintf("%+v", obj))
	return validate.Struct(obj)
}
