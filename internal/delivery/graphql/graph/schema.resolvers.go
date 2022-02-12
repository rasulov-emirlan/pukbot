package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/rasulov-emirlan/pukbot/internal/delivery/graphql/graph/generated"
	"github.com/rasulov-emirlan/pukbot/internal/puk"
)

func (r *pukResolver) CreatedAt(ctx context.Context, obj *puk.Puk) (string, error) {
	return obj.CreatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"), nil
}

func (r *pukResolver) UpdatedAt(ctx context.Context, obj *puk.Puk) (string, error) {
	return obj.UpdatedAt.Format("2006-01-02T15:04:05.999999999Z07:00"), nil
}

func (r *queryResolver) Puks(ctx context.Context, limit *int, page *int) ([]*puk.Puk, error) {
	return r.PukService.List(ctx, *page, *limit)
}

// Puk returns generated.PukResolver implementation.
func (r *Resolver) Puk() generated.PukResolver { return &pukResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type pukResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
