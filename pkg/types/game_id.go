package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type GameID uint64

func (r GameID) String() string {
	return fmt.Sprintf("%016x", uint64(r))
}

func (r GameID) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func GameIDFromString(s string) (GameID, error) {
	if len(s) != 16 {
		return 0, fmt.Errorf("wrong string length for GameID")
	}

	// Parse the string as a hexadecimal number.
	parsed, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("parse uint: %w", err)
	}

	return GameID(parsed), nil
}

func (r *GameID) UnmarshalJSON(data []byte) error {
	// Try unmarshalling the JSON data to a string.
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	// Parse the string as a hexadecimal number.
	parsed, err := GameIDFromString(s)
	if err != nil {
		return err
	}

	*r = parsed
	return nil
}
