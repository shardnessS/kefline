// Package odds — упражнение 01: ошибки в Go.
package odds

import (
	"errors"
	"fmt"
)

// ErrEmpty возвращается, когда строка с коэффициентом пустая.
var ErrEmpty = errors.New("пустой коэффициент")

// RangeError — коэффициент вне допустимого диапазона (1; 1000].
type RangeError struct {
	Value float64
}

func (e *RangeError) Error() string {
	return fmt.Sprintf("коэффициент %v вне диапазона (1; 1000]", e.Value)
}

// ParseDecimal разбирает десятичный коэффициент вида "1.85" или "2,10".
func ParseDecimal(s string) (float64, error) {
	panic("не решено: learn/01-errors")
}
