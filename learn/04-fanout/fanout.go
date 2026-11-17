// Package fanout — упражнение 04: раздача одного канала многим.
package fanout

// Broadcast отправляет каждое значение из in во все n выходных каналов с буфером buf.
// Когда in закрыт, закрывает все выходные каналы.
func Broadcast(in <-chan int, n, buf int) []<-chan int {
	panic("не решено: learn/04-fanout")
}
