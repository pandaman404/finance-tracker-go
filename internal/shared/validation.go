package shared

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) []string {
	var validationErrors validator.ValidationErrors
	var messages []string

	if errors, ok := err.(validator.ValidationErrors); ok {
		validationErrors = errors
	} else {
		return []string{ErrJSONBinding.Error()}
	}

	for _, e := range validationErrors {
		field := e.Field()
		switch e.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("El campo '%s' es requerido", field))
		case "max":
			messages = append(messages, fmt.Sprintf("El campo '%s' no puede superar %s caracteres", field, e.Param()))
		case "min":
			messages = append(messages, fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres", field, e.Param()))
		case "oneof":
			messages = append(messages, fmt.Sprintf("El campo '%s' debe ser uno de: %s", field, e.Param()))
		case "email":
			messages = append(messages, fmt.Sprintf("El campo '%s' debe ser un email válido", field))
		default:
			messages = append(messages, fmt.Sprintf("El campo '%s' es inválido", field))
		}
	}

	return messages
}
