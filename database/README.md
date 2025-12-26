# Database Executor

The `database` package provides a flexible wrapper around `database/sql` for query execution. It encapsulates the standard DB operations behind a `queryExecutor` interface, allowing for cleaner code and easier mocking.

## Installation

```go
import "github.com/aidapedia/gdk/database"
```

## Usage

### Initialization

Create a new executor instance by passing a standard `*sql.DB` connection.

```go
import (
	"database/sql"
	"github.com/aidapedia/gdk/database"
)

func main() {
	// Initialize your sql.DB connection
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname")
	if err != nil {
		panic(err)
	}

	// Create the executor
	exec := database.NewExecutor(db)
}
```

### Executing Queries

Use the `DB()` method to access the `queryExecutor` interface. This enables you to perform standard SQL operations like `Exec`, `Query`, and `QueryRow` (and their Context variants).

```go
type Repository struct{
    db *sql.DB
}
func (r *Repository) CreateUser(name, password string, opts ...database.Option) error {
    exec := database.NewExecutor(r.db)
    for _, opt := range opts {
        opt.Apply(&exec)
    }

    _, err := exec.DB().Exec("INSERT INTO users (name, password) VALUES ($1, $2)", name, password)
    if err != nil {
        return err
    }
    return nil
}
```

### Transaction Support

The package includes a `WithTransaction` option helper, which is useful when you need to swap the underlying executor with a transaction instance (`*sql.Tx`).

```go
tx, err := db.Begin()
if err != nil {
    return err
}

err = userRepository.CreateUser("John", "123456", database.WithTransaction(tx))
if err != nil {
    tx.Rollback()
    return 
}

tx.Commit()
return
```
