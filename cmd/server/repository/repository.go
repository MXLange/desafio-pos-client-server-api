package repository

import (
	"context"
	"database/sql"

	apperrors "github.com/MXLange/desafio-pos-client-server-api/cmd/server/app_errors"

	"github.com/MXLange/desafio-pos-client-server-api/pkg/types"
)


type PriceRepository struct {
	db *sql.DB
}

func New(db *sql.DB) (*PriceRepository, error) {
	if db == nil {
		return nil, apperrors.ErrNilDB
	}

	return &PriceRepository{db: db}, nil
}

// CreateTables creates the necessary tables for storing price data.
func (d *PriceRepository) CreateTables() error {
	_, err := d.db.Exec(`
			CREATE TABLE IF NOT EXISTS prices (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				code VARCHAR(10) NOT NULL,
				codein VARCHAR(10) NOT NULL,
				name VARCHAR(50) NOT NULL,
				high VARCHAR(20) NOT NULL,
				low VARCHAR(20) NOT NULL,
				varBid VARCHAR(20) NOT NULL,
				pctChange VARCHAR(20) NOT NULL,
				bid VARCHAR(20) NOT NULL,
				ask VARCHAR(20) NOT NULL,
				timestamp VARCHAR(20) NOT NULL,
				create_date VARCHAR(20) NOT NULL
			);
		`)
	
	if err != nil {
		return err
	}

	return d.clearTableAndResetID()
}

func (d *PriceRepository) clearTableAndResetID() error {
	_, err := d.db.Exec(`
		DELETE FROM prices;
		DELETE FROM sqlite_sequence WHERE name='prices';
	`)
	return err
}


// InsertPrice inserts a new price record into the database and returns the inserted record's ID.
func (d *PriceRepository) InsertPrice(ctx context.Context,price *types.Price) (*int64, error) {
	st, err := d.db.Prepare(`
		INSERT INTO prices (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return nil, err
	}

	res, err := st.ExecContext(ctx,
		price.USDBRL.Code,
		price.USDBRL.Codein,
		price.USDBRL.Name,
		price.USDBRL.High,
		price.USDBRL.Low,
		price.USDBRL.VarBid,
		price.USDBRL.PctChange,
		price.USDBRL.Bid,
		price.USDBRL.Ask,
		price.USDBRL.Timestamp,
		price.USDBRL.CreateDate,
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &id, nil

}