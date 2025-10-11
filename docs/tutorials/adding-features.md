# Adding Features

Learn how to extend your Goca projects with new features and functionality.

## Adding a New Feature

### 1. Generate the Feature

```bash
goca feature Comment user_id:uint:fk post_id:uint:fk content:text
```

This automatically:
- Creates entity in `internal/domain/`
- Generates repository in `internal/repository/`
- Creates service in `internal/usecase/`
- Generates handler in `internal/handler/http/`
- Updates DI container
- Registers routes

### 2. Add Relationships

Edit the domain entity to add relationships:

```go
// internal/domain/comment.go
type Comment struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null"`
    User      User      `gorm:"foreignKey:UserID"`
    PostID    uint      `gorm:"not null"`
    Post      Post      `gorm:"foreignKey:PostID"`
    Content   string    `gorm:"type:text;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 3. Customize Business Logic

Add custom methods to the service:

```go
// internal/usecase/comment_service.go

func (s *commentService) GetCommentsByPost(ctx context.Context, postID uint) ([]*CommentResponse, error) {
    comments, err := s.repo.FindByPostID(ctx, postID)
    if err != nil {
        return nil, fmt.Errorf("failed to get comments: %w", err)
    }
    
    var responses []*CommentResponse
    for _, comment := range comments {
        responses = append(responses, toCommentResponse(comment))
    }
    return responses, nil
}
```

### 4. Add Custom Repository Methods

```go
// internal/repository/postgres_comment_repository.go

func (r *PostgresCommentRepository) FindByPostID(ctx context.Context, postID uint) ([]*domain.Comment, error) {
    var comments []*domain.Comment
    err := r.db.WithContext(ctx).
        Where("post_id = ?", postID).
        Preload("User").
        Order("created_at DESC").
        Find(&comments).Error
    return comments, err
}
```

### 5. Add Custom HTTP Endpoints

```go
// internal/handler/http/comment_handler.go

func (h *CommentHandler) GetCommentsByPost(c *gin.Context) {
    postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }
    
    comments, err := h.service.GetCommentsByPost(c.Request.Context(), uint(postID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    comments,
    })
}
```

### 6. Register Custom Routes

```go
// internal/handler/http/routes.go

posts := api.Group("/posts")
{
    posts.POST("", container.PostHandler.CreatePost)
    posts.GET("/:id", container.PostHandler.GetPost)
    posts.GET("/:post_id/comments", container.CommentHandler.GetCommentsByPost)
}
```

## Adding Custom Middleware

### 1. Create Middleware

```go
// internal/middleware/auth.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization"})
            c.Abort()
            return
        }
        
        // Validate token here
        userID, err := validateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}
```

### 2. Apply Middleware

```go
// internal/handler/http/routes.go

protected := api.Group("/")
protected.Use(middleware.AuthMiddleware())
{
    protected.POST("/posts", container.PostHandler.CreatePost)
    protected.DELETE("/posts/:id", container.PostHandler.DeletePost)
}
```

## Adding Validation

### 1. Use Struct Tags

```go
// internal/usecase/dto.go

type CreatePostInput struct {
    Title   string   `json:"title" binding:"required,min=3,max=200"`
    Content string   `json:"content" binding:"required,min=10"`
    Tags    []string `json:"tags" binding:"max=10"`
}
```

### 2. Custom Validators

```go
// internal/handler/http/validators.go
package http

import (
    "github.com/go-playground/validator/v10"
    "regexp"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateEmail(fl validator.FieldLevel) bool {
    return emailRegex.MatchString(fl.Field().String())
}

// Register custom validator
func RegisterCustomValidators(v *validator.Validate) {
    v.RegisterValidation("custom_email", validateEmail)
}
```

## Adding Search/Filtering

### 1. Add Filter DTOs

```go
// internal/usecase/dto.go

type ListPostsInput struct {
    Search   string `form:"search"`
    Tags     string `form:"tags"`
    Page     int    `form:"page" binding:"min=1"`
    PageSize int    `form:"page_size" binding:"min=1,max=100"`
}
```

### 2. Implement Repository Method

```go
// internal/repository/postgres_post_repository.go

func (r *PostgresPostRepository) Search(ctx context.Context, input ListPostsInput) ([]*domain.Post, error) {
    query := r.db.WithContext(ctx)
    
    // Search by title or content
    if input.Search != "" {
        query = query.Where("title ILIKE ? OR content ILIKE ?",
            "%"+input.Search+"%",
            "%"+input.Search+"%")
    }
    
    // Filter by tags
    if input.Tags != "" {
        tags := strings.Split(input.Tags, ",")
        query = query.Where("tags && ?", pq.Array(tags))
    }
    
    // Pagination
    offset := (input.Page - 1) * input.PageSize
    
    var posts []*domain.Post
    err := query.
        Offset(offset).
        Limit(input.PageSize).
        Order("created_at DESC").
        Find(&posts).Error
    
    return posts, err
}
```

## Adding File Upload

### 1. Create Upload Handler

```go
// internal/handler/http/upload_handler.go
package http

import (
    "github.com/gin-gonic/gin"
    "path/filepath"
    "net/http"
)

type UploadHandler struct {
    uploadDir string
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }
    
    // Validate file type
    ext := filepath.Ext(file.Filename)
    if !isAllowedExtension(ext) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
        return
    }
    
    // Save file
    filename := generateUniqueFilename(file.Filename)
    path := filepath.Join(h.uploadDir, filename)
    
    if err := c.SaveUploadedFile(file, path); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success":  true,
        "filename": filename,
        "url":      "/uploads/" + filename,
    })
}
```

### 2. Register Upload Route

```go
api.POST("/upload", container.UploadHandler.UploadFile)
api.Static("/uploads", "./uploads")
```

## Adding Background Jobs

### 1. Create Job Service

```go
// internal/usecase/email_service.go
package usecase

import (
    "context"
    "fmt"
)

type EmailService interface {
    SendWelcomeEmail(ctx context.Context, userID uint) error
}

type emailService struct {
    userRepo repository.UserRepository
}

func (s *emailService) SendWelcomeEmail(ctx context.Context, userID uint) error {
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return err
    }
    
    // Send email logic here
    fmt.Printf("Sending welcome email to %s\n", user.Email)
    
    return nil
}
```

### 2. Use Goroutines for Async

```go
// internal/usecase/user_service.go

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*UserResponse, error) {
    user := &domain.User{
        Name:  input.Name,
        Email: input.Email,
    }
    
    if err := s.repo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Send welcome email asynchronously
    go func() {
        _ = s.emailService.SendWelcomeEmail(context.Background(), user.ID)
    }()
    
    return toUserResponse(user), nil
}
```

## See Also

- [Complete Tutorial](/tutorials/complete-tutorial) - Build from scratch
- [Best Practices](/guide/best-practices) - Code guidelines
- [Project Structure](/guide/project-structure) - Directory organization
