# Adding Features

Learn how to extend an existing Goca project with new features. This tutorial assumes you already have
a project created with `goca init` (see the [Complete Tutorial](/tutorials/complete-tutorial) if not).

Everything below was verified against a real generated project (`goca init blogapp --database postgres`,
then `goca feature User ...`, `goca feature Post ...`, `goca feature Comment ...`) and compiles with `go build ./...`.

## Adding a New Feature

### 1. Generate the Feature

```bash
goca feature Comment --fields "user_id:uint,post_id:uint,content:string"
```

Notes on field types: goca only accepts Go types for `--fields` (`string`, `int`, `uint`, `float64`, `bool`,
`time.Time`, etc.) — there is no `text` or `fk` type. A field named `xxx_id` of an integer type is treated as a
foreign key by convention, and goca automatically generates `FindByXxxID` lookup methods for it (see step 4).

Running this command from the root of an existing project automatically:

- Creates the entity in `internal/domain/comment.go` (plus `errors.go` / `validations.go` additions)
- Generates the repository interface and implementation in `internal/repository/`
  (e.g. `postgres_comment_repository.go` for `--database postgres`; the file/type prefix matches your
  configured database — `mongo_`, `dynamodb_`, `elasticsearch_`, `sqlserver_`, etc.)
- Creates the use case (interface + service + DTOs) in `internal/usecase/`
- Generates the HTTP handler in `internal/handler/http/comment_handler.go`
- Wires the repository → use case → handler chain into `internal/di/container.go`
- Registers `SetupCommentRoutes` in `internal/handler/http/routes.go` **and** adds the call to it in
  `cmd/server/main.go`, right above the `// goca:routes` marker

After generation, run `go mod tidy` and start the server; the new endpoints are live immediately:

```bash
go mod tidy
go run cmd/server/main.go
curl http://localhost:8080/api/v1/comments
```

The generated entity looks like this — no manual wiring is needed for basic CRUD:

```go
// internal/domain/comment.go
package domain

type Comment struct {
    ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID  uint   `json:"user_id" gorm:"not null" validate:"required"`
    PostID  uint   `json:"post_id" gorm:"not null" validate:"required"`
    Content string `json:"content" gorm:"type:varchar(255)" validate:"required"`
}

func (c *Comment) Validate() error {
    if c.Content == "" {
        return ErrInvalidCommentContent
    }
    return nil
}
```

### 2. Add Relationships

goca does not auto-generate GORM association fields (`belongs to` / `has many`) — add them yourself by editing
the entity:

```go
// internal/domain/comment.go
type Comment struct {
    ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
    UserID  uint   `json:"user_id" gorm:"not null" validate:"required"`
    PostID  uint   `json:"post_id" gorm:"not null" validate:"required"`
    Content string `json:"content" gorm:"type:varchar(255)" validate:"required"`

    // Manually added associations:
    User User `json:"-" gorm:"foreignKey:UserID"`
    Post Post `json:"-" gorm:"foreignKey:PostID"`
}
```

`cmd/server/main.go` already runs `db.AutoMigrate(&domain.Comment{}, &domain.Post{}, &domain.User{})` on
startup for every generated entity, so the new foreign keys are picked up on the next run.

### 3. Customize Business Logic

The generated use case is a plain interface + struct implementation — add methods directly:

```go
// internal/usecase/comment_usecase.go
type CommentUseCase interface {
    CreateComment(input CreateCommentInput) (CreateCommentOutput, error)
    GetComment(id int) (*domain.Comment, error)
    UpdateComment(id int, input UpdateCommentInput) error
    DeleteComment(id int) error
    ListComments() (ListCommentOutput, error)

    // Custom method added by you:
    GetCommentsByPost(postID uint) ([]domain.Comment, error)
}
```

```go
// internal/usecase/comment_service.go
func (c *commentService) GetCommentsByPost(postID uint) ([]domain.Comment, error) {
    return c.repo.FindManyByPostID(postID)
}
```

### 4. Add Custom Repository Methods

Because `post_id` ends in `_id`, goca already generated a lookup method for it — but it returns a single
record (`FindByPostID`). For a "list all comments for a post" query, add your own method to both the
interface and the implementation:

```go
// internal/repository/interfaces.go
type CommentRepository interface {
    Save(comment *domain.Comment) error
    FindByID(id int) (*domain.Comment, error)
    FindByUserID(userid uint) (*domain.Comment, error)
    FindByPostID(postid uint) (*domain.Comment, error) // generated because of post_id field
    FindByContent(content string) (*domain.Comment, error)
    Update(comment *domain.Comment) error
    Delete(id int) error
    FindAll() ([]domain.Comment, error)

    // Added by you:
    FindManyByPostID(postID uint) ([]domain.Comment, error)
}
```

```go
// internal/repository/postgres_comment_repository.go
func (p *postgresCommentRepository) FindManyByPostID(postID uint) ([]domain.Comment, error) {
    var comments []domain.Comment
    result := p.db.
        Where("post_id = ?", postID).
        Order("id DESC").
        Find(&comments)
    return comments, result.Error
}
```

### 5. Add Custom HTTP Endpoints

Generated handlers use `net/http` and `gorilla/mux` — not gin. Follow the same pattern for a new endpoint:

```go
// internal/handler/http/comment_handler.go
func (c *CommentHandler) GetCommentsByPost(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    postID, err := strconv.ParseUint(vars["post_id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid post ID", http.StatusBadRequest)
        return
    }

    comments, err := c.usecase.GetCommentsByPost(uint(postID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(comments)
}
```

### 6. Register Custom Routes

Add the route inside the generated `Setup*Routes` function in `internal/handler/http/routes.go` (goca already
created this function and wired it into `main.go` for you):

```go
// internal/handler/http/routes.go
func SetupCommentRoutes(router *mux.Router, uc usecase.CommentUseCase) {
    handler := NewCommentHandler(uc)

    commentRouter := router.PathPrefix("/comments").Subrouter()
    commentRouter.Use(corsMiddleware)
    commentRouter.Use(loggingMiddleware)

    commentRouter.HandleFunc("", handler.CreateComment).Methods("POST")
    commentRouter.HandleFunc("/{id}", handler.GetComment).Methods("GET")
    commentRouter.HandleFunc("/{id}", handler.UpdateComment).Methods("PUT")
    commentRouter.HandleFunc("/{id}", handler.DeleteComment).Methods("DELETE")
    commentRouter.HandleFunc("", handler.ListComments).Methods("GET")

    // Custom route added by you:
    commentRouter.HandleFunc("/by-post/{post_id}", handler.GetCommentsByPost).Methods("GET")
}
```

No `di.Container` route group or gin `Group()` call is involved — `SetupCommentRoutes` is invoked once from
`cmd/server/main.go`:

```go
// cmd/server/main.go
apiRouter := router.PathPrefix("/api/v1").Subrouter()
apphttp.SetupUserRoutes(apiRouter, container.UserUseCase())
apphttp.SetupPostRoutes(apiRouter, container.PostUseCase())
apphttp.SetupCommentRoutes(apiRouter, container.CommentUseCase())
// goca:routes -- feature routes are registered above this line
```

## Adding Custom Middleware

Generated route setups already apply `corsMiddleware` and `loggingMiddleware` (defined once in
`internal/handler/http/routes.go`). Add your own the same way, as a `func(http.Handler) http.Handler`:

```go
// internal/handler/http/routes.go (or a new internal/middleware/auth.go)
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Missing authorization", http.StatusUnauthorized)
            return
        }

        // Validate token here.
        next.ServeHTTP(w, r)
    })
}
```

Apply it on the subrouter for the routes that need protection:

```go
func SetupCommentRoutes(router *mux.Router, uc usecase.CommentUseCase) {
    handler := NewCommentHandler(uc)

    commentRouter := router.PathPrefix("/comments").Subrouter()
    commentRouter.Use(corsMiddleware)
    commentRouter.Use(loggingMiddleware)
    commentRouter.Use(authMiddleware) // protect all comment endpoints

    commentRouter.HandleFunc("", handler.CreateComment).Methods("POST")
    // ...
}
```

## Adding Validation

Generated entities carry a hand-written `Validate() error` method plus GORM/JSON struct tags — there is no
gin `binding` tag support, because handlers don't use gin. Extend `Validate()` for extra rules:

```go
// internal/domain/comment.go
func (c *Comment) Validate() error {
    if c.Content == "" {
        return ErrInvalidCommentContent
    }
    if len(c.Content) > 2000 {
        return ErrCommentContentTooLong // add this to internal/domain/errors.go
    }
    return nil
}
```

Call it from the use case before saving (the generated `CreateComment` already does this):

```go
// internal/usecase/comment_service.go
func (c *commentService) CreateComment(input CreateCommentInput) (CreateCommentOutput, error) {
    comment := domain.Comment{
        UserID:  input.UserID,
        PostID:  input.PostID,
        Content: input.Content,
    }

    if err := comment.Validate(); err != nil {
        return CreateCommentOutput{}, err
    }
    // ...
}
```

If you want field-level validation libraries (e.g. `github.com/go-playground/validator/v10`), call
`validate.Struct(comment)` inside `Validate()` — the entity already carries `validate:"required"` tags that
library understands.

## Adding Search/Filtering

Add query-parameter parsing directly in the handler (there is no gin `form` binding), then pass parsed values
to a new use case/repository method:

```go
// internal/handler/http/post_handler.go
func (h *PostHandler) SearchPosts(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query()
    search := q.Get("search")
    page, _ := strconv.Atoi(q.Get("page"))
    if page < 1 {
        page = 1
    }

    posts, err := h.usecase.SearchPosts(search, page)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(posts)
}
```

```go
// internal/repository/postgres_post_repository.go
func (p *postgresPostRepository) Search(search string, page, pageSize int) ([]domain.Post, error) {
    query := p.db.Model(&domain.Post{})

    if search != "" {
        query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
    }

    var posts []domain.Post
    result := query.
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Order("id DESC").
        Find(&posts)
    return posts, result.Error
}
```

Register the route in `SetupPostRoutes` the same way as step 6 above, e.g.
`postRouter.HandleFunc("/search", handler.SearchPosts).Methods("GET")`.

## Adding File Upload

There is no gin `c.FormFile`/`c.SaveUploadedFile` helper generated by goca — use the standard library's
`multipart` support directly in a handler that follows the same `func(w http.ResponseWriter, r *http.Request)`
signature as every other generated handler:

```go
// internal/handler/http/upload_handler.go
package http

import (
    "io"
    "net/http"
    "os"
    "path/filepath"
)

type UploadHandler struct {
    uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
    return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "No file uploaded", http.StatusBadRequest)
        return
    }
    defer file.Close()

    dst, err := os.Create(filepath.Join(h.uploadDir, header.Filename))
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
    w.Write([]byte(`{"success":true,"filename":"` + header.Filename + `"}`))
}
```

Register it directly on the router in `cmd/server/main.go` (above the `// goca:routes` marker), and serve the
uploaded files with `http.FileServer`:

```go
uploadHandler := apphttp.NewUploadHandler("./uploads")
apiRouter.HandleFunc("/upload", uploadHandler.UploadFile).Methods("POST")
router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
```

## Adding Background Jobs

This part of the stack is framework-agnostic — plain Go works the same regardless of gin vs. gorilla/mux.
Add a service and call it asynchronously from an existing use case:

```go
// internal/usecase/email_service.go
package usecase

import (
    "fmt"

    "github.com/example/blogapp/internal/repository"
)

type EmailService interface {
    SendWelcomeEmail(userID uint) error
}

type emailService struct {
    userRepo repository.UserRepository
}

func NewEmailService(userRepo repository.UserRepository) EmailService {
    return &emailService{userRepo: userRepo}
}

func (s *emailService) SendWelcomeEmail(userID uint) error {
    user, err := s.userRepo.FindByID(int(userID))
    if err != nil {
        return err
    }

    // Send email logic here.
    fmt.Printf("Sending welcome email to %s\n", user.Email)
    return nil
}
```

```go
// internal/usecase/user_service.go
func (s *userService) CreateUser(input CreateUserInput) (CreateUserOutput, error) {
    user := domain.User{
        Name:  input.Name,
        Email: input.Email,
    }

    if err := s.repo.Save(&user); err != nil {
        return CreateUserOutput{}, err
    }

    // Send welcome email asynchronously.
    go func() {
        _ = s.emailService.SendWelcomeEmail(user.ID)
    }()

    return CreateUserOutput{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}
```

Wire `emailService` into `internal/di/container.go` the same way the generated repositories/use cases are
wired (`c.emailService = usecase.NewEmailService(c.userRepo)`), and inject it into `userService` via its
constructor.

## See Also

- [Complete Tutorial](/tutorials/complete-tutorial) - Build from scratch
- [Best Practices](/guide/best-practices) - Code guidelines
- [Project Structure](/guide/project-structure) - Directory organization
