// Package null provides a convenient way of handling null values.
// Types in this package consider empty or zero input the same as null input.
// Types in this package will encode to their zero value, even if null.
// Use the nuller subpackage if you don't want this.
package null

import (
	"database/sql"
	"encoding/json"
)

// String is a nullable string.
type String struct {
	sql.NullString
}

// NewString creates a new String
func NewString(s string, valid bool) String {
	return String{
		NullString: sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// StringFrom creates a new String that will be null if s is blank.
func StringFrom(s string) String {
	return NewString(s, s != "")
}

// StringFrom creates a new String that be null if s is nil or blank.
// It will make s point to the String's value.
func StringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}
	str := NewString(*s, *s != "")
	return str
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input produces a null String.
// It also supports unmarshalling a sql.NullString.
func (s *String) UnmarshalJSON(data []byte) error {
	var err error
	var v interface{}
	json.Unmarshal(data, &v)
	switch x := v.(type) {
	case string:
		s.String = x
	case map[string]interface{}:
		err = json.Unmarshal(data, &s.NullString)
	case nil:
		s.Valid = false
		return nil
	}
	s.Valid = (err == nil) && (s.String != "")
	return err
}

// MarshalText implements encoding.TextMarshaler.
// It will encode a blank string when this String is null.
func (s String) MarshalText() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return []byte(s.String), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null String if the input is a blank string.
func (s *String) UnmarshalText(text []byte) error {
	s.String = string(text)
	s.Valid = s.String != ""
	return nil
}

// Pointer returns a pointer to this String's value, or a nil pointer if this String is null.
func (s String) Ptr() *string {
	if !s.Valid {
		return nil
	}
	return &s.String
}

// IsZero returns true for null or empty strings, for future omitempty support. (Go 1.4?)
func (s String) IsZero() bool {
	return !s.Valid || s.String == ""
}
