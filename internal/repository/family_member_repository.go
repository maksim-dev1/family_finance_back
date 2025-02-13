package repository

import (
	"database/sql"
	"time"

	"myapp/internal/models"
)

type FamilyMemberRepository interface {
	AddMember(familyMember *models.FamilyMember) error
	RemoveMember(id string) error
	GetMembersByFamily(familyID string) ([]*models.FamilyMember, error)
}

type familyMemberRepository struct {
	db *sql.DB
}

func NewFamilyMemberRepository(db *sql.DB) FamilyMemberRepository {
	return &familyMemberRepository{db: db}
}

func (r *familyMemberRepository) AddMember(familyMember *models.FamilyMember) error {
	query := `INSERT INTO family_members (id, family_id, user_id, role, joined_at)
              VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.db.Exec(query, familyMember.ID, familyMember.FamilyID, familyMember.UserID, familyMember.Role)
	return err
}

func (r *familyMemberRepository) RemoveMember(id string) error {
	query := `DELETE FROM family_members WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *familyMemberRepository) GetMembersByFamily(familyID string) ([]*models.FamilyMember, error) {
	query := `SELECT id, family_id, user_id, role, joined_at FROM family_members WHERE family_id = $1`
	rows, err := r.db.Query(query, familyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []*models.FamilyMember
	for rows.Next() {
		var member models.FamilyMember
		if err := rows.Scan(&member.ID, &member.FamilyID, &member.UserID, &member.Role, &member.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, &member)
	}
	return members, nil
}
