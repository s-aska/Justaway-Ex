package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// NullUint64 is a sql.Scanner for unsigned ints.
type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

// Scan implements the sql.Scanner interface.
func (n *NullUint64) Scan(src interface{}) error {
	if src == nil {
		n.Uint64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	s := asString(src)
	var err error
	n.Uint64, err = strconv.ParseUint(s, 10, 64)
	return err
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	return fmt.Sprintf("%v", src)
}

type JsonNullUInt64 struct {
	NullUint64
}

func (v JsonNullUInt64) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Uint64)
	} else {
		return json.Marshal(nil)
	}
}
