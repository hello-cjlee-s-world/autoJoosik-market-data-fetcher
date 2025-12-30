package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseInt(s string) int64 {
	s = strings.ReplaceAll(s, "+", "")
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func ParseFloat(s string) float64 {
	s = strings.ReplaceAll(s, "+", "")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

func ParseTradeTimestamp(tm string) time.Time {
	// date: 오늘 날짜
	today := time.Now().Format("2006-01-02")
	layout := "2006-01-02 150405" // yyyy-MM-dd HHMMSS
	combined := fmt.Sprintf("%s %s", today, tm)
	t, err := time.ParseInLocation(layout, combined, time.Local)
	if err != nil {
		return time.Now()
	}
	return t
}

func ParseSignedFloat(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSuffix(s, "%")
	if s == "" || s == "--" {
		return 0, nil
	}
	return strconv.ParseFloat(s, 64)
}
