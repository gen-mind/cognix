package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Map represents a key-value map where the keys are strings and the values can be of any type.
type Map map[string]interface{}

// Value returns the JSON representation of the Map as a driver.Value.
func (m Map) Value() (driver.Value, error) {
	j, err := json.Marshal(m)
	return string(j), err
}

// Scan reads the value from the database driver into the Map instance.
// It takes the value as an interface{} and attempts a type assertion to []byte.
// If the type assertion fails, it returns an error.
// Otherwise, it uses json.Unmarshal to decode the byte slice into the Map instance.
// It returns nil if the decoding is successful, otherwise it returns an error.
func (m *Map) Scan(val interface{}) error {
	value, ok := val.([]byte)
	if !ok {
		return fmt.Errorf("Type assertion .([]byte) failed.")
	}
	return json.Unmarshal(value, m)
}

// ToStruct converts a Map to a struct by marshaling the map to JSON and then unmarshaling it into the provided struct.
func (m *Map) ToStruct(s interface{}) error {
	buf, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, s)
}
