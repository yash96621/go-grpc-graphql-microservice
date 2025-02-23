package main

import "context"

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	return r.server, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, id *string) ([]*Product, error) {
	return r.server, nil
}
