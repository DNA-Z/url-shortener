package handler

import (
	"context"

	"github.com/DNA-Z/url-shortener/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExpandURL - возвращает оригинальный URL по ID (соответствует GET /{id})
func (h *GRPCHandler) ExpandURL(ctx context.Context, req *shortener.URLExpandRequest) (*shortener.URLExpandResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	_, err := h.getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	urlData, err := h.urlService.GetOriginURLByShortURL(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "url not found: %v", err)
	}

	if urlData.IsDeleted {
		return nil, status.Error(codes.NotFound, "url has been deleted")
	}

	return &shortener.URLExpandResponse{
		Result: urlData.OriginalURL,
	}, nil
}
