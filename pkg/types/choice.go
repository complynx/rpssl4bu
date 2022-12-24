package types

import (
	"encoding/json"
	"fmt"
)

type Choice byte

const (
	Rock Choice = iota
	Paper
	Scissors
	Lizard
	Spock

	Undefined Choice = 0xff
)

var choiceToString = map[Choice]string{
	Rock:     "rock",
	Paper:    "paper",
	Scissors: "scissors",
	Lizard:   "lizard",
	Spock:    "spock",
}

var stringToChoice = map[string]Choice{
	"rock":     Rock,
	"paper":    Paper,
	"scissors": Scissors,
	"lizard":   Lizard,
	"spock":    Spock,
}

func (r Choice) String() string {
	return choiceToString[r]
}

type jsonChoice struct {
	// ID is the unique identifier of the choice.
	ID int `json:"id"`
	// Name is the name of the choice.
	Name string `json:"name"`
}

func (r Choice) MarshalJSON() ([]byte, error) {
	return json.Marshal(jsonChoice{
		ID:   r.Int(),
		Name: r.String(),
	})
}

func (r *Choice) UnmarshalJSON(data []byte) error {
	// Try unmarshalling the JSON data to a struct.
	var jc jsonChoice
	if err := json.Unmarshal(data, &jc); err == nil {
		// If successful, set the Choice value to the ID from the struct.
		*r, err = IntToChoiceErr(jc.ID)
		return err
	}

	// Try unmarshalling the JSON data to a string.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// If successful, set the Choice value using the string.
		result, ok := stringToChoice[s]
		if !ok {
			return fmt.Errorf("invalid result: %s", s)
		}
		*r = result
		return nil
	}

	// Try unmarshalling the JSON data to an int.
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		// If successful, set the Choice value using the int.
		*r, err = IntToChoiceErr(i)
		return err
	}

	return fmt.Errorf("failed to unmarshal JSON data to struct, string, or int")
}

func (r Choice) Int() int {
	return int(r)
}

func IntToChoice(i int) Choice {
	c, err := IntToChoiceErr(i)
	if err != nil {
		return Undefined
	}
	return c
}

func IntToChoiceErr(i int) (Choice, error) {
	if i < 0 || i > 4 {
		return Undefined, fmt.Errorf("wrong choice ID")
	}
	return Choice(i), nil
}
