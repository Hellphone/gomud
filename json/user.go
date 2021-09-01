package json

const (
	StatusOffline Status = iota + 1
	StatusOnline
)

type Status int

type User struct {
	Login    string `json:"login"`
	Level    int    `json:"level"`
	Status   Status `json:"status"`
	Location string `json:"location"`
}
