package responseutil

import (
	"github.com/go-sql-driver/mysql"
	"strings"
)

func ExtractFieldFromError(err error, entityName string) string {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok || mysqlErr.Number != 1062 {
		return ""
	}

	// Extract the field name from the error message
	// Example error message: "Error 1062: Duplicate entry '0362356190' for key 'Account.phone'"
	msg := mysqlErr.Message
	prefix := "for key '"
	start := strings.Index(msg, prefix)

	if start == -1 {
		return "unknown_field"
	}

	start += len(prefix)
	end := strings.Index(msg[start:], "'")

	if end == -1 {
		return "unknown_field"
	}

	fullKey := msg[start : start+end] // e.g., "Account.phone"
	parts := strings.Split(fullKey, ".")

	if len(parts) != 2 {
		return "unknown_field"
	}

	// Ensure the entity name matches to extract the correct field
	if strings.ToLower(parts[0]) != strings.ToLower(entityName) {
		return "unknown_field"
	}

	return parts[1] // e.g., "phone"
}
