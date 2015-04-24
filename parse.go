package irmagician

import (
	"bytes"
	"strconv"
)

func ParseRawInt(resp []byte, base int) (int, error) {
	n, err := strconv.ParseInt(string(bytes.TrimSpace(resp)), base, 10)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
