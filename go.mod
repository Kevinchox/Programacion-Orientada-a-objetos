// Módulo principal del proyecto Go
module github.com/Kevinchox/Programacion-Orientada-a-objetos

// Versión mínima de Go requerida
// Puedes cambiar esto a tu versión instalada (ej: 1.20, 1.21)
go 1.23.0

toolchain go1.24.3 // Toolchain recomendada

require (
	github.com/google/uuid v1.6.0 // Para generación de UUIDs
	golang.org/x/crypto v0.39.0 // Para bcrypt en users/service.go
)
