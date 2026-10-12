package main

import (
	"fmt"
	"time"
)

type Deadline struct {
	Course  string    `json:"course"`
	Learner string    `json:"learner"`
	Phone   string    `json:"phone"`
	DueAt   time.Time `json:"due_at"`
}

func ShouldAlert(now, due time.Time) bool { return due.After(now) && due.Sub(now) <= 24*time.Hour }

func AlertMessage(d Deadline) string {
	return fmt.Sprintf("%s deadline is %s", d.Course, d.DueAt.Format("2006-01-02 15:04"))
}

func SendDeadlineAlert(c *SMSClient, d Deadline, now time.Time) (SendResult, error) {
	if !ShouldAlert(now, d.DueAt) {
		return SendResult{}, nil
	}
	return c.Send(d.Phone, AlertMessage(d), "deadline-"+d.Learner+"-"+d.Course)
}
