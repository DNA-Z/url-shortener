package handler

import (
	"context"

	"github.com/DNA-Z/url-shortener/api/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ListUserURLs - возвращает все URL пользователя (соответствует GET /api/user/urls)
func (h *GRPCHandler) ListUserURLs(ctx context.Context, _ *emptypb.Empty) (*shortener.UserURLsResponse, error) {
	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	urls, err := h.urlService.GetUserURLsByUserID(userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user URLs: %v", err)
	}

	protoURLs := make([]*shortener.URLData, len(urls))
	baseURL := h.normalizeBaseURL()
	for i, url := range urls {
		protoURLs[i] = &shortener.URLData{
			ShortUrl:    baseURL + url.ShortURL,
			OriginalUrl: url.OriginalURL,
		}
	}

	return &shortener.UserURLsResponse{
		Url: protoURLs,
	}, nil
}
