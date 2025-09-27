package validators

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var validatorObj *validator.Validate

func Init() {
	validatorObj = validator.New()
}

/*
Validate struct with tag, of request body, param, query, etc.
*/
func validateStructWithTags(data interface{}) error {
	errs := validatorObj.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			if err != nil {
				log.Println(err.Kind())
				// errMsg := fmt.Sprintf("Field '%s' : '%v' | Needs to pass '%s' validation", err.Field(), err.Value(), err.Tag())
				errMsg := fmt.Sprintf("Field '%s' | Needs to pass '%s' validation", err.Field(), err.Tag())
				return errors.New(errMsg)
			}
		}
	}

	return nil
}

/*
Validate uuid string input
*/
func ValidateUuid(uuid string) error {
	err := validatorObj.Var(uuid, "required,uuid")
	if err != nil {
		return errors.New("invalid uuid")
	}

	return nil
}

/*
Parse and validate body
*/
func ParseAndValidateBody(c echo.Context, out interface{}) error {
	if err := c.Bind(out); err != nil {
		return errors.New(err.Error())
	}
	return validateStructWithTags(out)
}

/*
Parse and validate params
*/
func ParseAndValidateQueryParam(c echo.Context, out interface{}) error {
	if err := c.Bind(out); err != nil {
		return errors.New(err.Error())
	}
	return validateStructWithTags(out)
}
