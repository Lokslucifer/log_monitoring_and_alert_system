package Customerrors

import (
	"errors"
)

// Server Error Messages
var (
	FileRotatedError = errors.New("log file rotated")

)
