package main

import (
	"fmt"
	"greengo/datachunk"
	"os"
)

func main() {

	// This will host all concat chunks
	var chunks []int

	// begin and end sequence to identify badge
	startCharacters := []string{"SOH", "STX"} // ABC
	endCharacters := []string{"ETX", "EOT"}   // XYZ

	// chunks of data
	slice := [][]int{{10, 20, 30, 40, 65}, {1, 66}, {1, 2, 200, 30, 40, 50}, {88, 3}, {4, 66, 36, 3, 95, 5}}

	fmt.Println("startCharacters:", startCharacters)
	fmt.Println("endCharacters:", endCharacters)

	// here we will call GetDataChunk function from our own private module greengo/datachunk
	for i := range slice {
		badge, isBadgeFound := datachunk.GetDataChunk(slice[i], &chunks, startCharacters, endCharacters)
		if isBadgeFound {
			fmt.Println("chunks:", chunks)
			fmt.Println("badge found:", badge)
			chunks = nil
			os.Exit(1)
		}
	}
	fmt.Println("badge not found")
	fmt.Println("chunks:", chunks)
	chunks = nil
}
