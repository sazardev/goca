package repository

import (
	"database/sql"
	"github.com/test/debugproject/internal/domain"

	_ "github.com/lib/pq"
)

type postgresTestUserFeatureRepository struct {
	db *sql.DB
}

func NewPostgresTestUserFeatureRepository(db *sql.DB) TestUserFeatureRepository {
	return &postgresTestUserFeatureRepository{
		db: db,
	}
}

func (p *postgresTestUserFeatureRepository) Save(testuserfeature *domain.TestUserFeature) error {
	// TODO: Customize this query based on your TestUserFeature entity fields
	query := `INSERT INTO testuserfeatures DEFAULT VALUES RETURNING id`
	err := p.db.QueryRow(query).Scan(&testuserfeature.ID)
	return err
}

func (p *postgresTestUserFeatureRepository) FindByID(id int) (*domain.TestUserFeature, error) {
	testuserfeature := &domain.TestUserFeature{}
	// TODO: Customize this query based on your TestUserFeature entity fields
	query := `SELECT id FROM testuserfeatures WHERE id = $1`
	err := p.db.QueryRow(query, id).Scan(&testuserfeature.ID)
	if err != nil {
		return nil, err
	}

	return testuserfeature, nil
}

func (p *postgresTestUserFeatureRepository) FindByEmail(email string) (*domain.TestUserFeature, error) {
	testuserfeature := &domain.TestUserFeature{}
	// TODO: Customize this query based on your TestUserFeature entity fields
	query := `SELECT id FROM testuserfeatures WHERE id = $1 LIMIT 1`
	err := p.db.QueryRow(query, email).Scan(&testuserfeature.ID)
	if err != nil {
		return nil, err
	}
	return testuserfeature, nil
}

func (p *postgresTestUserFeatureRepository) Update(testuserfeature *domain.TestUserFeature) error {
	// TODO: Customize this query based on your TestUserFeature entity fields
	query := `UPDATE testuserfeatures SET id = $1 WHERE id = $2`
	_, err := p.db.Exec(query, testuserfeature.ID, testuserfeature.ID)
	return err
}

func (p *postgresTestUserFeatureRepository) Delete(id int) error {
	query := `DELETE FROM testuserfeatures WHERE id = $1`
	_, err := p.db.Exec(query, id)
	return err
}

func (p *postgresTestUserFeatureRepository) FindAll() ([]domain.TestUserFeature, error) {
	// TODO: Customize this query based on your TestUserFeature entity fields
	query := `SELECT id FROM testuserfeatures`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testuserfeatures []domain.TestUserFeature
	for rows.Next() {
		var testuserfeature domain.TestUserFeature
		// TODO: Scan all fields of your TestUserFeature entity
		if err := rows.Scan(&testuserfeature.ID); err != nil {
			return nil, err
		}
		testuserfeatures = append(testuserfeatures, testuserfeature)
	}

	return testuserfeatures, nil
}
