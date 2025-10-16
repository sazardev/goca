# Database Support

Goca supports **8 production-ready database systems** with full repository pattern implementations. Choose the database that best fits your use case.

## Supported Databases

### SQL Databases

#### PostgreSQL
**Type:** SQL (ACID-compliant)  
**Driver:** GORM + postgres driver  
**Best For:** Traditional applications, transactional systems, most common choice

```bash
goca feature User --fields "name:string,email:string" --database postgres
```

**Features:**
- Full transaction support
- Complex joins and aggregations
- Advanced query capabilities
- Excellent performance for relational data

**Generated File:** `internal/repository/postgres_user_repository.go`

---

#### PostgreSQL JSON
**Type:** SQL + Document (semi-structured)  
**Driver:** GORM + gorm/datatypes (JSONB)  
**Best For:** Semi-structured data, configs, metadata, hierarchical objects

```bash
goca feature Config --fields "name:string,settings:string" --database postgres-json
```

**Special Methods:**
- `FindByJSONField(jsonField, value string)` - Query nested JSON with `@>` operator
- Full JSONB operator support

**Use Case Example:**
```go
// Store nested configuration
configs, _ := repo.FindByJSONField("settings.theme", "dark")
```

**Generated File:** `internal/repository/postgres_json_config_repository.go`

---

#### MySQL
**Type:** SQL (ACID-compliant)  
**Driver:** GORM + mysql driver  
**Best For:** Web applications, compatibility with existing infrastructure

```bash
goca feature Product --fields "name:string,price:float64" --database mysql
```

**Features:**
- InnoDB transactions
- Good scalability for web apps
- Strong compatibility
- Replication support

**Generated File:** `internal/repository/mysql_product_repository.go`

---

#### SQL Server
**Type:** SQL (Enterprise)  
**Driver:** GORM + mssql driver  
**Best For:** Enterprise applications, Microsoft ecosystem, legacy systems

```bash
goca feature Employee --fields "name:string,salary:float64" --database sqlserver
```

**Features:**
- T-SQL compatibility
- Enterprise security
- Advanced error handling
- Context-aware operations
- Enterprise compatibility

**Generated File:** `internal/repository/sqlserver_employee_repository.go`

---

#### SQLite
**Type:** SQL (Embedded)  
**Driver:** database/sql + SQLite driver  
**Best For:** Development, testing, embedded applications, single-file databases

```bash
goca feature Setting --fields "key:string,value:string" --database sqlite
```

**Features:**
- Single file storage (`.db`)
- No server required
- ACID compliance
- JSON marshaling for flexibility
- Perfect for prototyping

**Generated File:** `internal/repository/sqlite_setting_repository.go`

---

### NoSQL Databases

#### MongoDB
**Type:** Document NoSQL  
**Driver:** MongoDB official driver  
**Best For:** Document-oriented applications, flexible schemas, rapid iteration

```bash
goca feature Article --fields "title:string,content:string,tags:string" --database mongodb
```

**Features:**
- Flexible schema
- Rich query language
- Horizontal scalability
- Document transactions

**Generated File:** `internal/repository/mongodb_article_repository.go`

---

#### DynamoDB
**Type:** Key-Value NoSQL (Serverless)  
**Driver:** AWS SDK v2 + attributevalue  
**Best For:** Serverless AWS applications, auto-scaling, cloud-native architecture

```bash
goca feature Order --fields "orderID:string,total:float64" --database dynamodb
```

**Features:**
- AWS-managed serverless
- Auto-scaling
- Scan operations for queries
- Context-aware async operations
- Pay-per-request pricing

**Generated File:** `internal/repository/dynamodb_order_repository.go`

---

### Search Databases

#### Elasticsearch
**Type:** Full-text Search & Analytics  
**Driver:** go-elasticsearch v8 client  
**Best For:** Full-text search, analytics, logging systems, product search

```bash
goca feature Article --fields "title:string,content:string" --database elasticsearch
```

**Special Methods:**
- `FullTextSearch(query string)` - Multi-field full-text search
- Lucene query DSL support
- Result scoring and aggregations

**Use Case Example:**
```go
// Full-text search across multiple fields
results, _ := repo.FullTextSearch("golang elasticsearch tutorial")
```

**Generated File:** `internal/repository/elasticsearch_article_repository.go`

---

## Database Comparison

| Feature             | PostgreSQL  | MySQL       | MongoDB     | SQLite     | SQL Server | Elasticsearch | DynamoDB    |
| ------------------- | ----------- | ----------- | ----------- | ---------- | ---------- | ------------- | ----------- |
| **Type**            | SQL         | SQL         | Document    | SQL        | SQL        | Search        | Key-Value   |
| **ACID**            | ✅           | ✅           | ✅           | ✅          | ✅          | ❌             | Limited     |
| **Scalability**     | Vertical    | Horizontal  | Horizontal  | None       | Vertical   | Horizontal    | Unlimited   |
| **JSON Support**    | JSONB       | JSON        | Native      | Marshaling | Native     | Native        | Native      |
| **Transactions**    | ✅ Full      | ✅ InnoDB    | ✅ Multi-doc | ✅          | ✅          | ❌             | Limited     |
| **Server Required** | ✅           | ✅           | ✅           | ❌          | ✅          | ✅             | N/A (Cloud) |
| **Cost**            | Self-hosted | Self-hosted | Self-hosted | Free       | Enterprise | Self-hosted   | Pay-per-use |

---

## Quick Selection Guide

**Choose PostgreSQL if:**
- You need strong ACID guarantees
- Your data is highly relational
- You want the most feature-rich SQL database
- You're building traditional business applications

**Choose PostgreSQL JSON if:**
- You have semi-structured data
- You need flexible schema evolution
- You want to store hierarchical configs
- You want SQL with document flexibility

**Choose MySQL if:**
- You're integrating with existing MySQL infrastructure
- You need good web application performance
- You want something lightweight
- You're familiar with MySQL ecosystem

**Choose MongoDB if:**
- Your schema changes frequently
- You have flexible data structures
- You need horizontal scaling
- You're building rapid prototypes

**Choose SQLite if:**
- You're developing locally
- You need zero configuration database
- You're building embedded applications
- You want a single-file database

**Choose SQL Server if:**
- You're in Microsoft enterprise environment
- You need T-SQL specific features
- You're migrating from legacy SQL Server
- You need enterprise support

**Choose Elasticsearch if:**
- You need full-text search capabilities
- You're building search-heavy applications
- You need analytics capabilities
- You're centralizing application logs

**Choose DynamoDB if:**
- You're building serverless AWS applications
- You need unlimited auto-scaling
- You want managed cloud database
- You're comfortable with eventual consistency

---

## Migration Between Databases

One of the key benefits of Goca's repository pattern is the ability to switch databases without changing your business logic:

```bash
# Start with PostgreSQL
goca feature User --fields "name:string,email:string" --database postgres

# Later, switch to MongoDB (same feature structure)
goca repository User --database mongodb --implementation
```

Your use cases remain unchanged - only the repository implementation differs!

---

## Performance Considerations

### By Use Case

**High-Concurrency Writes:** PostgreSQL, MongoDB, DynamoDB

**Complex Queries:** PostgreSQL, PostgreSQL JSON, SQL Server

**Document Storage:** MongoDB, PostgreSQL JSON

**Full-Text Search:** Elasticsearch (specialized)

**Real-Time Analytics:** Elasticsearch

**Development/Testing:** SQLite

---

## Next Steps

- [Repository Pattern Guide](/guide/best-practices#repository-pattern)
- [Feature Generation](/commands/feature)
- [Repository Command](/commands/repository)
- [Complete Tutorial](/tutorials/complete-tutorial)
