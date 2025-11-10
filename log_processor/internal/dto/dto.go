package dto

import "time"

type FilterDTO struct {
	Levels []string  `json:"levels"`  // e.g. ["INFO", "ERROR"]
	Search string    `json:"search"`  // keyword search
	From   time.Time `json:"from"`    // start time
	To     time.Time `json:"to"`      // end time
	Limit  int       `json:"limit"`   // pagination limit
	Offset int       `json:"offset"`  // pagination offset
}