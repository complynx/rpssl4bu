package types

import (
	"bytes"
	"fmt"
	"testing"
)

func TestChoiceString(t *testing.T) {
	tests := []struct {
		choice Choice
		want   string
	}{
		{Rock, "rock"},
		{Paper, "paper"},
		{Scissors, "scissors"},
		{Lizard, "lizard"},
		{Spock, "spock"},
		{Choice(123), ""},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.choice), func(t *testing.T) {
			got := test.choice.String()
			if got != test.want {
				t.Errorf("Choice.String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestChoiceFromInt(t *testing.T) {
	tests := []struct {
		choice Choice
		i      int
	}{
		{Rock, 0},
		{Paper, 1},
		{Scissors, 2},
		{Lizard, 3},
		{Spock, 4},
		{Spock, 5},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Choice from int %d", test.i), func(t *testing.T) {
			got := IntToChoice(test.i)
			if got != test.choice {
				t.Errorf("IntToChoice() = %q, want %q", got, test.choice)
			}
		})
	}
}

func TestChoiceFromIntErr(t *testing.T) {
	tests := []struct {
		choice Choice
		i      int
		err    bool
	}{
		{Rock, 0, false},
		{Paper, 1, false},
		{Scissors, 2, false},
		{Lizard, 3, false},
		{Spock, 4, false},
		{Spock, 5, true},
		{Spock, -1, true},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Choice from int %d with err", test.i), func(t *testing.T) {
			got, err := IntToChoiceErr(test.i)
			if err != nil {
				if !test.err {
					t.Errorf("IntToChoiceErr() returned unexpected err %q", err)
				}
			} else if got != test.choice {
				t.Errorf("IntToChoice() = %q, want %q", got, test.choice)
			}
		})
	}
}

func TestChoiceMarshalJSON(t *testing.T) {
	tests := []struct {
		choice Choice
		want   []byte
		err    string
	}{
		{Rock, []byte(`{"id":0,"name":"rock"}`), ""},
		{Paper, []byte(`{"id":1,"name":"paper"}`), ""},
		{Scissors, []byte(`{"id":2,"name":"scissors"}`), ""},
		{Lizard, []byte(`{"id":3,"name":"lizard"}`), ""},
		{Spock, []byte(`{"id":4,"name":"spock"}`), ""},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.choice), func(t *testing.T) {
			got, err := test.choice.MarshalJSON()
			if err != nil {
				if test.err == "" {
					t.Errorf("Choice.MarshalJSON() returned unexpected error: %v", err)
				} else if err.Error() != test.err {
					t.Errorf("Choice.MarshalJSON() returned error %q, want %q", err.Error(), test.err)
				}
				return
			}
			if !bytes.Equal(got, test.want) {
				t.Errorf("Choice.MarshalJSON() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestChoiceUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Choice
		wantErr bool
	}{
		{
			name: "unmarshal from struct",
			data: []byte(`{"id":0,"name":"rock"}`),
			want: Rock,
		},
		{
			name:    "unmarshal from undefined",
			data:    []byte(`undefined`),
			wantErr: true,
		},
		{
			name:    "unmarshal from struct bad",
			data:    []byte(`{"id":6,"name":"rock"}`),
			wantErr: true,
		},
		{
			name: "unmarshal from string",
			data: []byte(`"rock"`),
			want: Rock,
		},
		{
			name: "unmarshal from int",
			data: []byte(`0`),
			want: Rock,
		},
		{
			name: "unmarshal from int",
			data: []byte(`1`),
			want: Paper,
		},
		{
			name: "unmarshal from int",
			data: []byte(`2`),
			want: Scissors,
		},
		{
			name: "unmarshal from int",
			data: []byte(`3`),
			want: Lizard,
		},
		{
			name: "unmarshal from int",
			data: []byte(`4`),
			want: Spock,
		},
		{
			name:    "unmarshal from invalid string",
			data:    []byte(`"invalid"`),
			wantErr: true,
		},
		{
			name:    "unmarshal from invalid int",
			data:    []byte(`5`),
			wantErr: true,
		},
		{
			name:    "unmarshal from invalid int",
			data:    []byte(`-1`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Choice
			err := c.UnmarshalJSON(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && c != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", c, tt.want)
			}
		})
	}
}
