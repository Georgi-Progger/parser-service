package annoucement

import (
	"context"
	"database/sql"
	"log"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) GetAnnoucement(ctx context.Context) (*[]Annoucement, error) {
	query := `
			SELECT * FROM Annoucement
			LIMIT 10;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Fatal("Not correct query", query)
	}
	defer rows.Close()
	annoucements := []Annoucement{}

	for rows.Next() {
		annoucement := Annoucement{}
		err := rows.Scan(
			&annoucement.Id,
			&annoucement.Model,
			&annoucement.Price,
			&annoucement.Year,
			&annoucement.Generation,
			&annoucement.Mileage,
			&annoucement.History,
			&annoucement.PTS,
			&annoucement.Owners,
			&annoucement.Condition,
			&annoucement.Modification,
			&annoucement.EngineVolume,
			&annoucement.EngineType,
			&annoucement.Transmission,
			&annoucement.Drive,
			&annoucement.Equipment,
			&annoucement.BodyType,
			&annoucement.Color,
			&annoucement.Steering,
			&annoucement.VIN,
			&annoucement.Exchange,
		)
		if err != nil {
			log.Fatal("Error scanning row:", err)
			return nil, err
		}

		annoucements = append(annoucements, annoucement)
	}

	return &annoucements, nil
}

func (r *repo) SetAnnoucement(ctx context.Context, annoucementInfo Annoucement) error {
	query := `
			INSERT INTO annoucement (
				Model, Price, Year, Generation, Mileage, History, PTS, Owners,
				Condition, Modification, Engine_Volume, Engine_Type, Transmission,
				Drive, Equipment, Body_Type, Color, Steering, VIN, Exchange
			) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`
	_, err := r.db.ExecContext(ctx, query,
		annoucementInfo.Model, annoucementInfo.Price, annoucementInfo.Year, annoucementInfo.Generation,
		annoucementInfo.Mileage, annoucementInfo.History, annoucementInfo.PTS, annoucementInfo.Owners,
		annoucementInfo.Condition, annoucementInfo.Modification, annoucementInfo.EngineVolume,
		annoucementInfo.EngineType, annoucementInfo.Transmission, annoucementInfo.Drive,
		annoucementInfo.Equipment, annoucementInfo.BodyType, annoucementInfo.Color,
		annoucementInfo.Steering, annoucementInfo.VIN, annoucementInfo.Exchange)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
