package notify

import (
	"errors"
	"testing"
)

var (
	_ Sink = (*MemorySink)(nil)
	_ Sink = FailingSink{}
)

func TestSendAllDeliversEverywhere(t *testing.T) {
	a, b := &MemorySink{}, &MemorySink{}
	u := Update{MatchID: "m1", Outcome: "П1", Odds: 1.85}

	if err := SendAll([]Sink{a, b}, u); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	for name, s := range map[string]*MemorySink{"a": a, "b": b} {
		if len(s.Got) != 1 || s.Got[0] != u {
			t.Errorf("приёмник %s получил %v, ждали [%v]", name, s.Got, u)
		}
	}
}

func TestSendAllContinuesAfterError(t *testing.T) {
	errBroken := errors.New("клиент отвалился")
	errSlow := errors.New("клиент тормозит")
	first, last := &MemorySink{}, &MemorySink{}
	u := Update{MatchID: "m2", Outcome: "X", Odds: 3.4}

	err := SendAll([]Sink{first, FailingSink{Err: errBroken}, FailingSink{Err: errSlow}, last}, u)

	if !errors.Is(err, errBroken) || !errors.Is(err, errSlow) {
		t.Errorf("ждали обе ошибки в результате, получили %v", err)
	}
	if len(first.Got) != 1 || len(last.Got) != 1 {
		t.Errorf("ошибка одного приёмника не должна мешать остальным: first=%v last=%v", first.Got, last.Got)
	}
}

func TestSendAllEmpty(t *testing.T) {
	if err := SendAll(nil, Update{}); err != nil {
		t.Errorf("без приёмников ждали nil, получили %v", err)
	}
}
