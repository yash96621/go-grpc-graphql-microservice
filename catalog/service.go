package catalog

import (
	"context"

	"github.com/segmentio/ksuid"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetProductsByIds(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CatalogService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &CatalogService{r}
}

func (s *CatalogService) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		ID:          ksuid.New().String(),
	}
	if err := s.repository.PutProduct(ctx, *p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *CatalogService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return s.repository.GetProductById(ctx, id)
}

func (s *CatalogService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListProducts(ctx, skip, take)

}

func (s *CatalogService) GetProductsByIds(ctx context.Context, ids []string) ([]Product, error) {

	return s.repository.ListProductsWithIDs(ctx, ids)
}

func (s *CatalogService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.SearchProducts(ctx, query, skip, take)
}
