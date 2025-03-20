package validator

import (
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate     *validator.Validate
	validateOnce sync.Once
)

func getValidator() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New() // Se ejecuta solo una vez
	})
	return validate
}

func ValidateStruct(s interface{}) error {
	validate := getValidator()
	err := validate.Struct(s)
	if err != nil {
		var msg string
		for _, err := range err.(validator.ValidationErrors) {
			if err.ActualTag() == "required" {
				msg += fmt.Sprintf("Field '%s' is required\n", err.Field())
			} else {
				msg += fmt.Sprintf("Field '%s' should be %s %s\n", err.Field(), err.Tag(), err.Param())
			}
			// log.Printf("Error en el campo '%s'", err.Field())
			// //fmt.Println(err.Namespace())
			// //fmt.Println(err.Field())
			// //fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// //fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// //fmt.Println(err.Kind())
			// //fmt.Println(err.Type())
			// //fmt.Println(err.Value())
			// fmt.Println(err.Param())

		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}
