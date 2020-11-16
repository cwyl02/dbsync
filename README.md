# dbsync
dbsync is a CLI tool that manages PostgreSQL db table schema migration written in golang.

## Assumptions:
1. No symlinks in patches directory
2. Each cli execution only handles patches to 1 **PostgreSQL** database
3. the SQL statement in patch json file is not over 10s of million lines(bad idea for a db patch anyways...) we slurp **the whole json file into memory**!
4. User answer of [Y/n] doesn't contain space(' ') character

## Build

You need to install Golang to build this program. The code is tested in `go 1.13`

```bash
go build -o dbsync cmd/dbsync/main.go
```

## Usage
```bash
 ./dbsync -path_dir migration-0.0.1 -db_conn conn_info.yaml
 ```

## Test
use the migrations in examples folder. The example assumes:
- you have a database "db0" in your PostgreSQL instance
- you have a "employee" table in the db 
    - ```sql
        CREATE TABLE EMPLOYEE(
        ID INT PRIMARY KEY     NOT NULL,
        NAME           TEXT    NOT NULL
        );
        ```
- you have a "example" user, and it has all necessary permissions to operate on employee

```bash
# test patch migration-0.0.1
./dbsync -db_conn examples/conn_info.yaml -path_dir examples/migration-0.0.1
# test a (revert) patch revert-0.0.1
./dbsync -db_conn examples/conn_info.yaml -path_dir examples/revert-0.0.1
# test another patch that depends on all patches above
./dbsync -db_conn examples/conn_info.yaml -path_dir examples/multiline-0.0.2
```