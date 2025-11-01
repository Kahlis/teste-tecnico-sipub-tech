package domain

import (
	"time"
)

type Movie struct {
	Id        int32     `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title"`
	Year      int       `json:"year" bson:"year"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
