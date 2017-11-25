try
===

Package try provides simple, throttled retrying.

Godoc: https://godoc.org/github.com/bruston/x/try

## Examples

Throttling with a fixed pause duration:

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "time"

    "github.com/bruston/x/try"
)

func main() {
    buf := &bytes.Buffer{}
    err := try.Do(5, func() error {
        resp, err := http.Get("https://google.com")
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        if _, err := io.Copy(buf, resp.Body); err != nil {
            return err
        }
        return nil
    }, try.Delay(time.Second))
    if err != nil {
        fmt.Println("failed and done retrying")
        return
    }
    fmt.Println(buf)
}
```

Backing off:

```go
package main

import (
    "fmt"
    "net"
    "time"

    "github.com/bruston/x/try"
)

func main() {
    backoff := &try.Backoff{
        Duration:   time.Millisecond * 200,
        Multiplier: 1,
    }

    var conn net.Conn
    err := try.Do(10, func() error {
        c, err := net.DialTimeout("tcp", "chat.freenode.net:6667", time.Second*5)
        if err != nil {
            return err
        }
        conn = c
        return nil
    }, backoff)
    if err != nil {
        fmt.Println("failed and done retrying")
        return
    }

    fmt.Println(conn.RemoteAddr())
    conn.Close()
}
```

# Throttler Interface

To use a custom Throttler just implement the Throttle method:

```go
type Throttler interface {
    Throttle()
}
```

Which can be passed to [Do](https://godoc.org/github.com/bruston/x/try#Do).
