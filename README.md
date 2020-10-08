# Poller
Poller is a file watching library which will produce the following events:

- CREATE
- WRITE
- CHMOD
- REMOVE


## Usage
```go
package main

import (
	"log"
	"os"
	"time"

	"github.com/gdbu/poller"
)

func main() {
	var (
		f   *os.File
		p   *poller.Poller
		err error
	)

	if p, err = poller.New("./.test_file", func(e poller.Event) {
		log.Println("Event received!", e)
	}); err != nil {
		out.Errorf("error during Init: %v", err)
        return
	}

	go p.Run(0)

	time.Sleep(300 * time.Millisecond)

	if f, err = os.Create("./.test_file"); err != nil {
		out.Errorf("error during Init: %v", err)
        return
	}

	time.Sleep(300 * time.Millisecond)

	if _, err = f.WriteString("Hello world!"); err != nil {
		out.Errorf("error during Init: %v", err)
        return
	}

	time.Sleep(300 * time.Millisecond)

	if err = f.Close(); err != nil {
		out.Errorf("error during Init: %v", err)
        return
	}

	time.Sleep(300 * time.Millisecond)

	if err = os.Chmod("./.test_file", 0655); err != nil {
		out.Errorf("error during Init: %v", err)
        return
	}

	time.Sleep(300 * time.Millisecond)

	if err = os.Remove("./.test_file"); err != nil {
		out.Errorf("error during Init: %v", err)
        return
	}

	time.Sleep(300 * time.Millisecond)
	p.Close()
}

```