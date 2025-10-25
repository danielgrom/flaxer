package repository

import (
	"database/sql"
	"errors"
	"fmt"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `create table if not exists flaxersettings(
	id integer primary key,
	projectsdirectory text null,
	flaxlocation text null)`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertFlaxerSettings(flaxerSettings FlaxerSettings) (*FlaxerSettings, error) {
	stmt := "insert into flaxersettings (projectsdirectory, flaxlocation) values (?, ?)"
	res, err := repo.Conn.Exec(stmt, flaxerSettings.ProjectsDirectory, flaxerSettings.FlaxLocation)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	flaxerSettings.ID = id
	return &flaxerSettings, nil
}

func (repo *SQLiteRepository) GetFlaxerSettings() (*FlaxerSettings, error) {
	row := repo.Conn.QueryRow("select id, projectsdirectory, flaxlocation from flaxersettings")
	var h FlaxerSettings
	err := row.Scan(&h.ID, &h.ProjectsDirectory, &h.FlaxLocation)
	if err != nil {
		if err == sql.ErrNoRows {
			// No settings found - this is normal for a fresh installation
			return nil, nil
		}
		// Actual database error
		return nil, err
	}
	return &h, nil
}

func (repo *SQLiteRepository) UpdateFlaxerSettings(id int64, update FlaxerSettings) error {
	if id == 0 {
		return errors.New("invalid updated id")
	}

	stmt := "update flaxersettings set projectsdirectory = ?, flaxlocation = ? where id = ?"
	res, err := repo.Conn.Exec(stmt, update.ProjectsDirectory, update.FlaxLocation, id)
	if err != nil {
		fmt.Printf("err=nil%s", err)
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("err affected")
		return err
	}

	if rowsAffected == 0 {
		fmt.Printf("0 rows")
		return errUpdateFailed
	}

	return nil
}
