# go-future - A (mutex/ condition based) future implementation in go.

## Usage
```
package main

import (
	"fmt"
	"github.com/cbuschka/go-future"
	"time"
)

func main() {

	f := future.NewFuture()

	go func() {
		time.Sleep(1 * time.Second)
		f.MustResolve("resolved")
	}()

	value, err := f.Await()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value is '%s'.", value)
}
```

## License
Copyright (c) 2021 by [Cornelius Buschka](https://github.com/cbuschka).

[MIT](./license.txt)
