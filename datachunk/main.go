package datachunk

import (
	"strconv"
	"strings"
)

// GetDataChunk : Main function will provide a badge if exist with chunks
func GetDataChunk(chunk []int, chunks *[]int, startWith []string, endWith []string) (myBadge []int, isMyBadgeFound bool) {

	// variable for badge in string and bounds conversion from ascii to string
	var badgeInString, downBoundary, upBoundary string
	// let's add current chunk to all chunks
	for k := range chunk {
		*chunks = append(*chunks, chunk[k])
	}

	// transform list of int into string to facilitate badge detection
	ChunksInString := sliceToString(*chunks)

	// Convert boundaries ascii into two strings
	for i := range startWith {
		number := startWith[i]
		downBoundary += AsciiToBase10(number) + " "
	}
	downBoundary = strings.TrimRight(downBoundary, " ")

	for i := range endWith {
		number := endWith[i]
		upBoundary += AsciiToBase10(number) + " "
	}
	upBoundary = strings.TrimRight(upBoundary, " ")

	// If badge is not found, badge variable is empty
	badgeInString, isMyBadgeFound = getSubString(ChunksInString, downBoundary, upBoundary)

	// revert badge string to slice of int
	myBadge = stringToSlice(badgeInString)

	return myBadge, isMyBadgeFound
}

// getSubString : Find a substring into a string based on inner and upper bounds
func getSubString(myString string, startBound string, endBound string) (myBadge string, found bool) {

	// Searching for occurrence of first upBoundary
	sliceWithStartBound := strings.Index(myString, startBound)
	// startBound was not found, exit
	if sliceWithStartBound == -1 {
		return "", false
	}

	// New string from startBound to end of original string, myString
	newString := myString[sliceWithStartBound+len(startBound):]

	sliceWithStartandEndBounds := strings.Index(newString, endBound)
	if sliceWithStartandEndBounds == -1 {
		return "", false
	}

	// to dig why we should begin slice from 1 to avoid empty slice with len = 1 for example (thanks to unit test)
	myBadge = newString[1:sliceWithStartandEndBounds]

	// if two bounds were founds but no element in between, this should be trapped (thanks to unit test)
	if len(myBadge) == 0 {
		return "", false
	}

	return myBadge, true
}

// sliceToString : Convert a slice into a string
func sliceToString(mySlice []int) (myList string) {

	// this will host final string
	var valuesText []string

	for i := range mySlice {
		number := mySlice[i]
		text := strconv.Itoa(int(number))
		valuesText = append(valuesText, text)
	}

	result := strings.Join(valuesText, " ")

	// We need to separate each new chunk with a space as soon as we concat multiple chunks
	if myList == "" {
		myList = result
	} else {
		myList += " " + result
	}

	return myList
}

// stringToSlice : Convert a string into a slice
func stringToSlice(myString string) (mySlice []int) {

	// this will host slice result
	numbers := strings.Fields(myString)

	for i := range numbers {
		number, _ := strconv.Atoi(numbers[i])
		mySlice = append(mySlice, number)
	}

	return mySlice
}

// AsciiToBase10 : Return base 10 value of an ascii character (string input)
func AsciiToBase10(myString string) (result string) {

	var intValue int
	switch myString {
	case "NUL":
		intValue = 0
	case "SOH":
		intValue = 1
	case "STX":
		intValue = 2
	case "ETX":
		intValue = 3
	case "EOT":
		intValue = 4
	case "ENQ":
		intValue = 5
	case "ACK":
		intValue = 6
	case "BEL":
		intValue = 7
	case "BS":
		intValue = 8
	case "HT":
		intValue = 9
	case "LF":
		intValue = 10
	case "VT":
		intValue = 11
	case "FF":
		intValue = 12
	case "CR":
		intValue = 13
	case "SO":
		intValue = 14
	case "SI":
		intValue = 15
	case "DLE":
		intValue = 16
	case "DC1":
		intValue = 17
	case "DC2":
		intValue = 18
	case "DC3":
		intValue = 19
	case "DC4":
		intValue = 20
	case "NAK":
		intValue = 21
	case "SYN":
		intValue = 22
	case "ETB":
		intValue = 23
	case "CAN":
		intValue = 24
	case "EM":
		intValue = 25
	case "SUB":
		intValue = 26
	case "ESC":
		intValue = 27
	case "FS":
		intValue = 28
	case "GS":
		intValue = 29
	case "RS":
		intValue = 30
	case "US":
		intValue = 31
	case "SPACE":
		intValue = 32
	default:
		intValue = int(myString[0])
	}

	result = strconv.Itoa(intValue)
	return result
}
