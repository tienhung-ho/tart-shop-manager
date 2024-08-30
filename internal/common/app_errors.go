package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type appError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func NewFullErrorResponse(statusCode int, root error, msg, log, key string) *appError {
	return &appError{
		StatusCode: statusCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewErrorResponse(root error, msg, log, key string) *appError {
	return &appError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewUnauthorized(root error, msg, log, key string) *appError {
	return &appError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func ErrDB(err error) *appError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrInvalidRequest(err error) *appError {
	return NewErrorResponse(err, "Invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInternal(err error) *appError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong with the server", err.Error(), "ErrInternal")
}

func TokenExpired(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("%s token expired", strings.ToLower(entity)),
		fmt.Sprintf("ErrTokenExpired%s", entity), entity)
}

func ErrCannotListEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity), entity)
}

func ErrCannotGetEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity), entity)
}

func ErrRecordExist(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannot create, record exist %s", entity), entity)
}

func ErrCannotUpdateEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity), entity)
}

func ErrEmailOrPasswordInvalid(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot login, wrong password or email %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotLogin%s", entity), entity)
}

func ErrCannotDeleteEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity), entity)
}

func ErrEntityDeleted(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Record has been deleted %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrRecordHasBeenDeleted%s", entity), entity)
}

func ErrCannotCreateEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreate%s", entity), entity)
}

func ErrNoPermission(entity string, err error) *appError {
	return NewErrorResponse(err,
		"You have no permission",
		"ErrNoPermission", entity)
}

var RecordNotFound = errors.New("record not found!")

func ErrNotFoundEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot not found %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotNotFound%s", entity), entity)
}

func (e *appError) RootError() error {
	if err, ok := e.RootErr.(*appError); ok {
		return err.RootError()
	}
	return e.RootErr
}

// Error implements the error interface for appError
func (e *appError) Error() string {
	return e.RootError().Error()
}

func ErrDuplicateEntry(entity, field string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Duplicate entry found for %s in %s", field, strings.ToLower(entity)),
		fmt.Sprintf("ErrDuplicateEntry%s%s", entity, strings.Title(field)), entity)
}

func ErrValidation(validationErrors validator.ValidationErrors) *appError {
	var errMsgs []string
	for _, err := range validationErrors {
		field := err.Field()
		tag := err.Tag()

		var errMsg string
		switch tag {
		case "required":
			errMsg = fmt.Sprintf("The field '%s' is required.", field)
		case "email":
			errMsg = fmt.Sprintf("The field '%s' must be a valid email address.", field)
		case "min":
			errMsg = fmt.Sprintf("The field '%s' must be at least %s characters long.", field, err.Param())
		case "eqfield":
			errMsg = fmt.Sprintf("The field '%s' must match the field '%s'.", field, err.Param())
		case "vietnamese_phone":
			errMsg = fmt.Sprintf("The field '%s' must be a valid Vietnamese phone number.", field)
		default:
			errMsg = fmt.Sprintf("The field '%s' failed validation with rule '%s'.", field, tag)
		}
		errMsgs = append(errMsgs, errMsg)
	}

	return NewErrorResponse(nil, "Validation failed", strings.Join(errMsgs, "; "), "VALIDATION_ERROR")
}

func ErrCanNotBindEntity(entity string, err error) *appError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot not bind %s or data is empty", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotNotBindData%s", entity), entity)
}
