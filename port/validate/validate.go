package validate

import "gopkg.in/go-playground/validator.v8"

var Validator *validator.Validate

func init() {
	Validator = validator.New(&validator.Config{TagName: "validate"})
}
