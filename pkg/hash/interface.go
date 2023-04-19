package hash

import "fmt"

type hashInterface []byte

// String returns the string value of an interface hash.
func (hash *hashInterface) String() string {
	return fmt.Sprintf("%x", *hash)
}

// Short returns the first 8 characters of an interface hash.
func (hash *hashInterface) Short() string {
	return hash.String()[:8]
}
