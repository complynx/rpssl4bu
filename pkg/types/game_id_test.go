package types

import (
	"bytes"
	"fmt"
	"testing"
)

func TestGIString(t *testing.T) {
	tests := []struct {
		gi   GameID
		want string
	}{
		{0x908, "0000000000000908"},
		{0xc33, "0000000000000c33"},
		{0x441, "0000000000000441"},
		{0xffffffffffffffff, "ffffffffffffffff"},
		{0, "0000000000000000"},
		{0x1af610fb9fe19ab, "01af610fb9fe19ab"},
		{0xa71ea4ac49b6fc7b, "a71ea4ac49b6fc7b"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.gi), func(t *testing.T) {
			got := test.gi.String()
			if got != test.want {
				t.Errorf("GameID.String() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGIMarshalJSON(t *testing.T) {
	tests := []struct {
		gi   GameID
		want []byte
	}{
		{0x1af610fb9fe19ab, []byte(`"01af610fb9fe19ab"`)},
		{0xa71ea4ac49b6fc7b, []byte(`"a71ea4ac49b6fc7b"`)},
		{0, []byte(`"0000000000000000"`)},
		{0xffffffffffffffff, []byte(`"ffffffffffffffff"`)},
		{0xc33, []byte(`"0000000000000c33"`)},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.gi), func(t *testing.T) {
			got, err := test.gi.MarshalJSON()
			if err != nil {
				t.Errorf("Choice.MarshalJSON() returned error %q", err.Error())
				return
			}
			if !bytes.Equal(got, test.want) {
				t.Errorf("Choice.MarshalJSON() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGIUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    GameID
		wantErr bool
	}{
		{
			name:    "unmarshal from struct",
			data:    []byte(`{"id":0,"name":"rock"}`),
			wantErr: true,
		},
		{
			name:    "unmarshal from undefined",
			data:    []byte(`undefined`),
			wantErr: true,
		},
		{
			name:    "unmarshal from int",
			data:    []byte(`34`),
			wantErr: true,
		},
		{
			name: "unmarshal from string ffffffffffffffff",
			data: []byte(`"ffffffffffffffff"`),
			want: 0xffffffffffffffff,
		},
		{
			name: "unmarshal from string 0000000000000000",
			data: []byte(`"0000000000000000"`),
			want: 0x0,
		},
		{
			name: "unmarshal from string a71ea4ac49b6fc7b",
			data: []byte(`"a71ea4ac49b6fc7b"`),
			want: 0xa71ea4ac49b6fc7b,
		},
		{
			name: "unmarshal from string 01af610fb9fe19ab",
			data: []byte(`"01af610fb9fe19ab"`),
			want: 0x1af610fb9fe19ab,
		},
		{
			name: "unmarshal from string 0000000000000c33",
			data: []byte(`"0000000000000c33"`),
			want: 0x0000000000000c33,
		},
		{
			name:    "unmarshal from string 000000000000c33 (too few zeros)",
			data:    []byte(`"000000000000c33"`),
			wantErr: true,
		},
		{
			name:    "unmarshal from string c33 (too few zeros)",
			data:    []byte(`"c33"`),
			wantErr: true,
		},
		{
			name:    "unmarshal from string 00000000000000c33 (too many zeros)",
			data:    []byte(`"00000000000000c33"`),
			wantErr: true,
		},
		{
			name:    "unmarshal from string a71ra4ac49b6fc7b (unexpected character)",
			data:    []byte(`"a71ra4ac49b6fc7b"`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c GameID
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
