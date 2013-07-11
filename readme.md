# Github Web Service hook server in Go

Crazy simple little server here.

```
package main

import (
  "github.com/derekpitt/githubservicehook"
  "log"
  "time"
)

func main() {
  processor := githubservicehook.New(":8080", func(p githubservicehook.Payload) {
    log.Println("starting: ", p.Repository.Name)
    time.Sleep(1 * time.Second)
    log.Println("done...")
  })

  log.Fatal(processor.Start())
}
```

## install

```
go install github.com/derekpitt/githubservicehook
```

## A little more detail

- Holds an internal queue so it will process requests in the order recieved while keeping githubs hook processor happy. (runs your handler in a go routine and returns as fast as possible to github)

- Each processor is contained, so you can run multiple servers in one process (each in a go routine, since Start() will block)

## TODO:

- Write some tests
