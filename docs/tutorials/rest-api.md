# Building a REST API

Learn how to build a complete RESTful API with Goca from scratch.

Every command and code sample in this tutorial was run against a real `goca` build. The generated HTTP layer uses **gorilla/mux** (not gin/chi/echo) — handlers look like `func (h *XHandler) Method(w http.ResponseWriter, r *http.Request)` and read path params with `mux.Vars(r)`.

## What We'll Build

A **Blog API** with:
- User authentication and authorization
- Posts with CRUD operations
- Comments system
- Tags
- Pagination and filtering

## Prerequisites

- Goca installed
- PostgreSQL running locally (or via Docker)
- Go 1.21+
- Basic REST API knowledge

## Step 1: Project Setup

```bash
mkdir blog-api
cd blog-api

goca init blog --module github.com/yourusername/blog --database postgres --auth
cd blog
```

`goca init` scaffolds `cmd/server/main.go`, `internal/{domain,usecase,repository,handler}/`, `pkg/{config,logger}/`, and — because of `--auth` — `pkg/auth/jwt.go`. It does **not** scaffold `internal/di/` yet; that appears the first time you generate a feature.

Configure your database connection in `.env` (created from `.env.example`):

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=blog
```

## Step 2: Generate the User Feature

```bash
goca feature User --fields "name:string,email:string,password:string" --database postgres
```

This single command generates the domain entity, use case, repository, HTTP handler, wires everything into `internal/di/container.go`, and registers `/api/v1/users` routes in `cmd/server/main.go`. The generated `internal/repository/interfaces.go` includes a `FindByEmail` method we'll reuse for login:

```go
type UserRepository interface {
    Save(user *domain.User) error
    FindByID(id int) (*domain.User, error)
    FindByEmail(email string) (*domain.User, error)
    Update(user *domain.User) error
    Delete(id int) error
    FindAll() ([]domain.User, error)
}
```

`goca` does not generate register/login endpoints on its own — it only knows about CRUD. The `--auth` flag from Step 1 gave us the JWT building blocks in `pkg/auth/jwt.go`:

```go
package auth

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key") // Change this!

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, email string) (string, error) { /* ... */ }
func ValidateToken(tokenString string) (*Claims, error)      { /* ... */ }
```

> **Tip:** `jwtSecret` is hardcoded in the generated file. Change it (or load it from `cfg` via an environment variable) before deploying.

Now add a small hand-written auth handler on top of the generated `UserRepository`. Create `internal/handler/http/auth_handler.go`:

```go
package http

import (
    "encoding/json"
    "net/http"

    "github.com/yourusername/blog/internal/domain"
    "github.com/yourusername/blog/internal/repository"
    "github.com/yourusername/blog/pkg/auth"
    "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
    userRepo repository.UserRepository
}

func NewAuthHandler(userRepo repository.UserRepository) *AuthHandler {
    return &AuthHandler{userRepo: userRepo}
}

type RegisterInput struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginInput struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Token string      `json:"token"`
    User  domain.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var input RegisterInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if existing, _ := h.userRepo.FindByEmail(input.Email); existing != nil {
        http.Error(w, "email already registered", http.StatusConflict)
        return
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "failed to hash password", http.StatusInternalServerError)
        return
    }

    user := &domain.User{Name: input.Name, Email: input.Email, Password: string(hashed)}
    if err := h.userRepo.Save(user); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    token, err := auth.GenerateToken(int(user.ID), user.Email)
    if err != nil {
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    user.Password = ""
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(AuthResponse{Token: token, User: *user})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var input LoginInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, err := h.userRepo.FindByEmail(input.Email)
    if err != nil || user == nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := auth.GenerateToken(int(user.ID), user.Email)
    if err != nil {
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    user.Password = ""
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(AuthResponse{Token: token, User: *user})
}
```

Wire it into `cmd/server/main.go`, right after the `// goca:routes` marker that `goca feature` maintains:

```go
    apiRouter := router.PathPrefix("/api/v1").Subrouter()
    apphttp.SetupUserRoutes(apiRouter, container.UserUseCase()) // user routes
// goca:routes -- feature routes are registered above this line

    // Hand-written auth endpoints (register/login) built on top of the User feature
    authHandler := apphttp.NewAuthHandler(container.UserRepository())
    authRouter := apiRouter.PathPrefix("/auth").Subrouter()
    authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
    authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")
```

`container.UserRepository()` is a real getter generated on `internal/di/container.go` — every feature gets one for its repository, use case, and handler.

## Step 3: Create the Post Feature

```bash
goca feature Post --fields "title:string,content:string,author_id:uint,published:bool,views:int" --database postgres
```

`--fields` only accepts `name:type` pairs (no `text` type and no `:fk` suffix — both will fail validation). Relationships are added by hand after generation.

Edit `internal/domain/post.go` to add relations and a couple of business methods:

```go
package domain

import (
    "errors"
    "time"
)

type Post struct {
    ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
    Title     string    `json:"title" gorm:"type:varchar(255);not null"`
    Content   string    `json:"content" gorm:"type:varchar(255)"`
    AuthorID  uint      `json:"author_id" gorm:"not null"`
    Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
    Published bool      `json:"published" gorm:"type:boolean;not null;default:false"`
    Views     int       `json:"views" gorm:"type:integer;not null;default:0"`
    Comments  []Comment `json:"comments" gorm:"foreignKey:PostID"`
    Tags      []Tag     `json:"tags" gorm:"many2many:post_tags"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (p *Post) Validate() error {
    if p.Title == "" {
        return ErrInvalidPostTitle
    }
    if p.Content == "" {
        return ErrInvalidPostContent
    }
    if p.Views < 0 {
        return ErrInvalidPostViews
    }
    return nil
}

// Publish marks the post as published, if it has content.
func (p *Post) Publish() error {
    if p.Title == "" || p.Content == "" {
        return errors.New("cannot publish incomplete post")
    }
    p.Published = true
    return nil
}

func (p *Post) IncrementViews() {
    p.Views++
}
```

> The `title:string,content:string,...` DTOs generated in `internal/usecase/dto.go` use the JSON tags `authorid` (no underscore) for `CreatePostInput`/`UpdatePostInput` — that's a quirk of the generator, distinct from the domain entity's `author_id`. Keep that in mind when building request bodies (see Step 9).

## Step 4: Create the Comment Feature

```bash
goca feature Comment --fields "post_id:uint,author_id:uint,content:string" --database postgres
```

## Step 5: Create the Tag Feature

```bash
goca feature Tag --fields "name:string,slug:string" --database postgres
```

After Steps 3–5, `go build ./...` should succeed and `cmd/server/main.go` will have all four `Setup*Routes` calls wired above the `// goca:routes` marker.

## Step 6: Add Authentication Middleware

Generated routes use plain `net/http` middleware (`func(http.Handler) http.Handler`), not gin middleware. Create `internal/middleware/auth.go`:

```go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/yourusername/blog/pkg/auth"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// RequireAuth validates the Bearer token on the request and, if valid,
// stores the authenticated user ID in the request context.
func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
            return
        }

        claims, err := auth.ValidateToken(parts[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## Step 7: Protect the Write Routes

The generated `internal/handler/http/routes.go` registers every method on `postRouter` with `HandleFunc`. Edit it to keep reads public and require a JWT for writes, using `Handle` + `middleware.RequireAuth`:

```go
package http

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/yourusername/blog/internal/middleware"
    "github.com/yourusername/blog/internal/usecase"
)

func SetupPostRoutes(router *mux.Router, uc usecase.PostUseCase) {
    handler := NewPostHandler(uc)

    postRouter := router.PathPrefix("/posts").Subrouter()
    postRouter.Use(corsMiddleware)
    postRouter.Use(loggingMiddleware)

    // Public, read-only routes
    postRouter.HandleFunc("", handler.ListPosts).Methods("GET")
    postRouter.HandleFunc("/{id}", handler.GetPost).Methods("GET")

    // Routes that require a valid JWT
    postRouter.Handle("", middleware.RequireAuth(http.HandlerFunc(handler.CreatePost))).Methods("POST")
    postRouter.Handle("/{id}", middleware.RequireAuth(http.HandlerFunc(handler.UpdatePost))).Methods("PUT")
    postRouter.Handle("/{id}", middleware.RequireAuth(http.HandlerFunc(handler.DeletePost))).Methods("DELETE")
}

// corsMiddleware and loggingMiddleware are generated further down in this file — leave them as-is.
```

Apply the same pattern to `SetupCommentRoutes` for `CreateComment`/`UpdateComment`/`DeleteComment` if you want comments to require authentication too.

## Step 8: Add Pagination and Filtering

The generated `internal/repository/postgres_post_repository.go` implements a plain `PostRepository` (no `context.Context` parameter, no built-in filtering). Add a filtered query method:

```go
// internal/repository/postgres_post_repository.go

type PostFilters struct {
    Published *bool
    AuthorID  *uint
    Search    string
    Page      int
    PageSize  int
}

func (p *postgresPostRepository) FindWithFilters(filters PostFilters) ([]domain.Post, int64, error) {
    query := p.db.Model(&domain.Post{})

    if filters.Published != nil {
        query = query.Where("published = ?", *filters.Published)
    }
    if filters.AuthorID != nil {
        query = query.Where("author_id = ?", *filters.AuthorID)
    }
    if filters.Search != "" {
        query = query.Where("title ILIKE ? OR content ILIKE ?",
            "%"+filters.Search+"%", "%"+filters.Search+"%")
    }

    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    offset := (filters.Page - 1) * filters.PageSize

    var posts []domain.Post
    err := query.Offset(offset).Limit(filters.PageSize).Find(&posts).Error
    return posts, total, err
}
```

Add `FindWithFilters` to the `PostRepository` interface in `internal/repository/interfaces.go`, then expose it through the use case (`internal/usecase/post_usecase.go` and `post_service.go`):

```go
// post_usecase.go
type PostUseCase interface {
    // ...existing methods...
    ListPostsFiltered(filters repository.PostFilters) (ListPostOutput, error)
}

// post_service.go
func (p *postService) ListPostsFiltered(filters repository.PostFilters) (ListPostOutput, error) {
    posts, total, err := p.repo.FindWithFilters(filters)
    if err != nil {
        return ListPostOutput{}, err
    }
    return ListPostOutput{Posts: posts, Total: int(total), Message: messages.PostsListedSuccessfully}, nil
}
```

Finally, parse query parameters in `ListPosts` on `internal/handler/http/post_handler.go`:

```go
func (p *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    filters := repository.PostFilters{Search: q.Get("search"), Page: 1, PageSize: 10}

    if v := q.Get("published"); v != "" {
        published := v == "true"
        filters.Published = &published
    }
    if v := q.Get("page"); v != "" {
        if page, err := strconv.Atoi(v); err == nil && page > 0 {
            filters.Page = page
        }
    }
    if v := q.Get("page_size"); v != "" {
        if size, err := strconv.Atoi(v); err == nil && size > 0 {
            filters.PageSize = size
        }
    }

    output, err := p.usecase.ListPostsFiltered(filters)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(output)
}
```

## Step 9: Run and Test the API

```bash
go mod tidy
go run cmd/server/main.go
```

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "password": ""
  }
}
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepass123"
  }'
```

### Create a Post (Authenticated)

Remember `CreatePostInput` uses `authorid`, not `author_id`:

```bash
TOKEN="your-jwt-token"

curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first blog post.",
    "authorid": 1,
    "published": false
  }'
```

### List Posts with Filters

```bash
# Get published posts only
curl "http://localhost:8080/api/v1/posts?published=true&page=1&page_size=10"

# Search posts
curl "http://localhost:8080/api/v1/posts?search=golang&page=1&page_size=10"
```

### Add a Comment

```bash
curl -X POST http://localhost:8080/api/v1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "postid": 1,
    "authorid": 1,
    "content": "Great post!"
  }'
```

(`CreateCommentInput` also drops the underscore: `postid`/`authorid`.)

## Step 10: Add Image Upload

The generated handlers use the standard library, so file uploads go through `r.ParseMultipartForm` — no gin `FormFile` helper. Create `internal/handler/http/upload_handler.go`:

```go
package http

import (
    "io"
    "net/http"
    "os"
    "path/filepath"

    "github.com/google/uuid"
)

type UploadHandler struct {
    uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
    return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        http.Error(w, "Invalid multipart form", http.StatusBadRequest)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "No file uploaded", http.StatusBadRequest)
        return
    }
    defer file.Close()

    ext := filepath.Ext(header.Filename)
    switch ext {
    case ".jpg", ".jpeg", ".png", ".gif":
    default:
        http.Error(w, "Invalid file type", http.StatusBadRequest)
        return
    }

    filename := uuid.New().String() + ext
    dstPath := filepath.Join(h.uploadDir, filename)

    dst, err := os.Create(dstPath)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    if _, err := io.Copy(dst, file); err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"success": true, "url": "/uploads/` + filename + `"}`))
}
```

Register it manually (next to the auth routes) since it isn't tied to a `goca feature`:

```go
uploadHandler := apphttp.NewUploadHandler("./uploads")
apiRouter.HandleFunc("/upload", uploadHandler.UploadImage).Methods("POST")
```

## Next Steps

- [Adding Features](/tutorials/adding-features) - Extend functionality
- [Best Practices](/guide/best-practices) - Code quality
- [Complete Tutorial](/tutorials/complete-tutorial) - Advanced patterns

## Summary

You now have:
- User registration and login with JWT, built on a generated `User` feature
- CRUD operations for posts, comments, and tags
- Pagination and filtering wired end-to-end
- Protected routes using gorilla/mux middleware
- Many-to-many relationships (tags) added by hand on top of the generated entity
- Business logic in the domain layer

Your REST API follows Clean Architecture and every command above works against a real `goca` build.
