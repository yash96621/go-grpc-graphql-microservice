package main

import "context"

type AccountResolver struct {
	server *Server
}

func (r *AccountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	return r.server, nil
}
