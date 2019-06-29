package main

import (
	"fmt"
)

func main() {
	fmt.Println("Starting example suite")

	s := simpleFive()

	// full
	fmt.Println("Deleting indices 0 to 3")
	s.DeleteBounds(0, 3)
	fmt.Println(s.values)

	fmt.Println("new start", s.start)
	s.Append(25)
	fmt.Println(s.values)

	// same test but with length instead of bounds

	s = simpleFive()

	s.DeleteCount(4)

	fmt.Println(s.values)

	fmt.Println("new start", s.start)
	s.Append(25)
	fmt.Println(s.values)

}

func simpleFive() *Slice {

	s := NewSlice(5, false)
	count := 0
	for i := 1; i <= 5; i++ {
		err := s.Append(i)
		if err != nil {
			count++
		}
	}
	return s
}
