// Package notify — упражнение 02: интерфейсы.
package notify

// Update — изменение коэффициента одного исхода.
type Update struct {
	MatchID string
	Outcome string
	Odds    float64
}

// Sink — приёмник обновлений: клиент, лог, очередь.
type Sink interface {
	Send(u Update) error
}

// MemorySink запоминает все полученные обновления.
type MemorySink struct {
	Got []Update
}

// TODO: метод Send для *MemorySink.

// FailingSink всегда отвечает ошибкой Err.
type FailingSink struct {
	Err error
}

// TODO: метод Send для FailingSink.

// SendAll отправляет u во все приёмники и возвращает все ошибки через errors.Join.
func SendAll(sinks []Sink, u Update) error {
	panic("не решено: learn/02-interfaces")
}
