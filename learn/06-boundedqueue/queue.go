// Package queue — упражнение 06: очередь ограниченного размера без блокировки отправителя.
package queue

import "context"

// Queue — очередь целых чисел ограниченного размера.
type Queue struct {
	// TODO: поля
}

// New создаёт очередь ёмкостью capacity.
func New(capacity int) *Queue {
	panic("не решено: learn/06-boundedqueue")
}

// TryPush кладёт v без ожидания. Возвращает false, если очередь полна и значение выброшено.
func (q *Queue) TryPush(v int) bool {
	panic("не решено: learn/06-boundedqueue")
}

// Pop ждёт значение. ok == false, если ctx закончился раньше.
func (q *Queue) Pop(ctx context.Context) (v int, ok bool) {
	panic("не решено: learn/06-boundedqueue")
}

// Len возвращает, сколько значений лежит в очереди сейчас.
func (q *Queue) Len() int {
	panic("не решено: learn/06-boundedqueue")
}

// Dropped возвращает, сколько значений выброшено из-за переполнения.
func (q *Queue) Dropped() int64 {
	panic("не решено: learn/06-boundedqueue")
}
