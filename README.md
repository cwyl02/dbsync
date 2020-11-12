## Assumptions:
1. No symlinks in patches directory
2. Each cli execution only handles patches to 1 *Postgres* database
3. the SQL statement in patch json file is not over 10s of million lines(bad idea for a db patch anyways...) we slurp the whole json file into memory!


## Build
```go build -o dbsync cmd/dbsync/main.go```