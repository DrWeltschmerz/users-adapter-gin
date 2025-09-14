package ginadapter_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ginadapter "github.com/DrWeltschmerz/users-adapter-gin/ginadapter"
	core "github.com/DrWeltschmerz/users-core"
	"github.com/gin-gonic/gin"
)

// mock implementations for UserRepository, RoleRepository, PasswordHasher, Tokenizer
// ... (for brevity, use in-memory or minimal mocks)

type mockUserRepo struct {
	users map[string]*core.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*core.User)}
}

func (m *mockUserRepo) Create(_ context.Context, user core.User) (*core.User, error) {
	user.ID = "1"
	m.users[user.ID] = &user
	return &user, nil
}
func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*core.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, core.ErrUserNotFound
}
func (m *mockUserRepo) GetByID(_ context.Context, id string) (*core.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, core.ErrUserNotFound
	}
	return u, nil
}
func (m *mockUserRepo) Update(_ context.Context, user core.User) (*core.User, error) {
	m.users[user.ID] = &user
	return &user, nil
}
func (m *mockUserRepo) List(_ context.Context) ([]core.User, error) {
	var out []core.User
	for _, u := range m.users {
		out = append(out, *u)
	}
	return out, nil
}
func (m *mockUserRepo) Delete(_ context.Context, id string) error { delete(m.users, id); return nil }
func (m *mockUserRepo) GetByUsername(_ context.Context, username string) (*core.User, error) {
	return nil, core.ErrUserNotFound
}

// ... similar minimal mocks for RoleRepository, PasswordHasher, Tokenizer

type mockHasher struct{}

func (m *mockHasher) Hash(pw string) (string, error) { return "hashed-" + pw, nil }
func (m *mockHasher) Verify(hashed, pw string) bool  { return hashed == "hashed-"+pw }

type mockTokenizer struct{}

func (m *mockTokenizer) GenerateToken(email, userID string) (string, error) { return "token", nil }
func (m *mockTokenizer) ValidateToken(token string) (string, error) {
	if token == "token" {
		return "1", nil
	}
	return "", core.ErrInvalidCredentials
}

type mockRoleRepo struct{}

func (m *mockRoleRepo) Create(_ context.Context, role core.Role) (*core.Role, error) {
	role.ID = "role1"
	return &role, nil
}
func (m *mockRoleRepo) GetByID(_ context.Context, id string) (*core.Role, error) {
	return &core.Role{ID: id, Name: "user"}, nil
}
func (m *mockRoleRepo) GetByName(_ context.Context, name string) (*core.Role, error) {
	return &core.Role{ID: "role1", Name: name}, nil
}
func (m *mockRoleRepo) Update(_ context.Context, role core.Role) (*core.Role, error) {
	return &role, nil
}
func (m *mockRoleRepo) List(_ context.Context) ([]core.Role, error) {
	return []core.Role{{ID: "role1", Name: "user"}}, nil
}
func (m *mockRoleRepo) Delete(_ context.Context, id string) error { return nil }

func setupRouter() *gin.Engine {
	userRepo := newMockUserRepo()
	roleRepo := &mockRoleRepo{}
	hasher := &mockHasher{}
	tokenizer := &mockTokenizer{}
	service := core.NewService(userRepo, roleRepo, hasher, tokenizer)
	r := gin.Default()
	ginadapter.RegisterRoutes(r, service, tokenizer)
	return r
}

func TestRegisterAndLogin(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := map[string]string{"email": "test@example.com", "username": "testuser", "password": "secret"}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(b))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	loginBody := map[string]string{"email": "test@example.com", "password": "secret"}
	b, _ = json.Marshal(loginBody)
	req, _ = http.NewRequest("POST", "/login", bytes.NewReader(b))
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == "" {
		t.Fatal("expected token in response")
	}
}
