# Database Support

Goca supports **8 production-ready database systems** with complete repository pattern implementations. Choose the database that best fits your architectural needs and use case.

## 📊 Complete Database Matrix

| Database            | Type               | Driver              | Best For                                | Generated File                         | Status  |
| ------------------- | ------------------ | ------------------- | --------------------------------------- | -------------------------------------- | ------- |
| **PostgreSQL**      | SQL                | GORM + postgres     | Traditional apps, transactional systems | `postgres_{entity}_repository.go`      | ✅ v1.0+ |
| **PostgreSQL JSON** | SQL + Document     | GORM + JSONB        | Semi-structured data, configs           | `postgres_json_{entity}_repository.go` | ✅ v2.1+ |
| **MySQL**           | SQL                | GORM + mysql        | Web apps, compatibility                 | `mysql_{entity}_repository.go`         | ✅ v1.0+ |
| **MongoDB**         | NoSQL Document     | Native driver       | Flexible schemas, rapid iteration       | `mongodb_{entity}_repository.go`       | ✅ v1.0+ |
| **SQLite**          | SQL Embedded       | database/sql        | Development, embedded apps              | `sqlite_{entity}_repository.go`        | ✅ v2.1+ |
| **SQL Server**      | SQL Enterprise     | GORM + mssql        | Enterprise systems, T-SQL               | `sqlserver_{entity}_repository.go`     | ✅ v2.1+ |
| **Elasticsearch**   | Search & Analytics | go-elasticsearch v8 | Full-text search, analytics             | `elasticsearch_{entity}_repository.go` | ✅ v2.1+ |
| **DynamoDB**        | NoSQL Serverless   | AWS SDK v2          | Serverless AWS, auto-scaling            | `dynamodb_{entity}_repository.go`      | ✅ v2.1+ |

---

## 🔍 Database Detailed Specifications

### PostgreSQL
**Version Support:** 10+  
**Driver:** `github.com/lib/pq` (via GORM)  
**Connection Pool:** Yes (configurable)  
**Transactions:** Full ACID support  
**JSONB Support:** Yes (use `postgres-json` variant)

```bash
goca feature User --database postgres
```

**Use When:**
- Building traditional business applications
- Need complex SQL queries and joins
- Require strong ACID guarantees
- Working with relational data

---

### PostgreSQL JSON
**Version Support:** 9.3+ (JSONB: 9.4+)  
**Driver:** `github.com/lib/pq` + `gorm/datatypes`  
**Special Methods:** `FindByJSONField(jsonField, value string)`

```bash
goca feature Config --database postgres-json
```

**Special Features:**
- JSONB operators (`@>`, `<@`, `?`, `||`)
- Nested object queries
- JSON array operations
- Full SQL+JSON hybrid querying

**Use When:**
- Storing hierarchical configurations
- Flexible, evolving schemas
- Semi-structured data with SQL querying
- Metadata and nested objects

---

### MySQL
**Version Support:** 5.7+, 8.0+  
**Driver:** `github.com/go-sql-driver/mysql` (via GORM)  
**Connection Pool:** Yes  
**Transactions:** InnoDB (configurable)  
**JSON Support:** Native JSON type

```bash
goca feature Product --database mysql
```

**Use When:**
- Integrating with existing MySQL infrastructure
- Building scalable web applications
- Need good performance for web workloads
- Familiar with MySQL ecosystem

---

### MongoDB
**Version Support:** 4.0+  
**Driver:** `go.mongodb.org/mongo-driver`  
**Transactions:** Multi-document ACID (4.0+)  
**Schema Flexibility:** Full  
**Replication:** Built-in

```bash
goca feature Article --database mongodb
```

**Use When:**
- Data model changes frequently
- Need flexible, document-oriented storage
- Building rapid prototypes
- Scaling horizontally with sharding
- Storing hierarchical or nested data

---

### SQLite
**Version Support:** 3.8.0+  
**Driver:** `github.com/mattn/go-sqlite3`  
**File-Based:** Single `.db` file  
**Server Required:** No  
**Perfect For:** Development, Testing, Embedded

```bash
goca feature Setting --database sqlite
```

**Features:**
- Zero configuration
- ACID transactions
- JSON marshaling for flexibility
- Cross-platform single file

**Use When:**
- Local development and testing
- Embedded applications
- Building desktop applications
- Simple CLI tools
- Prototyping before choosing production DB

---

### SQL Server
**Version Support:** 2012+  
**Driver:** `github.com/denisenkom/go-mssqldb` (via GORM)  
**T-SQL Support:** Yes  
**Enterprise Features:** All  
**Windows Integration:** Yes

```bash
goca feature Employee --database sqlserver
```

**Use When:**
- Enterprise environment with SQL Server infrastructure
- Requiring T-SQL specific features
- Migrating from legacy SQL Server systems
- Windows ecosystem integration needed
- Enterprise support required

---

### Elasticsearch
**Version Support:** 8.0+  
**Driver:** `github.com/elastic/go-elasticsearch/v8`  
**Type:** Search Engine + Analytics  
**Real-time:** Yes  
**Distributed:** Yes

**Special Methods:**
- `FullTextSearch(query string)` - Multi-field full-text search
- Lucene DSL support
- Result scoring

```bash
goca feature Article --database elasticsearch
```

**Use When:**
- Building full-text search functionality
- Need advanced analytics capabilities
- Centralizing application logs
- Building search-heavy applications (e-commerce)
- Real-time data analysis needed

**Note:** Elasticsearch is not suitable for primary data storage. Use with a main database.

---

### DynamoDB
**Version Support:** Current (AWS managed)  
**Driver:** `github.com/aws/aws-sdk-go-v2`  
**Serverless:** Yes  
**Auto-scaling:** Built-in  
**Cost Model:** Pay-per-request

```bash
goca feature Order --database dynamodb
```

**Use When:**
- Building serverless AWS applications
- Need unlimited auto-scaling
- Comfortable with eventual consistency
- Using AWS-managed infrastructure
- Cost varies with usage (good for unpredictable workloads)

---

## 🎯 Quick Selection Guide

### I'm building a...

**...Traditional Business Application**
→ Use **PostgreSQL**

**...Web Application with Existing MySQL**
→ Use **MySQL**

**...Startup with Rapidly Changing Data**
→ Use **MongoDB** or **PostgreSQL JSON**

**...E-commerce with Search**
→ Use **PostgreSQL** + **Elasticsearch**

**...Enterprise System**
→ Use **SQL Server** or **PostgreSQL**

**...Serverless AWS Application**
→ Use **DynamoDB**

**...Search-Heavy Application**
→ Use **Elasticsearch** (with primary DB)

**...Development/Testing Locally**
→ Use **SQLite**

---

## 🔄 Switching Databases

One of the key benefits of Goca is the **Repository Pattern** - you can change databases without modifying business logic:

```bash
# Start with PostgreSQL in development
goca feature User --fields "name:string,email:string" --database postgres

# Generate MongoDB repository for production scaling
goca repository User --database mongodb --implementation

# Your use cases, domain, and handlers remain unchanged!
```

---

## 📈 Performance Comparison

| Metric                  | PostgreSQL | MySQL       | MongoDB        | SQLite | SQL Server | Elasticsearch | DynamoDB  |
| ----------------------- | ---------- | ----------- | -------------- | ------ | ---------- | ------------- | --------- |
| **Read Performance**    | ⭐⭐⭐⭐       | ⭐⭐⭐⭐        | ⭐⭐⭐            | ⭐⭐⭐⭐   | ⭐⭐⭐⭐       | ⭐⭐⭐⭐⭐         | ⭐⭐⭐⭐      |
| **Write Performance**   | ⭐⭐⭐⭐       | ⭐⭐⭐⭐        | ⭐⭐⭐⭐           | ⭐⭐⭐    | ⭐⭐⭐⭐       | ⭐⭐⭐⭐          | ⭐⭐⭐⭐      |
| **Complex Queries**     | ⭐⭐⭐⭐⭐      | ⭐⭐⭐⭐        | ⭐⭐⭐            | ⭐⭐⭐    | ⭐⭐⭐⭐⭐      | N/A           | Limited   |
| **Full-Text Search**    | ⭐⭐         | ⭐⭐          | ⭐⭐⭐            | ⭐      | ⭐⭐         | ⭐⭐⭐⭐⭐         | N/A       |
| **Scaling**             | Vertical+  | Horizontal+ | Horizontal     | None   | Vertical+  | Horizontal    | Unlimited |
| **Transaction Support** | Full ACID  | InnoDB ACID | Multi-doc ACID | ACID   | Full ACID  | None          | Limited   |

---

## 🚀 Advanced Database Features

### Caching and Transactions
```bash
# Add transaction support
goca repository User --database postgres --transactions

# Add caching layer
goca repository User --database postgres --cache

# Both together
goca repository User --database postgres --transactions --cache
```

### Database-Specific Optimizations

**PostgreSQL:**
- Array operations
- Range types
- Full-text search (basic)
- JSON/JSONB operators

**MongoDB:**
- Aggregation pipeline
- Text search
- Geospatial queries
- Change streams

**Elasticsearch:**
- Analyzers and tokenizers
- Aggregations
- Suggestions
- Highlighting

---

## 📚 Related Documentation

- [goca repository Command](Command-Repository)
- [goca feature Command](Command-Feature)
- [Repository Pattern Guide](Best-Practices#repository-pattern)
- [Complete Tutorial](Complete-Tutorial)
