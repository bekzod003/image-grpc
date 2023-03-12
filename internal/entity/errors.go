package entity

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrLimitReached = status.Error(codes.PermissionDenied, "limit reached")
)
