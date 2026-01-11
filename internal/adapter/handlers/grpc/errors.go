package grpc

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidationErrorItem представляет отдельную ошибку валидации с полями path и reason.
type ValidationErrorItem struct {
	Path   string
	Reason string
}

// NewValidationErrorItem создает новый ValidationErrorItem с указанными path и reason.
func NewValidationErrorItem(path, reason string) ValidationErrorItem {
	return ValidationErrorItem{
		Path:   path,
		Reason: reason,
	}
}

// ValidationError возвращает gRPC ошибку с кодом InvalidArgument для ошибок валидации
// Принимает слайс ValidationErrorItem, где каждый элемент содержит path (путь к полю) и reason (причину ошибки).
func ValidationError(errors []ValidationErrorItem) error {
	if len(errors) == 0 {
		return status.Error(codes.InvalidArgument, "validation failed")
	}

	// Создаем статус с кодом InvalidArgument
	st := status.New(codes.InvalidArgument, "validation failed")

	// Преобразуем ValidationErrorItem в FieldViolation для BadRequest
	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, 0, len(errors))
	for _, err := range errors {
		fieldViolations = append(fieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       err.Path,
			Description: err.Reason,
		})
	}

	// Создаем BadRequest с FieldViolations
	badReq := &errdetails.BadRequest{
		FieldViolations: fieldViolations,
	}

	// Добавляем детали к статусу
	st, err := st.WithDetails(badReq)
	if err != nil {
		// Если не удалось добавить детали, возвращаем простую ошибку
		return status.Error(codes.InvalidArgument, "validation failed")
	}

	return st.Err()
}

// InternalError возвращает gRPC ошибку с кодом Internal для внутренних ошибок.
func InternalError(message string) error {
	return status.Error(codes.Internal, message)
}

// NotFoundError возвращает gRPC ошибку с кодом NotFound.
func NotFoundError(message string) error {
	return status.Error(codes.NotFound, message)
}
