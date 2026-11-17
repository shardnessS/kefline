package conflate

import (
	"fmt"
	"sync"
	"testing"
)

func TestLatestWinsKeepingPosition(t *testing.T) {
	c := New()
	c.Put(Change{Key: "П1", Value: 1.80, Version: 1})
	c.Put(Change{Key: "X", Value: 3.40, Version: 1})
	c.Put(Change{Key: "П1", Value: 1.75, Version: 2})
	c.Put(Change{Key: "П2", Value: 4.10, Version: 1})
	c.Put(Change{Key: "П1", Value: 1.70, Version: 3})

	got := c.Drain()
	want := []Change{
		{Key: "П1", Value: 1.70, Version: 3},
		{Key: "X", Value: 3.40, Version: 1},
		{Key: "П2", Value: 4.10, Version: 1},
	}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("Drain() = %v, ждали %v", got, want)
	}
	if m := c.Merged(); m != 2 {
		t.Errorf("Merged() = %d, ждали 2", m)
	}
}

func TestDrainClears(t *testing.T) {
	c := New()
	c.Put(Change{Key: "a", Value: 1, Version: 1})
	first := c.Drain()
	if len(first) != 1 {
		t.Fatalf("первый Drain: ждали 1 изменение, получили %v", first)
	}
	if again := c.Drain(); len(again) != 0 {
		t.Errorf("второй Drain должен быть пустым, получили %v", again)
	}
	c.Put(Change{Key: "b", Value: 2, Version: 1})
	if got := c.Drain(); len(got) != 1 || got[0].Key != "b" {
		t.Errorf("после очистки ждали только b, получили %v", got)
	}
	if first[0].Key != "a" {
		t.Errorf("Drain отдал наружу внутренний срез: старый результат изменился на %v", first)
	}
}

func TestIgnoresLateChange(t *testing.T) {
	c := New()
	c.Put(Change{Key: "П1", Value: 1.70, Version: 5})
	c.Put(Change{Key: "П1", Value: 1.90, Version: 3}) // опоздало

	got := c.Drain()
	if len(got) != 1 || got[0].Version != 5 || got[0].Value != 1.70 {
		t.Errorf("опоздавшее изменение должно игнорироваться, получили %v", got)
	}
}

func TestConcurrentPut(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for v := 1; v <= 100; v++ {
				c.Put(Change{Key: fmt.Sprintf("k%d", v%10), Value: float64(v), Version: int64(v)})
			}
		}(g)
	}
	wg.Wait()

	got := c.Drain()
	if len(got) != 10 {
		t.Fatalf("ждали 10 ключей, получили %d", len(got))
	}
	for _, ch := range got {
		// у каждого ключа последняя версия: 91..100
		if ch.Version < 91 {
			t.Errorf("ключ %s: версия %d, ждали последнюю из 91..100", ch.Key, ch.Version)
		}
	}
}
