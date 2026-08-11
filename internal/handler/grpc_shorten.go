package handler

import (
	"context"
	"errors"

	"github.com/DNA-Z/url-shortener/api/proto"
	cerrors "github.com/DNA-Z/url-shortener/internal/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShortenURL - сокращает URL (соответствует POST /api/shorten)
func (h *GRPCHandler) ShortenURL(ctx context.Context, req *shortener.URLShortenRequest) (*shortener.URLShortenResponse, error) {
	if req.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is required")
	}

	userID, err := h.getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication required: %v", err)
	}

	shortURL, err := h.urlService.Shorten(userID, req.Url)
	if err != nil {
		var conflictErr *cerrors.ConflictError
		if errors.As(err, &conflictErr) {
			return &shortener.URLShortenResponse{
				Result: h.normalizeBaseURL() + shortURL,
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "failed to shorten URL: %v", err)
	}

	fullShortURL := h.normalizeBaseURL() + shortURL

	return &shortener.URLShortenResponse{
		Result: fullShortURL,
	}, nil
}
