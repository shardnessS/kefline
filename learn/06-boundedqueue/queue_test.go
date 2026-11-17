package queue

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestQueueDropsWhenFull(t *testing.T) {
	q := New(2)
	if !q.TryPush(1) || !q.TryPush(2) {
		t.Fatal("первые два значения должны поместиться")
	}
	if q.TryPush(3) {
		t.Error("третье значение не должно поместиться в очередь ёмкостью 2")
	}
	if got := q.Dropped(); got != 1 {
		t.Errorf("Dropped() = %d, ждали 1", got)
	}
	if got := q.Len(); got != 2 {
		t.Errorf("Len() = %d, ждали 2", got)
	}

	ctx := context.Background()
	for _, want := range []int{1, 2} {
		v, ok := q.Pop(ctx)
		if !ok || v != want {
			t.Errorf("Pop() = (%d, %v), ждали (%d, true)", v, ok, want)
		}
	}
	if !q.TryPush(4) {
		t.Error("после Pop место должно освободиться")
	}
}

func TestQueuePopRespectsContext(t *testing.T) {
	q := New(1)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	start := time.Now()
	_, ok := q.Pop(ctx)
	if ok {
		t.Error("пустая очередь: ждали ok == false после таймаута")
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("Pop вернулся через %v, таймаут был 30ms", elapsed)
	}
}

func TestQueuePopWaitsForValue(t *testing.T) {
	q := New(1)
	go func() {
		time.Sleep(20 * time.Millisecond)
		q.TryPush(42)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	v, ok := q.Pop(ctx)
	if !ok || v != 42 {
		t.Errorf("Pop() = (%d, %v), ждали (42, true)", v, ok)
	}
}

func TestQueueConcurrentPush(t *testing.T) {
	q := New(100)
	var wg sync.WaitGroup
	var accepted sync.Map
	for g := 0; g < 10; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				if q.TryPush(g*1000 + i) {
					accepted.Store(g*1000+i, true)
				}
			}
		}(g)
	}
	wg.Wait()

	if q.Len() != 100 {
		t.Errorf("Len() = %d, ждали 100", q.Len())
	}
	if q.Dropped() != 400 {
		t.Errorf("Dropped() = %d, ждали 400 (500 отправок в очередь на 100)", q.Dropped())
	}
}
