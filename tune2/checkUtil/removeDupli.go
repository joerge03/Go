package main

import "fmt"

func main() {
	test := make(map[string]bool)

	testString := []string{"!23", "123", "51233", "123", "12532"}

	populate := []string{}

	for _, testStr := range testString {
		testasdf, testing := test[testStr]
		fmt.Println(testasdf, "asddf", testing)
		if testing := test[testStr]; testing {
			fmt.Println(testing)
			continue
		}

		test[testStr] = true
	}

	for testKey := range test {
		populate = append(populate, testKey)
	}

}
