package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Blue-Onion/ArtmeisterBackend/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type mockAdminRepo struct {
	database.UserRepository
	database.ArtRepository
	users     map[uuid.UUID]database.User
	arts      map[uuid.UUID]database.Art
	patchErr  error
	statusErr error
}

func (m *mockAdminRepo) PatchUserAdmin(ctx context.Context, arg database.PatchUserAdminParams) (database.PatchUserAdminRow, error) {
	if m.patchErr != nil {
		return database.PatchUserAdminRow{}, m.patchErr
	}
	u, ok := m.users[arg.ID]
	if !ok {
		return database.PatchUserAdminRow{}, sql.ErrNoRows
	}
	if arg.Role.Valid {
		u.Role = arg.Role.UserRole
	}
	if arg.Status.Valid {
		u.Status = arg.Status.AccountStatus
	}
	m.users[arg.ID] = u
	return database.PatchUserAdminRow{
		ID:     u.ID,
		Status: u.Status,
		Role:   u.Role,
	}, nil
}

func (m *mockAdminRepo) UpdateArtStatus(ctx context.Context, arg database.UpdateArtStatusParams) (database.UpdateArtStatusRow, error) {
	if m.statusErr != nil {
		return database.UpdateArtStatusRow{}, m.statusErr
	}
	a, ok := m.arts[arg.ID]
	if !ok {
		return database.UpdateArtStatusRow{}, sql.ErrNoRows
	}
	a.Status = arg.Status
	m.arts[arg.ID] = a
	return database.UpdateArtStatusRow{
		ID:     a.ID,
		Status: a.Status,
	}, nil
}

func newMockAdminRepo() *mockAdminRepo {
	return &mockAdminRepo{
		users: make(map[uuid.UUID]database.User),
		arts:  make(map[uuid.UUID]database.Art),
	}
}

func TestHandlerUserStatus(t *testing.T) {
	repo := newMockAdminRepo()
	h := &UserHandler{Repo: repo}

	userUUID := uuid.New()
	repo.users[userUUID] = database.User{
		ID:     userUUID,
		Name:   "Alice",
		Role:   database.UserRoleUser,
		Status: database.AccountStatusPending,
	}

	tests := []struct {
		name           string
		userIDParam    string
		queryRole      string
		queryStatus    string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success Approve User",
			userIDParam:    userUUID.String(),
			queryStatus:    "approved",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Success Change Role",
			userIDParam:    userUUID.String(),
			queryRole:      "admin",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Status Provided",
			userIDParam:    userUUID.String(),
			queryStatus:    "invalid-status",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Both Empty",
			userIDParam:    userUUID.String(),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Both Provided Fails",
			userIDParam:    userUUID.String(),
			queryRole:      "user",
			queryStatus:    "approved",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Repo Error Handle",
			userIDParam:    userUUID.String(),
			queryStatus:    "approved",
			mockErr:        errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid UUID Param",
			userIDParam:    "not-a-uuid",
			queryStatus:    "approved",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo.patchErr = tc.mockErr

			query := fmt.Sprintf("/admin/users/%s/status", tc.userIDParam)
			if tc.queryRole != "" || tc.queryStatus != "" {
				query += "?"
				if tc.queryRole != "" {
					query += "role=" + tc.queryRole
				}
				if tc.queryStatus != "" {
					if tc.queryRole != "" {
						query += "&"
					}
					query += "status=" + tc.queryStatus
				}
			}
			req := httptest.NewRequest(http.MethodPatch, query, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("user_id", tc.userIDParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandlerRole(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tc.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}

func TestHandlerArtStatus(t *testing.T) {
	repo := newMockAdminRepo()
	h := &ArtHandler{Repo: repo}

	artUUID := uuid.New()
	repo.arts[artUUID] = database.Art{
		ID:     artUUID,
		Name:   "Artwork",
		Status: database.ArtStatusPending,
	}

	tests := []struct {
		name           string
		artIDParam     string
		queryStatus    string
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success Approve Art",
			artIDParam:     artUUID.String(),
			queryStatus:    "approved",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid Query Status",
			artIDParam:     artUUID.String(),
			queryStatus:    "invalid-status",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Repo Failure",
			artIDParam:     artUUID.String(),
			queryStatus:    "approved",
			mockErr:        errors.New("db write failed"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid Art UUID",
			artIDParam:     "not-a-uuid",
			queryStatus:    "approved",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty Status Param",
			artIDParam:     artUUID.String(),
			queryStatus:    "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo.statusErr = tc.mockErr

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/admin/arts/%s/status?status=%s", tc.artIDParam, tc.queryStatus), nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("art_id", tc.artIDParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandlerArtStatus(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d: %s", tc.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}
