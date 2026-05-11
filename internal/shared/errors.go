package shared

import "errors"

var (
	ErrInternalServer  = errors.New("Error interno del servidor")
	ErrUnauthorized    = errors.New("No autorizado")
	ErrInvalidID       = errors.New("ID inválido")
	ErrJSONBinding     = errors.New("Error al procesar la solicitud")
	ErrTooManyRequests = errors.New("Demasiadas solicitudes, intente más tarde")
)
