package fanout

import (
	"sync"
	"testing"
	"time"
)

func TestBroadcastDeliversAllInOrder(t *testing.T) {
	in := make(chan int)
	outs := Broadcast(in, 3, 0)
	if len(outs) != 3 {
		t.Fatalf("ждали 3 выходных канала, получили %d", len(outs))
	}

	go func() {
		for i := 1; i <= 100; i++ {
			in <- i
		}
		close(in)
	}()

	var wg sync.WaitGroup
	results := make([][]int, len(outs))
	for i, out := range outs {
		wg.Add(1)
		go func(i int, out <-chan int) {
			defer wg.Done()
			for v := range out { // цикл закончится, только если канал закрыт
				results[i] = append(results[i], v)
			}
		}(i, out)
	}

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("выходные каналы не закрылись после закрытия входа")
	}

	for i, r := range results {
		if len(r) != 100 {
			t.Errorf("подписчик %d получил %d значений, ждали 100", i, len(r))
			continue
		}
		for j, v := range r {
			if v != j+1 {
				t.Errorf("подписчик %d: на месте %d значение %d, ждали %d", i, j, v, j+1)
				break
			}
		}
	}
}

func TestBroadcastSlowReaderBlocksOthers(t *testing.T) {
	// Это не ошибка решения, а свойство простой раздачи: запомни его до недели 3.
	in := make(chan int, 10)
	outs := Broadcast(in, 2, 1)
	fast, slow := outs[0], outs[1]
	_ = slow // медленный подписчик ничего не читает

	for i := 0; i < 5; i++ {
		in <- i
	}

	received := 0
	timeout := time.After(300 * time.Millisecond)
loop:
	for {
		select {
		case <-fast:
			received++
		case <-timeout:
			break loop
		}
	}
	if received >= 5 {
		t.Errorf("быстрый подписчик получил все %d значений, хотя медленный не читает: так не бывает при синхронной раздаче", received)
	}
	t.Logf("быстрый получил %d из 5, пока медленный спал", received)
}
