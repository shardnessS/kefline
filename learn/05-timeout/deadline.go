// Package deadline — упражнение 05: context, таймауты и отмена.
package deadline

import (
	"context"
	"time"
)

// FetchOdds вызывает fetch с контекстом, который истекает через timeout,
// и возвращается не позже таймаута, даже если fetch контекст не проверяет.
func FetchOdds(ctx context.Context, timeout time.Duration, fetch func(ctx context.Context) (float64, error)) (float64, error) {
	panic("не решено: learn/05-timeout")
}
