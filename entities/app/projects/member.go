package projects

type Member struct {
	ID uint `json:"id"`

	Account uint `json:"account"`
	Project uint `json:"project"`
	Role    uint `json:"role"`
}
