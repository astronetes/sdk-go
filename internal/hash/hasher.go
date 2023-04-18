package hash

import (
	"crypto/sha256"
	"encoding/json"
)

func Create[V any](value V) (string, error) {
	hashSrc, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	toHash := []interface{}{hashSrc}

	hash, err := generateHashFromInterfaces(toHash)
	if err != nil {
		return "", err
	}

	return hash.String(), nil
}

// GenerateHashFromInterfaces returns a hash sum based on a slice of given interfaces.
func generateHashFromInterfaces(interfaces []interface{}) (hashInterface, error) {
	hashSrc := make([]byte, len(interfaces))

	for _, in := range interfaces {
		chainElem, err := json.Marshal(in)
		if err != nil {
			return []byte{}, err
		}

		hashSrc = append(hashSrc, chainElem...)
	}

	hash := sha256.New()

	_, err := hash.Write(hashSrc)
	if err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
