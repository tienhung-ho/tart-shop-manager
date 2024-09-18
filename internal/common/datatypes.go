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
	StatusActive   Status = "Active"
	StatusInactive Status = "Inactive"
	StatusPending  Status = "Pending"
)

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

type CustomDate struct {
	time.Time
}

const dateFormat = "2006-01-02"

// Marshal JSON to ensure correct format when returning to client
func (cd *CustomDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", cd.Time.Format(dateFormat))), nil
}

// Unmarshal JSON to bind the incoming date string to time.Time
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	strInput := strings.Trim(string(b), "\"")
	parsedTime, err := time.Parse(dateFormat, strInput)
	if err != nil {
		return errors.New("invalid date format, use YYYY-MM-DD")
	}
	cd.Time = parsedTime
	return nil
}

// Implement the driver.Valuer interface for database serialization
func (cd CustomDate) Value() (driver.Value, error) {
	return cd.Time.Format(dateFormat), nil
}

// Implement the sql.Scanner interface for database deserialization
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		*cd = CustomDate{Time: time.Time{}}
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan date field")
	}

	parsedTime, err := time.Parse(dateFormat, str)
	if err != nil {
		return err
	}

	cd.Time = parsedTime
	return nil
}
