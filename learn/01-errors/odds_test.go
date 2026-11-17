package odds

import (
	"errors"
	"strconv"
	"strings"
	"testing"
)

func TestParseDecimalValid(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"1.85", 1.85},
		{" 2,10 ", 2.10},
		{"1000", 1000},
		{"1.01", 1.01},
	}
	for _, c := range cases {
		got, err := ParseDecimal(c.in)
		if err != nil {
			t.Errorf("ParseDecimal(%q): неожиданная ошибка %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseDecimal(%q) = %v, ждали %v", c.in, got, c.want)
		}
	}
}

func TestParseDecimalEmpty(t *testing.T) {
	for _, in := range []string{"", "   "} {
		_, err := ParseDecimal(in)
		if !errors.Is(err, ErrEmpty) {
			t.Errorf("ParseDecimal(%q): ждали ErrEmpty, получили %v", in, err)
		}
	}
}

func TestParseDecimalNotANumber(t *testing.T) {
	_, err := ParseDecimal("abc")
	if err == nil {
		t.Fatal("ждали ошибку для \"abc\"")
	}
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Errorf("ошибка должна оборачивать *strconv.NumError через %%w, получили %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "abc") {
		t.Errorf("в тексте ошибки нет исходной строки: %v", err)
	}
}

func TestParseDecimalOutOfRange(t *testing.T) {
	for _, in := range []string{"1", "0.5", "1000.01", "-3"} {
		_, err := ParseDecimal(in)
		var rangeErr *RangeError
		if !errors.As(err, &rangeErr) {
			t.Errorf("ParseDecimal(%q): ждали *RangeError, получили %v", in, err)
			continue
		}
		want, _ := strconv.ParseFloat(in, 64)
		if rangeErr.Value != want {
			t.Errorf("ParseDecimal(%q): RangeError.Value = %v, ждали %v", in, rangeErr.Value, want)
		}
	}
}
