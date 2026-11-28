package event

import (
	"encoding/json"
	"testing"
)

func TestParseMoney(t *testing.T) {
	tests := []struct {
		in      string
		want    Money
		wantErr bool
	}{
		{in: "0", want: 0},
		{in: "12", want: 12_000_000},
		{in: "0.5", want: 500_000},
		{in: "12.345600", want: 12_345_600},
		{in: "0.000001", want: 1},
		{in: "1.2345678", wantErr: true},
		{in: "", wantErr: true},
		{in: ".5", wantErr: true},
		{in: "5.", wantErr: true},
		{in: "-1", wantErr: true},
		{in: "+1", wantErr: true},
		{in: "1e3", wantErr: true},
		{in: " 1", wantErr: true},
		{in: "1,5", wantErr: true},
		{in: "99999999999999999999", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseMoney(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMoney(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseMoney(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}

func TestMoneyString(t *testing.T) {
	tests := []struct {
		in   Money
		want string
	}{
		{in: 0, want: "0.000000"},
		{in: 1, want: "0.000001"},
		{in: 12_345_600, want: "12.345600"},
		{in: -1_500_000, want: "-1.500000"},
	}
	for _, tt := range tests {
		if got := tt.in.String(); got != tt.want {
			t.Errorf("Money(%d).String() = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestMoneyJSON(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    Money
		wantErr bool
	}{
		{name: "string", in: `"1.25"`, want: 1_250_000},
		{name: "number rejected", in: `1.25`, wantErr: true},
		{name: "integer number rejected", in: `1`, wantErr: true},
		{name: "bool rejected", in: `true`, wantErr: true},
		{name: "bad string", in: `"abc"`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Money
			err := json.Unmarshal([]byte(tt.in), &m)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal(%s) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if !tt.wantErr && m != tt.want {
				t.Errorf("Unmarshal(%s) = %d, want %d", tt.in, m, tt.want)
			}
		})
	}

	out, err := json.Marshal(Money(12_345_600))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `"12.345600"` {
		t.Errorf("Marshal = %s, want \"12.345600\"", out)
	}
}
