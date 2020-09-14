package main

import (
	"log"
	"os"
	"time"

	"github.com/hatchify/poller"
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
		log.Fatal(err)
	}

	go p.Run(0)

	time.Sleep(300 * time.Millisecond)

	if f, err = os.Create("./.test_file"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if _, err = f.WriteString("Hello world!"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if err = os.Chmod("./.test_file", 0655); err != nil {
		log.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)

	if err = os.Remove("./.test_file"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(300 * time.Millisecond)
	p.Close()
}
