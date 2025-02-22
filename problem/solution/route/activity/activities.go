package activity

import (
	"math"
	"strconv"
)

func Round(time float64) string {
	if time == math.MaxFloat64 {
		return "oo"
	}
	return strconv.FormatInt(int64(math.Round(time)), 10)
}
