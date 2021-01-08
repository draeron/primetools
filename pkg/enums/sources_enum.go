// Code generated by go-enum
// DO NOT EDIT!

package enums

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

const (
	// ITunes is a SourceType of type ITunes
	ITunes SourceType = iota
	// PRIME is a SourceType of type PRIME
	PRIME
	// File is a SourceType of type File
	File
	// Traktor is a SourceType of type Traktor
	Traktor
)

const _SourceTypeName = "ITunesPRIMEFileTraktor"

var _SourceTypeNames = []string{
	_SourceTypeName[0:6],
	_SourceTypeName[6:11],
	_SourceTypeName[11:15],
	_SourceTypeName[15:22],
}

// SourceTypeNames returns a list of possible string values of SourceType.
func SourceTypeNames() []string {
	tmp := make([]string, len(_SourceTypeNames))
	copy(tmp, _SourceTypeNames)
	return tmp
}

var _SourceTypeMap = map[SourceType]string{
	0: _SourceTypeName[0:6],
	1: _SourceTypeName[6:11],
	2: _SourceTypeName[11:15],
	3: _SourceTypeName[15:22],
}

// String implements the Stringer interface.
func (s SourceType) String() string {
	if str, ok := _SourceTypeMap[s]; ok {
		return str
	}
	return fmt.Sprintf("SourceType(%d)", s)
}

var _SourceTypeValue = map[string]SourceType{
	_SourceTypeName[0:6]:                    0,
	strings.ToLower(_SourceTypeName[0:6]):   0,
	_SourceTypeName[6:11]:                   1,
	strings.ToLower(_SourceTypeName[6:11]):  1,
	_SourceTypeName[11:15]:                  2,
	strings.ToLower(_SourceTypeName[11:15]): 2,
	_SourceTypeName[15:22]:                  3,
	strings.ToLower(_SourceTypeName[15:22]): 3,
}

// ParseSourceType attempts to convert a string to a SourceType
func ParseSourceType(name string) (SourceType, error) {
	if x, ok := _SourceTypeValue[name]; ok {
		return x, nil
	}
	return SourceType(0), fmt.Errorf("%s is not a valid SourceType, try [%s]", name, strings.Join(_SourceTypeNames, ", "))
}

// MarshalText implements the text marshaller method
func (s SourceType) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements the text unmarshaller method
func (s *SourceType) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseSourceType(name)
	if err != nil {
		return err
	}
	*s = tmp
	return nil
}

// Scan implements the Scanner interface.
func (s *SourceType) Scan(value interface{}) error {
	var name string

	switch v := value.(type) {
	case string:
		name = v
	case []byte:
		name = string(v)
	case nil:
		*s = SourceType(0)
		return nil
	}

	tmp, err := ParseSourceType(name)
	if err != nil {
		return err
	}
	*s = tmp
	return nil
}

// Value implements the driver Valuer interface.
func (s SourceType) Value() (driver.Value, error) {
	return s.String(), nil
}
