# Building a REST API

Learn how to build a complete RESTful API with Goca from scratch.

## What We'll Build

A **Blog API** with:
- User authentication and authorization
- Posts with CRUD operations
- Comments system
- Tags and categories
- Pagination and filtering
- Image uploads

## Prerequisites

- Goca installed
- PostgreSQL or MySQL
- Go 1.21+
- Basic REST API knowledge

## Step 1: Project Setup

```bash
mkdir blog-api
cd blog-api

goca init blog --module github.com/yourusername/blog --database postgres --auth
cd blog
```

The `--auth` flag includes JWT authentication scaffolding.

## Step 2: User Authentication

The authentication system is already generated. Let's configure it:

Edit `internal/usecase/auth_service.go`:

```go
package usecase

import (
    "context"
    "errors"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    "time"
)

type AuthService interface {
    Register(ctx context.Context, input RegisterInput) (*AuthResponse, error)
    Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
    ValidateToken(tokenString string) (*Claims, error)
}

type authService struct {
    userRepo repository.UserRepository
    jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
    return &authService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
    }
}

type RegisterInput struct {
    Name     string `json:"name" binding:"required,min=2"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}

type Claims struct {
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
    // Check if user exists
    existing, _ := s.userRepo.FindByEmail(ctx, input.Email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    
    // Create user
    user := &domain.User{
        Name:     input.Name,
        Email:    input.Email,
        Password: string(hashedPassword),
    }
    
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // Generate token
    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        Token: token,
        User:  toUserResponse(user),
    }, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
    user, err := s.userRepo.FindByEmail(ctx, input.Email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Verify password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    
    // Generate token
    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        Token: token,
        User:  toUserResponse(user),
    }, nil
}

func (s *authService) generateToken(user *domain.User) (string, error) {
    claims := Claims{
        UserID: user.ID,
        Email:  user.Email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.jwtSecret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

## Step 3: Create Post Feature

```bash
goca feature Post title:string content:text author_id:uint:fk published:bool views:int
```

Add relationships in `internal/domain/post.go`:

```go
type Post struct {
    ID        uint      `gorm:"primaryKey"`
    Title     string    `gorm:"not null;index"`
    Content   string    `gorm:"type:text"`
    AuthorID  uint      `gorm:"not null"`
    Author    User      `gorm:"foreignKey:AuthorID"`
    Published bool      `gorm:"default:false"`
    Views     int       `gorm:"default:0"`
    Comments  []Comment `gorm:"foreignKey:PostID"`
    Tags      []Tag     `gorm:"many2many:post_tags"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Business methods
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

## Step 4: Create Comment Feature

```bash
goca feature Comment post_id:uint:fk author_id:uint:fk content:text
```

## Step 5: Create Tag Feature

```bash
goca feature Tag name:string slug:string
```

## Step 6: Add Authentication Middleware

Create `internal/middleware/auth.go`:

```go
package middleware

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "yourproject/internal/usecase"
)

func AuthMiddleware(authService usecase.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Extract token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }
        
        token := parts[1]
        
        // Validate token
        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Store user info in context
        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        
        c.Next()
    }
}
```

## Step 7: Update Routes

Edit `internal/handler/http/routes.go`:

```go
func SetupRoutes(router *gin.Engine, container *di.Container) {
    api := router.Group("/api/v1")
    
    // Public routes
    auth := api.Group("/auth")
    {
        auth.POST("/register", container.AuthHandler.Register)
        auth.POST("/login", container.AuthHandler.Login)
    }
    
    // Public posts (read-only)
    publicPosts := api.Group("/posts")
    {
        publicPosts.GET("", container.PostHandler.ListPosts)
        publicPosts.GET("/:id", container.PostHandler.GetPost)
    }
    
    // Protected routes
    protected := api.Group("")
    protected.Use(middleware.AuthMiddleware(container.AuthService))
    {
        // Posts management
        posts := protected.Group("/posts")
        {
            posts.POST("", container.PostHandler.CreatePost)
            posts.PUT("/:id", container.PostHandler.UpdatePost)
            posts.DELETE("/:id", container.PostHandler.DeletePost)
            posts.POST("/:id/publish", container.PostHandler.PublishPost)
        }
        
        // Comments
        comments := protected.Group("/comments")
        {
            comments.POST("", container.CommentHandler.CreateComment)
            comments.PUT("/:id", container.CommentHandler.UpdateComment)
            comments.DELETE("/:id", container.CommentHandler.DeleteComment)
        }
    }
}
```

## Step 8: Add Pagination

Edit `internal/repository/postgres_post_repository.go`:

```go
type PostFilters struct {
    Published  *bool
    AuthorID   *uint
    Tag        string
    Search     string
    Page       int
    PageSize   int
}

func (r *PostgresPostRepository) FindWithFilters(ctx context.Context, filters PostFilters) ([]*domain.Post, int64, error) {
    query := r.db.WithContext(ctx).Model(&domain.Post{})
    
    // Apply filters
    if filters.Published != nil {
        query = query.Where("published = ?", *filters.Published)
    }
    
    if filters.AuthorID != nil {
        query = query.Where("author_id = ?", *filters.AuthorID)
    }
    
    if filters.Search != "" {
        query = query.Where("title ILIKE ? OR content ILIKE ?",
            "%"+filters.Search+"%",
            "%"+filters.Search+"%")
    }
    
    if filters.Tag != "" {
        query = query.Joins("JOIN post_tags ON post_tags.post_id = posts.id").
            Joins("JOIN tags ON tags.id = post_tags.tag_id").
            Where("tags.slug = ?", filters.Tag)
    }
    
    // Count total
    var total int64
    query.Count(&total)
    
    // Pagination
    offset := (filters.Page - 1) * filters.PageSize
    
    var posts []*domain.Post
    err := query.
        Preload("Author").
        Preload("Tags").
        Offset(offset).
        Limit(filters.PageSize).
        Order("created_at DESC").
        Find(&posts).Error
    
    return posts, total, err
}
```

## Step 9: Test the API

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
    "email": "john@example.com"
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

```bash
TOKEN="your-jwt-token"

curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my first blog post.",
    "published": false
  }'
```

### List Posts with Filters

```bash
# Get published posts only
curl "http://localhost:8080/api/v1/posts?published=true&page=1&page_size=10"

# Search posts
curl "http://localhost:8080/api/v1/posts?search=golang&page=1&page_size=10"

# Filter by tag
curl "http://localhost:8080/api/v1/posts?tag=programming&page=1&page_size=10"
```

### Add a Comment

```bash
curl -X POST http://localhost:8080/api/v1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "post_id": 1,
    "content": "Great post!"
  }'
```

## Step 10: Add Image Upload

Create `internal/handler/http/upload_handler.go`:

```go
package http

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "path/filepath"
    "github.com/google/uuid"
)

type UploadHandler struct {
    uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
    return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }
    
    // Validate file type
    ext := filepath.Ext(file.Filename)
    if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
        return
    }
    
    // Generate unique filename
    filename := uuid.New().String() + ext
    path := filepath.Join(h.uploadDir, filename)
    
    // Save file
    if err := c.SaveUploadedFile(file, path); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "url": "/uploads/" + filename,
    })
}
```

## Next Steps

- [Adding Features](/tutorials/adding-features) - Extend functionality
- [Best Practices](/guide/best-practices) - Code quality
- [Complete Tutorial](/tutorials/complete-tutorial) - Advanced patterns

## Summary

You now have:
-  User authentication with JWT
-  CRUD operations for posts
-  Comments system
-  Pagination and filtering
-  Protected routes with middleware
-  Many-to-many relationships (tags)
-  Business logic in domain layer

Your REST API follows Clean Architecture and is production-ready! 
