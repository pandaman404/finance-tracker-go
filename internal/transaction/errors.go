package transaction

import "errors"

var (
	ErrInvalidType          = errors.New("Tipo de transacción inválido")
	ErrInvalidAmount        = errors.New("Monto inválido")
	ErrInvalidDescription   = errors.New("Descripción inválida")
	ErrTransactionNotFound  = errors.New("Transacción no encontrada")
	ErrTransactionsNotFound = errors.New("No se encontraron transacciones para este usuario")
	ErrAccountNotFound      = errors.New("Cuenta no encontrada")
	ErrCategoryNotFound     = errors.New("Categoría no encontrada")
	ErrCategoryTypeMismatch = errors.New("La categoría no corresponde al tipo de transacción")
)
