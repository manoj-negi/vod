// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error)
	DeleteUser(ctx context.Context, id int32) error
	GetUser(ctx context.Context, id int32) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	ListVideos(ctx context.Context, name string) ([]ListVideosRow, error)
}

var _ Querier = (*Queries)(nil)