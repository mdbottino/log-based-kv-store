package main

import (
	"fmt"

	"github.com/mdbottino/log-based-kv-store/store"
)

func main() {
	s := store.NewStore()

	s.Set("banana", "pijama")

	val, _ := s.Get("banana")
	fmt.Println("'Banana' => " + val)

	_, err := s.Get("Not here")
	fmt.Println(err)
}
