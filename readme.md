# go-future - A (mutext based) future implementation in go.

## Usage
```
import "github.com/cbuschka/go-future"

[...]

f := NewFuture()

go func() {
  time.Sleep(1 * time.Millisecond)
  f.Resolve("resolved")
}()

value, err := f.Await()
if err != nil {
  panic(err)
}

fmt.Printf("Value is '%s'.\n", value)
```

## License
Copyright (c) 2021 by [Cornelius Buschka](https://github.com/cbuschka).

[MIT](./license.txt)
