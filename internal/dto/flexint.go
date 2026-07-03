package dto

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexInt accepts a JSON number or a numeric JSON string, since some clients
// send integer fields (e.g. venue type) as strings.
type FlexInt int

// UnmarshalJSON accepts a JSON number or a numeric JSON string.
func (n *FlexInt) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*n = FlexInt(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("must be a number or a numeric string")
	}
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("must be a number or a numeric string")
	}
	*n = FlexInt(parsed)
	return nil
}

// MarshalJSON encodes the value as a JSON number.
func (n FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(n))
}
