package models

import (
	"fmt"
	"time"
)

type Episode struct {
	Number     string    // e.g., "#037"
	Date       time.Time // Parsed date
	URL        string    // Episode URL
	IsListened bool      // Whether the episode has been listened to
}

func (e Episode) String() string {
	listenedMark := "[ ]"
	if e.IsListened {
		listenedMark = "[x]"
	}
	
	// Format: [ ] 06/11 [#38](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55)
	if e.URL != "" {
		return fmt.Sprintf("%s %02d/%02d [%s](%s)", listenedMark, e.Date.Month(), e.Date.Day(), e.Number, e.URL)
	}
	return fmt.Sprintf("%s %02d/%02d [%s]", listenedMark, e.Date.Month(), e.Date.Day(), e.Number)
}