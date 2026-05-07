package account

import "errors"

var (
	ErrInvalidName      = errors.New("Nombre de cuenta inválido")
	ErrInvalidType      = errors.New("Tipo de cuenta inválido")
	ErrUserNotFound     = errors.New("Usuario no encontrado")
	ErrAccountExists    = errors.New("Cuenta ya existe")
	ErrNotFound         = errors.New("Cuenta no encontrada")
	ErrAccountsNotFound = errors.New("No se encontraron cuentas para este usuario")
)
