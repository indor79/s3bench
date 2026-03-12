package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}
	mul := int64(1)
	switch {
	case strings.HasSuffix(s, "KIB"):
		mul = 1024
		s = strings.TrimSuffix(s, "KIB")
	case strings.HasSuffix(s, "MIB"):
		mul = 1024 * 1024
		s = strings.TrimSuffix(s, "MIB")
	case strings.HasSuffix(s, "GIB"):
		mul = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "GIB")
	case strings.HasSuffix(s, "KB"):
		mul = 1000
		s = strings.TrimSuffix(s, "KB")
	case strings.HasSuffix(s, "MB"):
		mul = 1000 * 1000
		s = strings.TrimSuffix(s, "MB")
	case strings.HasSuffix(s, "GB"):
		mul = 1000 * 1000 * 1000
		s = strings.TrimSuffix(s, "GB")
	case strings.HasSuffix(s, "B"):
		mul = 1
		s = strings.TrimSuffix(s, "B")
	}
	v, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil || v <= 0 {
		return 0, fmt.Errorf("invalid size: %q", s)
	}
	return v * mul, nil
}
