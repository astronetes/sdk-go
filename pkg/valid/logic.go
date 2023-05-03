package valid

import "fmt"

func OnlyOneOfThemIsInitialized(items ...interface{}) (int, error) {
	notNilIndex := -1
	for i, item := range items {
		if item != nil {
			if notNilIndex > -1 {
				return -1, fmt.Errorf("more than one element has been provided")
			}
			notNilIndex = i
		}
	}
	if notNilIndex < 0 {
		return -1, fmt.Errorf("none of the elements has been provided")
	}
	return notNilIndex, nil
}
