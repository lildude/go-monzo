package main

import (
	"fmt"

	"github.com/gurparit/go-monzo/monzo"
)

func main() {
	m := monzo.Monzo{}

	fmt.Printf("%+v\n", m)
}
