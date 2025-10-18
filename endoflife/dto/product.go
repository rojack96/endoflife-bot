package dto

import "time"

type Product struct {
	Name                 string    `json:"name"`
	Release              string    `json:"release"`
	Released             time.Time `json:"released"`
	EndOfActiveSupport   time.Time `json:"end_of_active_support"`
	EndOfSecuritySupport time.Time `json:"end_of_security_support"`
	Latest               struct {
		Version string    `json:"version"`
		Date    time.Time `json:"date"`
		Link    string    `json:"link"`
	} `json:"latest"`
}
