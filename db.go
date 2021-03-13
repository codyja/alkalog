package main

import (
	"context"
	"github.com/codyja/alkatronic/api"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type PostgresAlkatronic struct {
	db *pgxpool.Pool
}

func NewPostgresAlkatronic(addr string) (*PostgresAlkatronic, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		log.Fatalf("unable to parse connection string %s", err)
	}

	conn, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("unable to connect to database %s", err)
	}

	return &PostgresAlkatronic{db: conn}, nil
}

func (a *PostgresAlkatronic) InsertRecord(r api.Record, d api.Device) error {
	s := `INSERT INTO alkatronic(acid_used, create_time, device_id, device_name, is_deleted, is_hidden, kh, note, record_id, remaining_reagent, solution_added, lower_ref, upper_ref)
          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	      ON CONFLICT ON CONSTRAINT alkatronic_record_id_key
	      DO NOTHING;`
	_, err := a.db.Exec(context.Background(), s,
		r.AcidUsed, r.CreateTime, r.DeviceID, d.FriendlyName, r.IsDeleted, r.IsHidden, r.KhValue, r.Note, r.RecordID, r.RemainingReagent, r.SolutionAdded, d.LowerKh, d.UpperKh)
	if err != nil {
		log.Fatalf("error writing record to database: %s", err)
	}

	return nil
}
