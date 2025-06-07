# pg_db

Simple helper package to insert and upsert data into PostgreSQL using [Bun ORM](https://bun.uptrace.dev/) in Go.

This is a small utility package for personal and internal use — but feel free to use or adapt it.

## Features

- Connect to PostgreSQL via Bun + pgx
- Generic `ToDb` function to insert or upsert slices of data
- Automatically builds `ON CONFLICT DO UPDATE` clause based on primary key and columns
- No hardcoding of column names — uses `information_schema` dynamically

## Installation

```bash
go get github.com/iambence/pg_db
````

## Usage

```go
package main
import (
    "fmt"
    "github.com/imbence/pg_db"
)

func main() {
    err := pg_db.ConnectToDb("postgres://user:pass@host:5432/dbname?sslmode=disable")
    if err != nil {
        panic(err)
    }

    type MyTable struct {
        ID    int    `bun:"id,pk"`
        Name  string `bun:"name"`
        Value int    `bun:"value"`
    }

    data := []MyTable{
        {ID: 1, Name: "foo", Value: 100},
        {ID: 2, Name: "bar", Value: 200},
    }

    rows, err := pg_db.ToDb(data, "my_table", "public")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Inserted/Updated %d rows\n", rows)
}
```