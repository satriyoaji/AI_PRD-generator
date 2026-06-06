package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func Connect(databaseURL string) (*pgx.Conn, error) {
	if databaseURL == "" {
		log.Println("No DATABASE_URL provided, using mock database")
		return nil, nil
	}

	var err error
	conn, err = pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	log.Println("Connected to database")
	return conn, nil
}

func GetConn() *pgx.Conn {
	return conn
}

func InitSchema(conn *pgx.Conn) error {
	if conn == nil {
		return nil
	}

	schema := `
	CREATE TABLE IF NOT EXISTS transcripts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS generations (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		transcript_id UUID REFERENCES transcripts(id),
		type VARCHAR(50) NOT NULL,
		content JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_generations_transcript ON generations(transcript_id);
	CREATE INDEX IF NOT EXISTS idx_generations_type ON generations(type);
	`

	_, err := conn.Exec(context.Background(), schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Println("Database schema initialized")
	return nil
}

func Close() {
	if conn != nil {
		conn.Close(context.Background())
	}
}