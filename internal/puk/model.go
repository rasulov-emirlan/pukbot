package puk

import (
	"net/url"
	"time"
)

type Puk struct {
	ID       int64 `json:"id"`
	ChatID   int64 `json:"-"`
	AuthorID int64 `json:"-"`

	AudioURL string `json:"audioURL"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewPuk(chatID, authorID int64, audioURL string) (*Puk, error) {
	_, err := url.ParseRequestURI(audioURL)
	if err != nil {
		return nil, err
	}

	return &Puk{
		ChatID:    chatID,
		AuthorID:  authorID,
		AudioURL:  audioURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}
