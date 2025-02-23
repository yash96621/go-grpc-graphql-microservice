package catalog

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/olivere/elastic/v7"
)

var (
	ErrNotFound = errors.New("Entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error)
}

type elasticRepository struct {
	client *elastic.Client
}

type ProductDocument struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client: client}, nil
}

func (r *elasticRepository) Close() {
	r.client.Stop()
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
		Index("products").
		Id(p.ID).
		BodyJson(ProductDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, err := r.client.Get().
		Index("products").
		Id(id).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if !res.Found {
		return nil, ErrNotFound
	}
	var doc ProductDocument
	if err := json.Unmarshal(*&res.Source, &doc); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        doc.Name,
		Description: doc.Description,
		Price:       doc.Price,
	}, err
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("products").
		Query(elastic.NewMatchAllQuery()).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []Product
	for _, hit := range res.Hits.Hits {
		var doc ProductDocument
		if err := json.Unmarshal(*&hit.Source, &doc); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        doc.Name,
			Description: doc.Description,
			Price:       doc.Price,
		})
	}
	return products, nil
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	var products []Product
	for _, id := range ids {
		p, err := r.GetProductById(ctx, id)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("products").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, hit := range res.Hits.Hits {
		var doc ProductDocument
		if err := json.Unmarshal(*&hit.Source, &doc); err != nil {
			return nil, err
		}
		products = append(products, Product{
			ID:          hit.Id,
			Name:        doc.Name,
			Description: doc.Description,
			Price:       doc.Price,
		})
	}
	return products, nil
}
