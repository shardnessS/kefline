# Go после Java: шпаргалка

Только то, что понадобится в kefline. Примеры минимальные.

## Ошибки вместо исключений

```go
v, err := strconv.ParseFloat(s, 64)
if err != nil {
    return 0, fmt.Errorf("разбор коэффициента %q: %w", s, err) // %w сохраняет причину
}
```

| Java | Go |
|---|---|
| `throw new X()` | `return ..., err` |
| `catch (X e)` | `errors.As(err, &target)` |
| `e == CONST` | `errors.Is(err, ErrConst)` |
| `getCause()` | `errors.Unwrap(err)`, но обычно хватает `Is` и `As` |
| `RuntimeException` | `panic`, только для «так не бывает» |

## Типы и интерфейсы

```go
type Match struct {          // класс без наследования
    ID      string
    Version int64
}

func (m *Match) Bump() { m.Version++ }        // метод на указателе меняет объект
func (m Match) Key() string { return m.ID }   // метод на значении получает копию

type Sink interface { Send(u Update) error }  // никаких implements
```

- Нет наследования, есть встраивание: `type Market struct { Base; ... }`.
- Интерфейсы объявляют там, где их используют, а не там, где реализуют. Маленькие: один-два метода.
- Конструктор — обычная функция `NewX(...) *X`.
- Нет DI-контейнера. Зависимости передаются в конструктор, собирается всё в `main`.

## Видимость

С большой буквы — видно из других пакетов (`public`), с маленькой — только внутри пакета. Папка `internal/` закрыта для внешних модулей.

## Горутины и каналы

| Java | Go |
|---|---|
| `new Thread(r).start()` | `go f()` |
| `ExecutorService` на N потоков | N горутин читают из канала |
| `BlockingQueue` | канал `make(chan T, n)` |
| `CountDownLatch` | `sync.WaitGroup` |
| `synchronized` | `sync.Mutex` |
| `AtomicLong` | `atomic.Int64` |
| `Future.get(timeout)` | `select` с `ctx.Done()` или `time.After` |
| `Thread.interrupt()` | `context` с отменой |

```go
select {
case v := <-results:     // пришёл результат
    return v, nil
case <-ctx.Done():       // таймаут или отмена
    return 0, ctx.Err()
}
```

Правила:
- Закрывает канал тот, кто в него пишет.
- Чтение из закрытого канала сразу отдаёт нулевое значение, `for v := range ch` заканчивается на закрытии.
- Запись в закрытый канал — паника.
- Горутина, которая навсегда ждёт на канале, — утечка. Ей всегда нужен путь выхода: `ctx.Done()` или буфер.
- `go test -race` ловит гонки. Запускай его всегда, когда есть горутины.

## context

- Первый параметр всего, что может ждать: `func Load(ctx context.Context, id string)`.
- `context.WithTimeout`, `context.WithCancel` и сразу `defer cancel()`.
- В `main` корневой контекст отменяется по Ctrl+C: `signal.NotifyContext(context.Background(), os.Interrupt)`.

## Тесты

```go
func TestApply(t *testing.T) {
    cases := []struct {
        name string
        in   int64
        want bool
    }{
        {"новая версия", 5, true},
        {"старая версия", 3, false},
    }
    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            if got := apply(c.in); got != c.want {
                t.Errorf("apply(%d) = %v, ждали %v", c.in, got, c.want)
            }
        })
    }
}
```

- Тест лежит рядом с кодом: `match.go` и `match_test.go`.
- Нет JUnit и Mockito: стандартный `testing`, моки пишутся руками через интерфейсы.
- Бенчмарк: `func BenchmarkX(b *testing.B)`, запуск `go test -bench=.`.

## Инструменты

| Задача | Команда |
|---|---|
| Сборка всех пакетов | `go build ./...` |
| Тесты | `go test ./...`, с гонками `go test -race ./...` |
| Проверка кода | `go vet ./...`, позже `golangci-lint run` |
| Форматирование | `gofmt -w .` (спорить о стиле не нужно) |
| Зависимость | `go get github.com/jackc/pgx/v5` |
| Убрать лишние зависимости | `go mod tidy` |
| Профиль | `go test -cpuprofile cpu.out`, потом `go tool pprof cpu.out` |

## Стиль

- Короткие имена в узкой области: `i`, `v`, `ctx`, `err`. Понятные имена в широкой.
- Без геттеров `GetX`: поле `X` или метод `X()`.
- Ошибки обрабатываются сразу, без `else` после `return`.
- Комментарий к экспортируемому имени начинается с самого имени: `// Apply применяет изменение...`.
