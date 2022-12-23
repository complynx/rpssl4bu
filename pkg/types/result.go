package types

import (
	"encoding/json"
	"fmt"
)

type Result byte

const (
	Win Result = iota
	Lose
	Tie
)

var resultToString = map[Result]string{
	Win:  "win",
	Lose: "lose",
	Tie:  "tie",
}

var stringToResult = map[string]Result{
	"win":  Win,
	"lose": Lose,
	"tie":  Tie,
}

func (r Result) String() string {
	return resultToString[r]
}

func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Result) UnmarshalJSON(data []byte) error {
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*r, err = IntToResultErr(i)
		return err
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	result, ok := stringToResult[s]
	if !ok {
		return fmt.Errorf("invalid result: %s", s)
	}
	*r = result
	return nil
}

func (r Result) Int() int {
	return int(r)
}

func IntToResult(i int) Result {
	if i < 0 || i > 2 {
		return Tie
	}
	return Result(i)
}

func IntToResultErr(i int) (Result, error) {
	if i < 0 || i > 2 {
		return Tie, fmt.Errorf("wrong result ID")
	}
	return Result(i), nil
}
