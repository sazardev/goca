# Quick Start: Using .goca.yaml Configuration

**For Users Who Want to Get Started Quickly**

---

## 🚀 5-Minute Setup

### Step 1: Create `.goca.yaml` in Your Project

```yaml
project:
  name: my-awesome-api
  module: github.com/myuser/awesome-api

database:
  type: postgres
  features:
    timestamps: true
    soft_delete: true

generation:
  validation:
    enabled: true
    library: validator

architecture:
  naming:
    files: snake_case
```

### Step 2: Generate Your First Entity

```bash
goca entity Product --fields "name:string,price:float64,description:string"
```

### Step 3: Check the Generated File

```bash
cat internal/domain/product.go
```

**You should see:**
- ✅ Filename: `product.go` (snake_case from config)
- ✅ Validation tags: `validate:"required"` on fields
- ✅ Timestamps: `CreatedAt`, `UpdatedAt` fields
- ✅ Soft delete: `DeletedAt`, `SoftDelete()`, `IsDeleted()`

---

## 🎯 That's It!

All subsequent commands will use these settings automatically:

```bash
goca entity Order --fields "customer_id:int,total:float64"
goca repository Order --database postgres
goca handler Order --type http
goca usecase Order --operations "create,read,update,delete"
```

---

## 💡 Common Configurations

### Configuration for Microservices

```yaml
project:
  name: user-service
  module: github.com/company/user-service

database:
  type: postgres
  features:
    timestamps: true
    soft_delete: true
    uuid: true

generation:
  validation:
    enabled: true
  documentation:
    swagger:
      enabled: true

architecture:
  naming:
    files: snake_case

features:
  logging:
    enabled: true
    level: info
  monitoring:
    enabled: true
```

### Configuration for Monolithic Apps

```yaml
project:
  name: shop-app
  module: github.com/shop/app

database:
  type: mysql
  features:
    timestamps: true
    soft_delete: false

generation:
  validation:
    enabled: true

architecture:
  naming:
    files: PascalCase
```

### Minimal Configuration (Just the Essentials)

```yaml
project:
  name: quick-api
  module: github.com/user/quick-api

database:
  type: postgres
```

---

## ❓ FAQ

### Q: Do I need a config file?
**A:** No! GOCA works fine without it. Config is optional but recommended for consistency.

### Q: Can I override config with CLI flags?
**A:** Yes! CLI flags always take precedence:
```bash
goca entity Product --fields "name:string" --validation=false
# This disables validation even if config enables it
```

### Q: What if I make a typo in the YAML?
**A:** GOCA will show an error message and use default values. Run `goca config validate` to check.

### Q: Can I use the config in multiple projects?
**A:** Each project should have its own `.goca.yaml`. You can copy and modify as needed.

### Q: Where can I find all available config options?
**A:** See `docs/YAML_STRUCTURE_REFERENCE.md` for the complete reference.

---

## 🆘 Troubleshooting

### Config file not being loaded?
```bash
# Check if file exists and is named correctly
ls .goca.yaml

# Validate the YAML syntax
goca config validate

# View current config
goca config show
```

### Features not applying?
Make sure your YAML structure is correct:

❌ **WRONG:**
```yaml
generation:
  timestamps:
    enabled: true
```

✅ **CORRECT:**
```yaml
database:
  features:
    timestamps: true
```

See `docs/YAML_STRUCTURE_REFERENCE.md` for the complete correct structure.

---

## 📚 Learn More

- **Full Reference:** `docs/YAML_STRUCTURE_REFERENCE.md`
- **Migration Guide:** `docs/migration-guide.md`
- **Complete Tutorial:** `wiki/Complete-Tutorial.md`
- **Configuration System:** `docs/configuration-system.md`

---

## 🎉 Happy Coding!

With `.goca.yaml`, you configure once and generate consistent code across your entire team. Enjoy! 🚀
