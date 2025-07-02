# pgdb

Simple helper package to insert and upsert data into PostgreSQL using [Bun ORM](https://bun.uptrace.dev/) in Go.

This is a small utility package for personal and internal use — but feel free to use or adapt it.

## Features

- Connect to PostgreSQL via Bun + pgx
- Generic `ToDb` function to insert or upsert slices of data
- Automatically builds `ON CONFLICT DO UPDATE` clause based on primary key and columns
- No hardcoding of column names — uses `information_schema` dynamically
- set DEBUG_SQL=true to log SQL queries for debugging

## Installation

```bash
go get github.com/imbence/pgdb
````