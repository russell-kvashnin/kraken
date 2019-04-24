package error

import "fmt"

const (
	defaultErrorDomain = "DEFAULT_DOMAIN"
	defaultErrCode     = "ERROR"

	ErrLvlFatal   = "fatal"
	ErrLvlError   = "error"
	ErrLvlWarning = "warning"

	errorMessageTemplate = "[Level]: '%s' [Domain]: '%s' [ErrorCode]: '%s' [Message]: '%s'"
)

// Advanced error structure
type err struct {
	level    string
	domain   string
	code     string
	original error
	details  map[string]string
}

// Advanced error interface
type Error interface {
	Level() string
	Domain() string
	Code() string
	Original() error
	Details() map[string]string
	Error() string
}

// Constructor
func NewErr(lvl string, domain string, code string, original error, details map[string]string) Error {
	if lvl == "" {
		lvl = ErrLvlError
	}

	if domain == "" {
		domain = defaultErrorDomain
	}

	if code == "" {
		code = defaultErrCode
	}

	err := new(err)
	err.level = lvl
	err.domain = domain
	err.code = code
	err.original = original
	err.details = details

	return err
}

// Std error interface impl
func (err err) Error() string {
	var (
		errString string
	)

	errString = fmt.Sprintf(
		errorMessageTemplate,
		err.level, err.domain,
		err.code,
		err.original.Error())

	return errString
}

func (err err) Level() string {
	return err.level
}

func (err err) Domain() string {
	return err.domain
}

func (err err) Code() string {
	return err.code
}

func (err err) Original() error {
	return err.original
}

func (err err) Details() map[string]string {
	return err.details
}

func (err err) DetailsString() string {
	var details string

	for key, value := range err.details {
		details += fmt.Sprintf("[%s]: '%s' ", key, value)
	}

	return details
}
