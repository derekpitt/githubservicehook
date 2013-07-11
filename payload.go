package githubservicehook

import (
	"encoding/json"
	"time"
)

type Payload struct {
	Before string `json:"before"`
	After  string `json:"after"`
	Ref    string `json:"ref"`

	Commits    []commit   `json:"commits"`
	Repository repository `json:"repository"`
}

type commit struct {
	Id        string    `json:"id"`
	Author    person    `json:"author"`
	Committer person    `json:"committer"`
	Distinct  bool      `json:"distinct"`
	Message   string    `json:"message"`
	Modified  []string  `json:"modified"`
	Removed   []string  `json:"removed"`
	Timestamp time.Time `json:"timestamp"`
	Url       string    `json:"url"`
}

type repository struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Private     bool   `json:"private"`
	Owner       person `json:"owner"`
}

type person struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"username"`
}

func parsePayload(body string) (payload Payload, err error) {
	err = json.Unmarshal([]byte(body), &payload)
	return
}
