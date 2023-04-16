package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: rm [file or directory]")
        os.Exit(1)
    }

    path := os.Args[1]

    err := os.RemoveAll(path)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
