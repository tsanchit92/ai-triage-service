package incidents

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Title           string    `db:"title" json:"title"`
	Description     string    `db:"description" json:"description"`
	AffectedService string    `db:"affected_service" json:"affected_service"`
	AISeverity      string    `db:"ai_severity" json:"ai_severity"`
	AICategory      string    `db:"ai_category" json:"ai_category"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
}

type CreateIncidentRequest struct {
	Title           string `json:"title" validate:"required"`
	Description     string `json:"description" validate:"required"`
	AffectedService string `json:"affected_service" validate:"required"`
}
