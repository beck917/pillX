// new_server
package main

import (
	"fmt"
	"../../pillX"
)

func main() {
	test1 := pillx.ReturnStr()
	fmt.Printf("ReturnStr from package1: %s\n", test1)
	fmt.Println("Hello World!")
}
