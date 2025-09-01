package incidents

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Insert(ctx context.Context, i Incident) error
	List(ctx context.Context) ([]Incident, error)
}

type MySQLRepository struct{ db *sqlx.DB }

// NewRepository returns a MySQL-compatible repository
func NewRepository(db *sqlx.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// Insert adds a new incident to the database
func (r *MySQLRepository) Insert(ctx context.Context, i Incident) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO incidents 
			(id, title, description, affected_service, ai_severity, ai_category, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`, i.ID.String(), i.Title, i.Description, i.AffectedService, i.AISeverity, i.AICategory)
	return err
}

// List fetches all incidents from the database
func (r *MySQLRepository) List(ctx context.Context) ([]Incident, error) {
	var out []Incident
	err := r.db.SelectContext(ctx, &out, `
		SELECT id, title, description, affected_service, ai_severity, ai_category, created_at
		FROM incidents
		ORDER BY created_at DESC
	`)
	return out, err
}

// NewIncident creates a new Incident struct with a generated UUID
func NewIncident(title, desc, svc, sev, cat string) Incident {
	return Incident{
		ID:              uuid.New(),
		Title:           title,
		Description:     desc,
		AffectedService: svc,
		AISeverity:      sev,
		AICategory:      cat,
	}
}
