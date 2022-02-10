package puk

import (
	"context"
	"io"

	"github.com/rasulov-emirlan/pukbot/pkg/logger"
)

type Service interface {
	Create(ctx context.Context, chatID, authorID int64, f io.Reader) (*Puk, error)
	List(ctx context.Context, page, limit int) ([]*Puk, error)
}

type service struct {
	l logger.Logger
	r Repository
}

func NewService(l logger.Logger, r Repository) Service {
	return &service{
		l: l,
		r: r,
	}
}

func (s *service) Create(ctx context.Context, chatID, authorID int64, f io.Reader) (*Puk, error) {
	s.l.Infof("Create - chatID: %s, authorID: %s;", chatID, authorID)
	return s.r.Create(ctx, chatID, authorID, f)
}

func (s *service) List(ctx context.Context, page, limit int) ([]*Puk, error) {
	s.l.Infof("List - page: %s, limit %s;", page, limit)
	return s.r.List(ctx, page, limit)
}
