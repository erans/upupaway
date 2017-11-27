package utils

import uuid "github.com/satori/go.uuid"

// GenerateUploadID generates a unique upload ID
func GenerateUploadID() string {
	return uuid.NewV4().String()
}

// SafeCastToString safely casts an interface to a string, otherwise it will return an empty string
func SafeCastToString(v interface{}) string {
	if v != nil {
		return v.(string)
	}

	return ""
}
