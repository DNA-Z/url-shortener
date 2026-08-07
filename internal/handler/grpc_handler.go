package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/DNA-Z/url-shortener/api/proto"
	"github.com/DNA-Z/url-shortener/internal/auth"
	"github.com/DNA-Z/url-shortener/internal/middleware"
	"github.com/DNA-Z/url-shortener/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

// GRPCHandler - gRPC обработчик для сервиса сокращения URL
type GRPCHandler struct {
	shortener.UnimplementedShortenerServiceServer
	urlService *service.URLStorage
	baseURL    string
}

// NewGRPCHandler - создает новый gRPC обработчик
func NewGRPCHandler(urlService *service.URLStorage, baseURL string) *GRPCHandler {
	return &GRPCHandler{
		urlService: urlService,
		baseURL:    baseURL,
	}
}

func (h *GRPCHandler) getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	if userID, err := middleware.GetUserIDFromContext(ctx); err == nil {
		return userID, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, errors.New("missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return uuid.Nil, errors.New("missing authorization header")
	}

	authHeader := authHeaders[0]

	const prefix = "Bearer "
	if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
		return uuid.Nil, errors.New("invalid authorization header format")
	}

	token := authHeader[len(prefix):]

	userIDStr, err := auth.GetUserID(token)
	if err != nil {
		return uuid.Nil, errors.New("invalid token")
	}

	if userIDStr == "" && !auth.IsAuthEnabled() {
		return uuid.Nil, nil
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return userID, nil
}

func (h *GRPCHandler) normalizeBaseURL() string {
	baseURL := h.baseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return baseURL
}
