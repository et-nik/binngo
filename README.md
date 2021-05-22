BinnGo
=======================

Binary serializer. Implements easy to use encoding and decoding of Binn. This package is very similar to the standard go
encoder packages like `encoding/json`. BinnGo uses reflection.

Original C Binn Library: https://github.com/liteserver/binn

Binn Specification: https://github.com/liteserver/binn/blob/master/spec.md

## Work In Progress notification

This package in under development. Encoding and decoding complex and nested structures may not work as expected.

## Installation

Run the following command to install the package:

```
go get -u github.com/et-nik/binngo
```

## How To Use

### Reading Binn data

```go
package main

import (
	"fmt"
	"github.com/et-nik/binngo"
)

func main() {
	binnBinary := []byte{
		0xE0,                          // [type] list (container)
		23,                            // [size] container total size
		0x02,                          // [count] items
		0xA0,                          // [type] = string
		0x05,                          // [size]
		'h', 'e', 'l', 'l', 'o', 0x00, // [data] (null terminated)
		0xA0,                          // [type] = string
		0x05,                          // [size]
		'w', 'o', 'r', 'l', 'd', 0x00, // [data] (null terminated)
	}
	items := []string{}

	err := binngo.Unmarshal(binnBinary, &items)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%+v", items)
}
```

### Writing Binn data

```go
package main

import (
	"fmt"
	"io/ioutil"
	"github.com/et-nik/binngo"
)

func main() {
	structure := struct {
		Val1 int64
		Val2 string
	}{
		123,
		"value",
	}

	binnBinary, err := binngo.Marshal(structure)
	if err != nil {
		fmt.Println("error:", err)
	}

	err = ioutil.WriteFile("/path/to/binfile", binnBinary, 0644)
	if err != nil {
		fmt.Println("error:", err)
	}
}
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
