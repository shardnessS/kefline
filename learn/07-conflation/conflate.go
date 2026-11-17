// Package conflate — упражнение 07: склейка изменений, побеждает последнее.
package conflate

// Change — изменение значения по ключу, например коэффициента исхода.
type Change struct {
	Key     string
	Value   float64
	Version int64
}

// Conflator хранит по одному ожидающему изменению на ключ.
type Conflator struct {
	// TODO: поля
}

// New создаёт пустую склейку.
func New() *Conflator {
	panic("не решено: learn/07-conflation")
}

// Put добавляет изменение. Ожидающее изменение с тем же ключом заменяется,
// если версия нового не меньше. Позиция ключа в очереди не меняется.
func (c *Conflator) Put(ch Change) {
	panic("не решено: learn/07-conflation")
}

// Drain возвращает ожидающие изменения в порядке первого появления ключей и очищает склейку.
func (c *Conflator) Drain() []Change {
	panic("не решено: learn/07-conflation")
}

// Merged возвращает, сколько изменений было заменено более новыми.
func (c *Conflator) Merged() int64 {
	panic("не решено: learn/07-conflation")
}
