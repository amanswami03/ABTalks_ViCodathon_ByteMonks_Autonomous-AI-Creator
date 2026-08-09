package models

import "time"

type RejectedTopic struct {
    TopicID string    `json:"topicId"`
    Title   string    `json:"title"`
    Reason  string    `json:"reason"`
    Time    time.Time `json:"time"`
}
