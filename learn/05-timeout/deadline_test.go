package deadline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestFetchOddsFast(t *testing.T) {
	got, err := FetchOdds(context.Background(), 100*time.Millisecond, func(ctx context.Context) (float64, error) {
		return 1.95, nil
	})
	if err != nil || got != 1.95 {
		t.Errorf("получили (%v, %v), ждали (1.95, nil)", got, err)
	}
}

func TestFetchOddsPassesErrorThrough(t *testing.T) {
	errProvider := errors.New("провайдер недоступен")
	_, err := FetchOdds(context.Background(), 100*time.Millisecond, func(ctx context.Context) (float64, error) {
		return 0, errProvider
	})
	if !errors.Is(err, errProvider) {
		t.Errorf("ждали ошибку провайдера, получили %v", err)
	}
}

func TestFetchOddsTimeoutRespectingFetch(t *testing.T) {
	start := time.Now()
	_, err := FetchOdds(context.Background(), 50*time.Millisecond, func(ctx context.Context) (float64, error) {
		select {
		case <-time.After(time.Second):
			return 2.0, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("ждали context.DeadlineExceeded, получили %v", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("вернулись через %v, таймаут был 50ms", elapsed)
	}
}

func TestFetchOddsTimeoutIgnoringFetch(t *testing.T) {
	start := time.Now()
	_, err := FetchOdds(context.Background(), 50*time.Millisecond, func(ctx context.Context) (float64, error) {
		time.Sleep(500 * time.Millisecond) // плохой fetch: контекст не смотрит
		return 2.0, nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("ждали context.DeadlineExceeded, получили %v", err)
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("вернулись через %v: FetchOdds ждал медленный fetch вместо таймаута", elapsed)
	}
}

func TestFetchOddsParentCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	_, err := FetchOdds(ctx, time.Second, func(ctx context.Context) (float64, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("ждали context.Canceled, получили %v", err)
	}
}
