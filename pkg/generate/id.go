// Package generate provides helpers for generating IDs and identifiers.
package generate

import (
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

const defaultIDLength = 30

// ID generates a new id with provided prefix. Default length is 24.
func ID(prefix string, length ...int) string {
	idLength := defaultIDLength

	if len(length) > 0 {
		idLength = length[0]
	}

	lenghtWithNoPrefix := idLength - len(prefix) - 1

	id, _ := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", lenghtWithNoPrefix)
	finalID := fmt.Sprintf("%s_%s", prefix, id)

	return finalID
}
