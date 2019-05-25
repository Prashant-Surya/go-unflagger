package flagger

import (
	"fmt"
	"time"
)

func extractDate(date string) *time.Time{
	format := "2006_01_02"
	dateTime, err := time.Parse(format, date)
	if err != nil {
		fmt.Println("ERROR ", err)
		return nil
	}
	return &dateTime
}
