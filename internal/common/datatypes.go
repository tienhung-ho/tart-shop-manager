package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
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
