package flagger

import (
	"fmt"
	"time"
)

func extractDate(dateFormat, date string) *time.Time{
	dateTime, err := time.Parse(dateFormat, date)
	if err != nil {
		fmt.Println("ERROR ", err)
		return nil
	}
	return &dateTime
}
