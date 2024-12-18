package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const DBPath = "./cards.db"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	Migrate()
}

func Migrate() {
	query := `
    CREATE TABLE IF NOT EXISTS cards (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        word TEXT UNIQUE,
        example TEXT,
        description TEXT,
        power INTEGER,
        revision_date TIMESTAMP
    );
    `
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
}

// NOTE: check if word exists or not, with CheckIfWordExists, in the func that calls AddCard
// to avoid error from UNIQUE constraint.
// Also trim spaces around wrod and check if it's empty or not.
func AddCard(card *Card) error {
	query := `
    INSERT INTO cards (word, example, description, power, revision_date)
    VALUES (?, ?, ?, ?, ?);
    `
	_, err := db.Exec(query,
		card.word,
		card.example,
		card.description,
		card.power,
		card.revisionDate,
	)
	return err
}

func CheckIfWordExists(word string) (bool, error) {
	query := `SELECT 1 FROM cards WHERE word LIKE ?;`
	err := db.QueryRow(query, word).Scan(new(int))
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

func DeleteCard(id int) error {
	query := `DELETE FROM cards WHERE id = ?;`
	_, err := db.Exec(query, id)
	return err
}

func UpdateCard(card *Card) error {
	query := `
    UPDATE cards
    SET
        word = ?,
        example = ?,
        description = ?,
        power = ?,
        revision_date = ?
    WHERE id = ?;
    `
	_, err := db.Exec(query,
		card.word,
		card.example,
		card.description,
		card.power,
		card.revisionDate,
		card.id,
	)
	return err
}

func GetNextCard() (*Card, error) {
	query := `
    SELECT *
    FROM cards
    WHERE revision_date <= DATE('now')
    ORDER BY revision_date
    LIMIT 1;
    `
	card := &Card{}
	if err := db.QueryRow(query).Scan(
		&card.id,
		&card.word,
		&card.example,
		&card.description,
		&card.power,
		&card.revisionDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return card, nil
}
