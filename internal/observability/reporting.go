package observability

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Report struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
}

func NewReport(identifier, message, level string) Report {
	reportID := fmt.Sprintf("%s-%s", identifier, uuid.New().String())
	return Report{
		ID:        reportID,
		Message:   message,
		Level:     level,
		Timestamp: time.Now(),
	}
}
