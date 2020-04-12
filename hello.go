package main

import (
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/testingtony/hello/morestrings"
)

func main() {
	fmt.Println(morestrings.ReverseRunes("Hello, world."))

	fmt.Println(cmp.Diff("Hello World", "Hello world"))

	time.Sleep(20 * time.Second)
}
