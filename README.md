# mmdb

Download MindMax database using license key

## usage

```golang
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	
	"github.com/loafoe/mmdb"
)

func main() {
	file, err := ioutil.TempFile("", "*.mmdb")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = mmdb.Download(file.Name(), os.Getenv("LICENSE_KEY"))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("Done: %s\n", file.Name())
}
```

## license

License is MIT
