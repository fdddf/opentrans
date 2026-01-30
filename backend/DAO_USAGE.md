# GORM Gen DAO Usage Guide

## Overview
This project now uses GORM Gen to generate type-safe DAO (Data Access Object) code for all database models. This provides better performance, type safety, and reduces SQL injection risks.

## Generated Models
GORM Gen will generate type-safe query code for the following models:
- User
- Project
- Translation
- UserActivity
- ProviderConfig
- App
- AppLocalization
- Subscription
- AppUser
- AppProviderConfig
- TranslationQueue

## How to Generate Code

### 1. Run the Generation Command
```bash
# From the backend directory
cd backend
go run main.go generate-dao
```

This will generate the DAO code in `./internal/dao/query/` directory.

### 2. Using the Generated Code

After running the generation command, you can use the generated DAOs in your services:

```go
import (
    "github.com/fdddf/xcstrings-translator/internal/dao/query"
    
    // Import your models
    "github.com/fdddf/xcstrings-translator/internal/database"
)

// Initialize the generated query
var db *gorm.DB // your database connection
q := query.Use(db)

// Example: Query all users
users, err := q.User.Find()
if err != nil {
    // handle error
}

// Example: Find a user by ID
user, err := q.User.Where(q.User.ID.Eq(1)).First()
if err != nil {
    // handle error
}

// Example: Create a new user
newUser := &database.User{
    Username: "example",
    Email:    "example@example.com",
    // ... other fields
}
err = q.User.Create(newUser)
if err != nil {
    // handle error
}

// Example: Update user
_, err = q.User.Where(q.User.ID.Eq(1)).Update(q.User.Username, "newname")
if err != nil {
    // handle error
}

// Example: Delete user
_, err = q.User.Where(q.User.ID.Eq(1)).Delete()
if err != nil {
    // handle error
}
```

## Benefits of Using GORM Gen

1. **Type Safety**: All queries are type-safe, reducing runtime errors
2. **Performance**: Optimized SQL queries with better prepared statement usage
3. **Auto-completion**: IDE auto-completion for all database operations
4. **SQL Injection Protection**: All queries are parameterized by default
5. **Maintainability**: Changes to models automatically update generated queries

## Migration Files Structure

The migration files have been reorganized to follow golang-migrate standards:

```
backend/migrations/
├── 000001_init_schema.up.sql    # Schema creation
└── 000001_init_schema.down.sql  # Schema deletion
```

This structure allows for better management of database schema changes and provides clearer migration steps.

## Running Migrations

```bash
# Apply all pending migrations
cd backend
go run main.go migrate

# Rollback the last migration (for development)
go run main.go migrate rollback
```

## Updated Services Using Generated DAO

Some services have been updated to use the generated DAO code for better performance and type safety. For example:

### AppService Examples

#### CreateApp with generated DAO:
```go
// Check if bundle ID already exists
existingApp, err := s.Query.App.Where(s.Query.App.BundleID.Eq(bundleID)).First()
if err == nil && existingApp != nil {
    return nil, errors.New("bundle ID already exists")
} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, fmt.Errorf("failed to check bundle ID: %v", err)
}

// Check user's subscription limit
user, err := s.Query.User.Where(s.Query.User.ID.Eq(userID)).First()
if err != nil {
    return nil, fmt.Errorf("failed to get user: %v", err)
}

// Check if user has reached the app limit based on subscription
if user.CurrentAppCount >= user.MaxApps {
    return nil, errors.New("app limit reached for your subscription")
}

// Create the app
err = s.Query.App.Create(app)
if err != nil {
    return nil, fmt.Errorf("failed to create app: %v", err)
}

// Update user's app count
_, err = s.Query.User.Where(s.Query.User.ID.Eq(userID)).Update(s.Query.User.CurrentAppCount, user.CurrentAppCount+1)
```

#### GetAppsByUser with generated DAO:
```go
apps, err := s.Query.App.Where(s.Query.App.UserID.Eq(userID)).Order(s.Query.App.CreatedAt.Desc()).Find()
if err != nil {
    return nil, fmt.Errorf("failed to retrieve apps: %v", err)
}

// Convert slice of pointers to slice of values
appSlice := make([]database.App, len(apps))
for i, app := range apps {
    appSlice[i] = *app
}

return appSlice, nil
```

## Benefits of Using GORM Gen

1. **Type Safety**: All queries are type-safe, reducing runtime errors
2. **Performance**: Optimized SQL queries with better prepared statement usage
3. **Auto-completion**: IDE auto-completion for all database operations
4. **SQL Injection Protection**: All queries are parameterized by default
5. **Maintainability**: Changes to models automatically update generated queries
6. **Better Error Handling**: More detailed error information
7. **Reduced Boilerplate**: Less code needed for common database operations
8. **Prepared Statements**: Better security and performance
9. **Join Support**: Better support for complex queries with joins
10. **Transaction Support**: Better transaction management capabilities

## Service Updates

The following services have been updated to use the generated DAO code:
- AppService (CreateApp, GetAppsByUser, GetUsersForApp methods)
- Additional services can be updated following the same pattern

To update additional services, follow these patterns:
1. Add the Query field to the service struct
2. Initialize the query in the SetService function: `query.Use(db.DB)`
3. Replace raw GORM queries with generated query methods
4. Handle the differences in return types (generated queries return pointers to structs)
