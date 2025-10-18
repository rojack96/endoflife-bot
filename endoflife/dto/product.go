package dto

type Product struct {
	Name                 string  `json:"name"`
	Release              string  `json:"release"`
	Released             string  `json:"released"`
	EndOfActiveSupport   *string `json:"end_of_active_support,omitempty"`
	EndOfSecuritySupport *string `json:"end_of_security_support,omitempty"`
	Latest               struct {
		Version string `json:"version"`
		Date    string `json:"date"`
		Link    string `json:"link"`
	} `json:"latest"`
}
