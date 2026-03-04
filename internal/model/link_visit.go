package model

import "time"

type LinkVisit struct {
	ID        int64     `json:"id"`
	LinkId    int64     `json:"link_id"`
	Ip        string    `json:"ip"`
	UserAgent *string   `json:"user_agent,omitempty"`
	Referer   *string   `json:"referer,omitempty"`
	Status    int64     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
