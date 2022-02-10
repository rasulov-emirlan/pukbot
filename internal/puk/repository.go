package puk

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/rasulov-emirlan/pukbot/pkg/fs"
)

const folder = "12zhCcj4sc3dW7PidEYnM34JMFY5AN3JI"

type Repository interface {
	Create(ctx context.Context, chatID, authorID int64, f io.Reader) (*Puk, error)
	List(ctx context.Context, page, limit int) ([]*Puk, error)
}

type repository struct {
	fs   fs.FileSystem
	conn *pgx.Conn
}

func NewRepository(fs fs.FileSystem, conn *pgx.Conn) Repository {
	return &repository{
		fs:   fs,
		conn: conn,
	}
}

func (r *repository) Create(ctx context.Context, chatID, authorID int64, f io.Reader) (*Puk, error) {
	filename := fmt.Sprintf("%d___%d", chatID, authorID)
	url, err := r.fs.UploadFile(filename, "audio/mp3", f, folder)
	if err != nil {
		return nil, fmt.Errorf("Repository: %v", err)
	}
	url = fmt.Sprintf("https://drive.google.com/uc?export=view&id=%s", url)
	puk, err := NewPuk(chatID, authorID, url)
	if err != nil {
		return nil, err
	}
	err = r.conn.QueryRow(ctx, `
	INSERT INTO puks(
		chat_id, author_id, audio_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`, chatID, authorID, url, puk.CreatedAt, puk.UpdatedAt).Scan(&puk.ID)
	return puk, err
}

func (r *repository) List(ctx context.Context, page, limit int) ([]*Puk, error) {
	rows, err := r.conn.Query(ctx, `
	SELECT id, chat_id, author_id, audio_url, created_at, updated_at
	FROM puks
	LIMIT $1 OFFSET $2;`, limit, page*limit)
	if err != nil {
		return nil, err
	}

	var (
		puks                 []*Puk
		id, chatID, authorID int64
		url                  string
		createdAt, updatedAt time.Time
	)

	for rows.Next() {
		if err := rows.Scan(
			&id,
			&chatID,
			&authorID,
			&url,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		puks = append(puks, &Puk{
			ID:        id,
			ChatID:    chatID,
			AuthorID:  authorID,
			AudioURL:  url,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}
	return puks, nil
}
