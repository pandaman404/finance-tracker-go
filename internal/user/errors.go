package user

import "errors"

var (
	ErrEmailExists        = errors.New("El email ya está registrado")
	ErrNotFound           = errors.New("Usuario no encontrado")
	ErrInvalidCredentials = errors.New("Email o contraseña incorrectos")
)
