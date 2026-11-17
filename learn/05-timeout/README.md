# 05. Таймауты и отмена

В Java отмену делают через `Future.cancel` и прерывание потока. В Go всё, что может долго работать, принимает первым параметром `ctx context.Context`. Через него приходят таймаут и отмена.

Задача: `FetchOdds` вызывает `fetch` с контекстом, который истекает через `timeout`.
- `fetch` успел: вернуть его результат.
- Не успел: вернуть ошибку, по которой `errors.Is(err, context.DeadlineExceeded)` даёт `true`.
- Родительский `ctx` отменили: `errors.Is(err, context.Canceled)`.
- **Даже если `fetch` не смотрит на контекст и спит дольше**, `FetchOdds` должен вернуться по таймауту, а не ждать его.

Подсказки:
- `ctx, cancel := context.WithTimeout(parent, timeout)` и сразу `defer cancel()`.
- Запусти `fetch` в горутине, результат положи в канал, ожидай через `select` либо канал, либо `ctx.Done()`.
- Канал результата сделай с буфером 1. Иначе горутина с медленным `fetch` после таймаута навсегда зависнет на отправке: это утечка горутины.

Что почитать: [Go by Example: Context](https://gobyexample.com/context), [Timeouts](https://gobyexample.com/timeouts), [Go Concurrency Patterns: Context](https://go.dev/blog/context).
