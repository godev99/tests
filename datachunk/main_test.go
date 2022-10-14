package datachunk

import (
	"strconv"
	"testing"
)

func TestGetDataChunk(t *testing.T) {
	tests := []struct {
		name      string
		startWith []string
		endWith   []string
		chunk     []int
	}{
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 20, 30, 40, 3}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 1, 65, 66, 65}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{1, 66, 65, 40, 2}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 6, 65, 2, 6}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 10, 1, 3, 2, 4, 65}},
		{name: "KO", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 3, 30, 40, 1}},
		{name: "KO", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 20, 65, 3, 65}},
		{name: "OK", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{1, 66, 65, 3, 4}},
		{name: "KO", startWith: []string{"SOH"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 4, 65, 1, 6}},
		{name: "OK", startWith: []string{"STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 10, 1, 20, 30, 3, 4}},
		{name: "KO", startWith: []string{"STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 1, 30, 40, 66}},
		{name: "KO", startWith: []string{"STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{1, 20, 65, 66, 65}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 66, 65, 40, 65}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{10, 6, 65, 40, 6}},
		{name: "KO", startWith: []string{"SOH", "STX"}, endWith: []string{"ETX", "EOT"}, chunk: []int{65, 10, 57, 20, 30, 40, 65}},
	}
	t.Log("Testing GetDataChunk with various scenario")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {
				var chunks []int

				t.Logf("\tTest %d", i)
				t.Log("\t\tchunk: ", tt.chunk)
				t.Log("\t\tstartWith: ", tt.startWith)
				t.Log("\t\tendCharacters: ", tt.endWith)
				badge, found := GetDataChunk(tt.chunk, &chunks, tt.startWith, tt.endWith)
				if !found {
					t.Error("\t\tbadge not found")
				} else {
					t.Log("\t\tbadge found:", badge)
				}

			}
			t.Run(tt.name, tf)
		}
	}
}

func TestSliceToString(t *testing.T) {
	tests := []struct {
		name    string
		mySlice []int
	}{
		{name: "65, 54, 678, 23, 0, 1", mySlice: []int{65, 54, 678, 23, 0, 1}},
		{name: "-1, 9, 54, 0, 54", mySlice: []int{-1, 9, 54, 0, 54}},
		{name: "0", mySlice: []int{0}},
		{name: "235334, 564, 6776", mySlice: []int{235334, 564, 6776}},
	}
	t.Log("Testing sliceToString with various scenario")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest %d", i)
				t.Log("\t\tslice: ", tt.mySlice)

				var myString string
				myString = sliceToString(tt.mySlice)
				t.Log("\t\tstring: ", myString)

				if myString == "" {
					t.Error("\t\tstring is empty")
				}

			}
			t.Run(tt.name, tf)
		}
	}
}

func TestStringToSlice(t *testing.T) {
	tests := []struct {
		myString string
	}{
		{myString: "65, 54, 678, 23, 0, 1"},
		{myString: "-1, 9, 54, 0, 54"},
		{myString: "0"},
		{myString: "235334, 564, 6776"},
	}
	t.Log("Testing stringToSlice with various scenario")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest %d", i)
				t.Log("\t\tstring: ", tt.myString)

				var mySlice []int
				mySlice = stringToSlice(tt.myString)
				t.Log("\t\tslice: ", mySlice)

				if len(mySlice) == 0 {
					t.Error("\t\tslice is empty")
				}

			}
			t.Run(tt.myString, tf)
		}
	}
}

func TestGetSubString(t *testing.T) {
	tests := []struct {
		name      string
		myString  string
		startWith string
		endWith   string
	}{
		{name: "65 54 678 23 0 1 (65 1)", myString: "65 54 678 23 0 1", startWith: "65", endWith: "1"},
		{name: "1 9 54 0 54 (9 -1)", myString: "-1 9 54 0 54", startWith: "9", endWith: "-1"},
		{name: "0 (65 1)", myString: "0", startWith: "65", endWith: "1"},
		{name: "12 434 65 564 6776 (65 564)", myString: "12 434 65 564 6776", startWith: "65", endWith: "564"},
		{name: "43 65 56 67 45 12 (56 12)", myString: "43 65 56 67 45 12", startWith: "56", endWith: "12"},
		{name: "43 65 56 67 45 12 (56 45)", myString: "43 65 56 67 45 12", startWith: "56", endWith: "45"},
	}
	t.Log("Testing sliceToString with various scenario")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest %d", i)
				t.Log("\t\tlist: ", tt.myString)
				t.Log("\t\tfirst delimiter: ", tt.startWith)
				t.Log("\t\tlast delimiter : ", tt.endWith)

				substring, success := getSubString(tt.myString, tt.startWith, tt.endWith)

				if success {
					t.Log("\t\tsubstring was found:", substring)
				} else {
					t.Error("\t\tsubstring was not found")
				}

			}
			t.Run(tt.name, tf)
		}
	}
}

func TestAsciiToBase10(t *testing.T) {
	tests := []struct {
		myString string
	}{
		{myString: "SOH"},
		{myString: "STX"},
		{myString: "EOT"},
		{myString: "EQ"},
		{myString: "?"},
		{myString: "ETX"},
		{myString: "NUL"},
		{myString: "ENQ"},
		{myString: "LF"},
		{myString: "LIFO"},
		{myString: "'"},
		{myString: "DLE"},
		{myString: "DC1"},
		{myString: "SYN"},
		{myString: "LIO"},
		{myString: "}"},
		{myString: "EM"},
		{myString: "NAK"},
		{myString: "CAN"},
		{myString: "RLT"},
		{myString: "ESC"},
	}
	t.Log("Testing TestAsciiToBase10 with various scenario")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest %d", i)
				t.Log("\t\tstring: ", tt.myString)

				result := AsciiToBase10(tt.myString)
				resultInt, _ := strconv.Atoi(result)

				if result != "" {
					if resultInt < 33 && len(tt.myString) > 1 {
						t.Log("\t\tascii code was found:", result)
					} else if resultInt > 33 && len(tt.myString) > 1 {
						t.Log("\t\tlen myString:", len(tt.myString))
						t.Log("\t\tresult:", result)
						t.Log("\t\tresultInt:", resultInt)
						t.Error("\t\tascii symbol is invalid.")
					}

				} else {
					t.Error("\t\tascii code was not found")
				}

			}
			t.Run(tt.myString, tf)
		}
	}
}
