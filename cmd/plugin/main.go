package main

import (
	"fmt"
	"os"
)

func main() {

}

func PluginMain() {
	args := os.Args
	for _, arg := range args {
		fmt.Println(arg)
	}
}
