package formaterror

import (
	"strings"
)

// errorMessages is a map of strings to strings used to store error messages
var errorMessages = make(map[string]string)

// err is a variable used to store errors
var err error

// FormatError takes a string as an argument and returns a map of strings to strings
// containing error messages.
//
// It checks if the string contains certain keywords and adds the corresponding
// error messages to the map. If the map is empty, it adds an "Incorrect Details"
// error message.
func FormatError(errString string) map[string]string {
	if strings.Contains(errString, "username") {
		errorMessages["Taken_username"] = "Username Already Taken"
	}
	if strings.Contains(errString, "email") {
		errorMessages["Taken_email"] = "Email Already Taken"
	}
	if strings.Contains(errString, "hashedPassword") {
		errorMessages["Incorrect_password"] = "Incorrect Password"
	}
	if strings.Contains(errString, "record not found") {
		errorMessages["No_record"] = "No Record Found"
	}
	if len(errorMessages) > 0 {
		return errorMessages
	}
	if len(errorMessages) == 0 {
		errorMessages["Incorrect_details"] = "Incorrect Details"
		return errorMessages
	}
	return nil
}
