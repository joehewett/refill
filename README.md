# refill

Rehydrate Go structs from arbitrary text-based data using language models.

This is a work in progress, and is not yet usable.

# Usage

```sh
$ go get github.com/joehewett/refill
```

```go
package main

import (
	"fmt"
	"github.com/joehewett/refill"
)

type Person struct {
	Age  int
	Name string
	Address string
	Occupation string
}

func main() {
	text := "My name is Joe, I am 23 years old and by day i'm an engineer."
	person := Person{}
	refill.Refill(&person, text)

	fmt.Printf("%+v", person)
}
```

Output:

```
{Age:23 Name:Joe Address: Occupation:engineer}
```

# Disclaimer

This project is experimental. I have no idea if it will work for your use case. I have no idea if it will work for my use case.

# License

[MIT](LICENSE)
