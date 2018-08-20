package timezone

type Timezone struct {
	Abbreviation string   `json:"abbreviation"`
	Name         string   `json:"name"`
	Locations    []string `json:"locations"`
	Offset       int	  `json:"offset"`
}

