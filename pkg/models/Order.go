package models

type Orders struct {
	ID       int64  `json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}
