package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Status string

const (
	StatusAll      Status = "All"
	StatusActive   Status = "Active"
	StatusInactive Status = "Inactive"
	StatusPending  Status = "Pending"
)

var validStatuses = map[Status]bool{
	StatusActive:   true,
	StatusInactive: true,
	StatusPending:  true,
}

func (s Status) IsValid() bool {
	_, ok := validStatuses[s]
	return ok
}

func (s *Status) UnmarshalJSON(b []byte) error {
	var strValue string
	if err := json.Unmarshal(b, &strValue); err != nil {
		return err
	}

	status := Status(strValue)
	if !status.IsValid() {
		return ErrInvalidStatus("data", fmt.Errorf("invalid status value: %s", strValue))
	}

	*s = status
	return nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	if !s.IsValid() {
		return nil, ErrInvalidStatus("data", fmt.Errorf("invalid status value: %s", s))
	}
	return json.Marshal(string(s))
}

func (s *Status) Scan(value interface{}) error {
	if value == nil {
		*s = StatusPending
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*s = Status(v)
	case string:
		*s = Status(v)
	default:
		return fmt.Errorf("unsupported Scan value for Status: %v", value)
	}
	return nil
}

func (s Status) Value() (driver.Value, error) {
	return string(s), nil
}

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
	GenderOther  Gender = "Other"
)

func (g *Gender) Scan(value interface{}) error {
	if value == nil {
		*g = "" // Hoặc gán giá trị mặc định nào đó
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*g = Gender(v)
	case string:
		*g = Gender(v)
	default:
		return fmt.Errorf("unsupported Scan value for Gender: %v", value)
	}
	return nil
}

func (g Gender) Value() (driver.Value, error) {
	return string(g), nil
}

type JSON json.RawMessage

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// CustomDate bao gồm cả ngày, giờ và múi giờ
type CustomDate struct {
	time.Time
}

// Định dạng thời gian chuẩn RFC3339
const DateTimeFormatRFC3339 = time.RFC3339 // "2006-01-02T15:04:05Z07:00"

// Định dạng ngày duy nhất
const DateFormat = "2006-01-02"

// Định dạng datetime cho MySQL
const DateTimeFormatMySQL = "2006-01-02 15:04:05"

// MarshalJSON để đảm bảo định dạng chính xác khi trả về cho client
func (cd *CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte("null"), nil
	}
	// Sử dụng RFC3339 cho JSON
	return []byte(fmt.Sprintf("\"%s\"", cd.Time.Format(DateTimeFormatRFC3339))), nil
}

// UnmarshalJSON để bind chuỗi ngày-giờ từ client vào time.Time
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	strInput := strings.Trim(string(b), "\"")
	if strInput == "null" || strInput == "" {
		cd.Time = time.Time{}
		return nil
	}

	// Thử parse với RFC3339
	parsedTime, err := time.Parse(DateTimeFormatRFC3339, strInput)
	if err != nil {
		// Nếu không thành công, thử parse với định dạng ngày duy nhất
		parsedTime, err = time.Parse(DateFormat, strInput)
		if err != nil {
			return errors.New("invalid datetime format, use RFC3339 or 'YYYY-MM-DD'")
		}
	}
	cd.Time = parsedTime
	return nil
}

// Implement the driver.Valuer interface for database serialization
func (cd CustomDate) Value() (driver.Value, error) {
	if cd.Time.IsZero() {
		return nil, nil
	}
	// Sử dụng định dạng MySQL khi lưu vào DB
	return cd.Time.Format(DateTimeFormatMySQL), nil
}

// Implement the sql.Scanner interface for database deserialization
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		*cd = CustomDate{Time: time.Time{}}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		cd.Time = v
		return nil
	case string:
		parsedTime, err := time.Parse(DateTimeFormatMySQL, v)
		if err != nil {
			// Nếu không thành công, thử parse với RFC3339
			parsedTime, err = time.Parse(DateTimeFormatRFC3339, v)
			if err != nil {
				// Thử parse với định dạng ngày duy nhất
				parsedTime, err = time.Parse(DateFormat, v)
				if err != nil {
					return errors.New("invalid datetime format")
				}
			}
		}
		cd.Time = parsedTime
		return nil
	default:
		return errors.New("unsupported type for CustomDate")
	}
}

// Implement the encoding.TextUnmarshaler interface
func (cd *CustomDate) UnmarshalText(text []byte) error {
	strInput := strings.Trim(string(text), "\"")
	if strInput == "" {
		cd.Time = time.Time{}
		return nil
	}

	// Thử parse với RFC3339
	parsedTime, err := time.Parse(DateTimeFormatRFC3339, strInput)
	if err != nil {
		// Nếu không thành công, thử parse với định dạng ngày duy nhất
		parsedTime, err = time.Parse(DateFormat, strInput)
		if err != nil {
			return fmt.Errorf("invalid datetime format, use RFC3339 or 'YYYY-MM-DD': %w", err)
		}
	}
	cd.Time = parsedTime
	return nil
}
