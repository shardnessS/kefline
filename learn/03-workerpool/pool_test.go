package pool

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestMapKeepsOrder(t *testing.T) {
	items := []int{5, 1, 4, 2, 3}
	got := Map(items, 3, func(x int) int {
		time.Sleep(time.Duration(x) * time.Millisecond) // разное время, чтобы порядок завершения перемешался
		return x * 10
	})
	want := []int{50, 10, 40, 20, 30}
	if len(got) != len(want) {
		t.Fatalf("длина %d, ждали %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, ждали %d", i, got[i], want[i])
		}
	}
}

func TestMapLimitsConcurrency(t *testing.T) {
	const workers = 4
	var running, maxRunning atomic.Int64

	items := make([]int, 40)
	Map(items, workers, func(x int) int {
		now := running.Add(1)
		for {
			prev := maxRunning.Load()
			if now <= prev || maxRunning.CompareAndSwap(prev, now) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		running.Add(-1)
		return x
	})

	if m := maxRunning.Load(); m > workers {
		t.Errorf("одновременно работало %d горутин, лимит %d", m, workers)
	}
	if m := maxRunning.Load(); m < 2 {
		t.Errorf("одновременно работала %d горутина: параллелизма нет", m)
	}
}

func TestMapEdgeCases(t *testing.T) {
	if got := Map(nil, 4, func(x int) int { return x }); len(got) != 0 {
		t.Errorf("пустой вход: ждали пустой результат, получили %v", got)
	}
	got := Map([]int{1, 2, 3}, 0, func(x int) int { return x + 1 })
	if len(got) != 3 || got[0] != 2 || got[2] != 4 {
		t.Errorf("workers=0 должно работать как 1, получили %v", got)
	}
}
