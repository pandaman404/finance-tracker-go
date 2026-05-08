package category

import "errors"

var (
	ErrInvalidName        = errors.New("El nombre de la categoría no puede estar vacío")
	ErrInvalidType        = errors.New("Tipo de categoría inválida")
	ErrCategoryExists     = errors.New("Ya existe una categoría con ese nombre para este usuario")
	ErrCategoryNotFound   = errors.New("Categoría no encontrada")
	ErrCategoriesNotFound = errors.New("No se encontraron categorías para este usuario")
	ErrUserNotFound       = errors.New("Usuario no encontrado")
)
