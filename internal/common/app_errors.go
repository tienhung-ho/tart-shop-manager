package common

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"time"
)

type AppError struct {
	StatusCode int       `json:"status_code"`
	RootErr    error     `json:"-"`
	Message    string    `json:"message"`
	Log        string    `json:"log"`
	Key        string    `json:"error_key"`
	Timestamp  time.Time `json:"timestamp"`
}

func NewFullErrorResponse(statusCode int, root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
		Timestamp:  time.Now(), // Add the current time
	}
}

func NewErrorResponse(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
		Timestamp:  time.Now(), // Add the current time
	}
}

func NewUnauthorized(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
		Timestamp:  time.Now(), // Add the current time
	}
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "Invalid request", err.Error(), "ErrInvalidRequest")
}

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong with the server", err.Error(), "ErrInternal")
}

func TokenExpired(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("%s token expired", strings.ToLower(entity)),
		fmt.Sprintf("ErrTokenExpired%s", entity), entity)
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity), entity)
}

func ErrCannotGetReport(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot get report %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGetReport%s", entity), entity)
}

func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity), entity)
}

func ErrRecordExist(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot handle %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannot handle, record exist %s", entity), entity)
}

func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity), entity)
}

func ErrCannotSort(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Idvalid field %s on params", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotSort&ListItem%s", entity), entity)
}

func ErrCannotGenerateKey(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot Create key of %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreateKey%s", entity), entity)
}

func ErrInvalidStatus(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Invalid status for %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrInvalidStatus%s", entity),
		entity)
}

func ErrInvalidGender(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Invalid gender for %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrInvalidGender%s", entity),
		entity)
}

func ErrInvalidAccount(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Your account status has been disabled: %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrInvalid%s", entity),
		entity)
}

func ErrInvalidEnum(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Invalid enum for %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrInvalidEnum%s", entity),
		entity)
}

func ErrEmailInvalid(entity string, err error) *AppError {
	return NewFullErrorResponse(http.StatusForbidden, err,
		fmt.Sprintf("Cannot login, wrong email %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotLogin%s", entity), entity)
}

func ErrPasswordInvalid(entity string, err error) *AppError {
	return NewFullErrorResponse(http.StatusForbidden, err,
		fmt.Sprintf("Cannot login, wrong password %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotLogin%s", entity), entity)
}

func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity), entity)
}

func ErrEntityDeleted(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Record has been deleted %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrRecordHasBeenDeleted%s", entity), entity)
}

func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreate%s", entity), entity)
}

func ErrNoPermission(entity string, err error) *AppError {
	return NewErrorResponse(err,
		"You don't have permission to access this resource",
		"ErrNoPermission", entity)
}

var RecordNotFound = errors.New("record not found")

func ErrNotFoundEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot not found %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotNotFound%s", entity), entity)
}

func ErrOutOffQuantity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Out off quantity %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrOutOffQuantity%s", entity), entity)
}

func ErrNotFoundToken(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot not found %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotNotFound%s", entity), entity)
}
func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}

// Error implements the error interface for AppError
func (e *AppError) Error() string {
	return e.RootError().Error()
}

func ErrDuplicateEntry(entity, field string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Duplicate entry found for %s in %s", field, strings.ToLower(entity)),
		fmt.Sprintf("ErrDuplicateEntry%s%s", entity, strings.Title(field)), entity)
}

func ErrValidation(validationErrors validator.ValidationErrors) *AppError {
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

func ErrCanNotBindEntity(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Cannot not bind %s or data is empty", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotNotBindData%s", entity), entity)
}

var (
	ErrCloudConnectionFailedSentinel = fmt.Errorf("failed to connect to cloud service")
	// ... other sentinel errors
)

// Cloud-related error responses

func ErrCannotUploadFile(entity string, err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err,
		fmt.Sprintf("Cannot upload %s", strings.ToLower(entity)),
		err.Error(), fmt.Sprintf("ERR_CANNOT_UPLOAD_%s", strings.ToUpper(entity)))
}

func ErrCloudConnectionFailed(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err,
		"Failed to connect to cloud service",
		err.Error(), "ERR_CLOUD_CONNECTION_FAILED")
}

func ErrInvalidCloudCredentials(err error) *AppError {
	return NewFullErrorResponse(http.StatusUnauthorized, err,
		"Invalid cloud credentials",
		err.Error(), "ERR_INVALID_CLOUD_CREDENTIALS")
}

func ErrFileTooLarge(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("%s file size exceeds the limit", strings.ToLower(entity)),
		err.Error(), fmt.Sprintf("ERR_%s_FILE_TOO_LARGE", strings.ToUpper(entity)))
}

func ErrUnsupportedFileType(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("%s file type is not supported", strings.ToLower(entity)),
		err.Error(), fmt.Sprintf("ERR_UNSUPPORTED_%s_FILE_TYPE", strings.ToUpper(entity)))
}

func ErrCannotDeleteFile(entity string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("%s can not delete file", strings.ToLower(entity)),
		err.Error(), fmt.Sprintf("ERR_DELETE_%s_FILE", strings.ToUpper(entity)))
}

func ErrCloudTimeout(err error) *AppError {
	return NewFullErrorResponse(http.StatusGatewayTimeout, err,
		"Cloud service timeout",
		err.Error(), "ERR_CLOUD_TIMEOUT")
}

func ErrCloudServiceUnavailable(err error) *AppError {
	return NewFullErrorResponse(http.StatusServiceUnavailable, err,
		"Cloud service is unavailable",
		err.Error(), "ERR_CLOUD_SERVICE_UNAVAILABLE")
}

// common/errors.go

func ErrForeignKeyConstraint(entity, field string, err error) *AppError {
	return NewErrorResponse(err,
		fmt.Sprintf("Invalid %s: %s does not exist.", field, entity),
		fmt.Sprintf("ErrForeignKeyConstraint%s%s", entity, strings.Title(field)),
		entity)
}
