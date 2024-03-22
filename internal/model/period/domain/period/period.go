package period

type Period struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Year      int    `json:"year"`
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}
