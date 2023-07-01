package helpers

import (
	"fmt"
	"math/rand"
	"strings"
)

func CreateConfirmationCode() string {
	float := rand.Float64()
	str := fmt.Sprintf("%f", float)
	return strings.Replace(str, "0.", "", 1)
}
