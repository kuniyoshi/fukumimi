package models

import (
	"fmt"
	"time"
)

type Episode struct {
	Number     string    // e.g., "#037"
	Date       time.Time // Parsed date
	Year       int       // Year from the listing page like "2025.06.02"
	URL        string    // Episode URL
	IsListened bool      // Whether the episode has been listened to
}

func (e Episode) String() string {
	listenedMark := "[ ]"
	if e.IsListened {
		listenedMark = "[x]"
	}
	
	// Format: [ ] 2025-06-11 [#38](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55)
	if e.URL != "" {
		return fmt.Sprintf("%s %s [%s](%s)", listenedMark, e.Date.Format("2006-01-02"), e.Number, e.URL)
	}
	return fmt.Sprintf("%s %s [%s]", listenedMark, e.Date.Format("2006-01-02"), e.Number)
}