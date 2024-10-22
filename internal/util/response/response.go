package responseutil

import (
	"github.com/go-sql-driver/mysql"
	"reflect"
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

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	// Nếu là con trỏ, ta phải lấy giá trị thực sự của nó
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Bỏ qua trường hợp không thể xuất ra được (unexported field)
		if !field.CanInterface() {
			continue
		}

		// Kiểm tra giá trị mặc định và bỏ qua
		isDefault := false
		switch field.Kind() {
		case reflect.String:
			isDefault = field.String() == ""
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			isDefault = field.Int() == 0
		case reflect.Float32, reflect.Float64:
			isDefault = field.Float() == 0
		case reflect.Slice, reflect.Map:
			isDefault = field.Len() == 0
		case reflect.Ptr:
			isDefault = field.IsNil()
		}

		if !isDefault {
			result[fieldType.Tag.Get("json")] = field.Interface()
		}
	}

	return result
}
