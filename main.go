package main

import (
	"fmt"

	"github.com/mdbottino/log-based-kv-store/store"
)

func main() {
	s := store.NewStore("./data")

	s.Set("banana", "pijama")
	s.Set("another", "value")
	s.Set("more", "values")
	s.Set("banana", "pijama")

	val, _ := s.Get("banana")
	fmt.Println("'Banana' => " + val)

	s.Set("dummy", "content")
	val, _ = s.Get("dummy")
	fmt.Println("'dummy' => " + val)

	_, err := s.Get("Not here")
	fmt.Println(err)
}
