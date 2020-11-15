## Assumptions:
1. No symlinks in patches directory
2. Each cli execution only handles patches to 1 **PostgreSQL** database
3. the SQL statement in patch json file is not over 10s of million lines(bad idea for a db patch anyways...) we slurp **the whole json file into memory**!

## Build

You need to install Golang to build this program. The code is tested in `go 1.13`

```bash
go build -o dbsync cmd/dbsync/main.go
```

## Usage
```bash
 ./dbsync -path_dir migration-0.0.1 -db_conn conn_info.yaml
 ```

## Example
use the migrations in examples folder. The migration assumes you have a "employee" table in the db you connect to

```sql
CREATE TABLE EMPLOYEE(
   ID INT PRIMARY KEY     NOT NULL,
   NAME           TEXT    NOT NULL
);
```