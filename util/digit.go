package util

import (
	"strconv"
)

func IsDigit(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}
