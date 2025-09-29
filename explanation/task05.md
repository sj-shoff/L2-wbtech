# task05

Что выведет программа?

Объяснить вывод программы.

```go
package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  // ... do something
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
    println("error")
    return
  }
  println("ok")
}
```

Вывод:
```
error
```

Объяснение:

Тип error является интерфейсом - он содержит информацию о типе и значении.

Функция test() возвращает *customError (указатель), и в данном случае возвращается nil.

Однако при присваивании err = test() переменная err типа error получает:

Динамический тип: *customError

Динамическое значение: nil

В условии if err != nil проверяется, является ли интерфейс error пустым. Интерфейс считается пустым только если и тип, и значение равны nil.

В данном случае:

Тип: *customError (не nil)

Значение: nil

Поэтому интерфейс error не пустой, и условие if err != nil выполняется.