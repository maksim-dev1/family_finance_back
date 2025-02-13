package repository

import (
	"database/sql"
	"errors"
	"time"

	"myapp/internal/models"
)

type FamilyRepository interface {
	CreateFamily(family *models.Family) error
	GetFamilyByID(id string) (*models.Family, error)
	GetFamiliesByUser(userID string) ([]*models.Family, error)
}

type familyRepository struct {
	db *sql.DB
}

func NewFamilyRepository(db *sql.DB) FamilyRepository {
	return &familyRepository{db: db}
}

func (r *familyRepository) CreateFamily(family *models.Family) error {
	query := `INSERT INTO families (id, name, description, created_at, updated_at)
              VALUES ($1, $2, $3, NOW(), NOW())`
	_, err := r.db.Exec(query, family.ID, family.Name, family.Description)
	return err
}

func (r *familyRepository) GetFamilyByID(id string) (*models.Family, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM families WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var family models.Family
	err := row.Scan(&family.ID, &family.Name, &family.Description, &family.CreatedAt, &family.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("семья не найдена")
		}
		return nil, err
	}
	return &family, nil
}

func (r *familyRepository) GetFamiliesByUser(userID string) ([]*models.Family, error) {
	query := `SELECT f.id, f.name, f.description, f.created_at, f.updated_at
              FROM families f
              JOIN family_members fm ON f.id = fm.family_id
              WHERE fm.user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var families []*models.Family
	for rows.Next() {
		var family models.Family
		if err := rows.Scan(&family.ID, &family.Name, &family.Description, &family.CreatedAt, &family.UpdatedAt); err != nil {
			return nil, err
		}
		families = append(families, &family)
	}
	return families, nil
}
