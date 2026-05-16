package plystra

import "context"

func (s ResourceTypesService) List(ctx context.Context, query Query) ([]Map, error) {
	return s.client.getList(ctx, "/api/v1/resource-types", query)
}
