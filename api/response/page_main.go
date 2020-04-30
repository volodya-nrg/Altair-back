package response

type PageMain struct {
	Last PageMainLast `json:"last"`
}
type PageMainLast struct {
	AdsFull []*AdFull `json:"adsFull"`
}
