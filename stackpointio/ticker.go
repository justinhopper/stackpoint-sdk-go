package stackpointio

import (
    "fmt"
    "time"
)

func ShowTicker(limit int) {
    fmt.Printf("\033[s")
    ticker := time.Tick(time.Second/4)
    spinner := []string{`\`,`|`,`/`,`-`}
    for i := 0; i <= limit; i++ {
        <-ticker
        fmt.Printf("\033[u")
        fmt.Printf("%v", spinner[i%4])
    }
    fmt.Printf("\n")
}
