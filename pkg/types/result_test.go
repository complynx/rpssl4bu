package types

import (
	"bytes"
	"fmt"
	"testing"
)

func TestResultString(t *testing.T) {
	tests := []struct {
		res  Result
		want string
	}{
		{Tie, "tie"},
		{Win, "win"},
		{Lose, "lose"},
		{Result(123), ""},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.res), func(t *testing.T) {
			got := test.res.String()
			if got != test.want {
				t.Errorf("Result.String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestResultSwap(t *testing.T) {
	tests := []struct {
		res  Result
		want Result
	}{
		{Tie, Tie},
		{Win, Lose},
		{Lose, Win},
		{Result(123), Result(123)},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.res), func(t *testing.T) {
			got := test.res.Swap()
			if got != test.want {
				t.Errorf("Result.Swap() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestResultFromInt(t *testing.T) {
	tests := []struct {
		res Result
		i   int
	}{
		{Tie, 2},
		{Lose, 1},
		{Win, 0},
		{Unknown, 3},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Result from int %d", test.i), func(t *testing.T) {
			got := IntToResult(test.i)
			if got != test.res {
				t.Errorf("IntToResult() = %q, want %q", got, test.res)
			}
		})
	}
}

func TestResultFromIntErr(t *testing.T) {
	tests := []struct {
		res Result
		i   int
		err bool
	}{
		{Win, 0, false},
		{Lose, 1, false},
		{Tie, 2, false},
		{Tie, 3, true},
		{Tie, -1, true},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Result from int %d with err", test.i), func(t *testing.T) {
			got, err := IntToResultErr(test.i)
			if err != nil {
				if !test.err {
					t.Errorf("IntToResultErr() returned unexpected err %q", err)
				}
			} else if got != test.res {
				t.Errorf("IntToResult() = %q, want %q", got, test.res)
			}
		})
	}
}

func TestResultMarshalJSON(t *testing.T) {
	tests := []struct {
		res  Result
		want []byte
		err  string
	}{
		{Win, []byte(`"win"`), ""},
		{Lose, []byte(`"lose"`), ""},
		{Tie, []byte(`"tie"`), ""},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.res), func(t *testing.T) {
			got, err := test.res.MarshalJSON()
			if err != nil {
				if test.err == "" {
					t.Errorf("Result.MarshalJSON() returned unexpected error: %v", err)
				} else if err.Error() != test.err {
					t.Errorf("Result.MarshalJSON() returned error %q, want %q", err.Error(), test.err)
				}
				return
			}
			if !bytes.Equal(got, test.want) {
				t.Errorf("Result.MarshalJSON() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestResultUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    Result
		wantErr bool
	}{
		{
			name:    "unmarshal from undefined",
			data:    []byte(`undefined`),
			wantErr: true,
		},
		{
			name: "unmarshal from string",
			data: []byte(`"win"`),
			want: Win,
		},
		{
			name: "unmarshal from int",
			data: []byte(`0`),
			want: Win,
		},
		{
			name: "unmarshal from int",
			data: []byte(`1`),
			want: Lose,
		},
		{
			name: "unmarshal from int",
			data: []byte(`2`),
			want: Tie,
		},
		{
			name:    "unmarshal from invalid string",
			data:    []byte(`"invalid"`),
			wantErr: true,
		},
		{
			name:    "unmarshal from invalid int",
			data:    []byte(`3`),
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
			var c Result
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
