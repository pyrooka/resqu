# resqu
resqu (RESt from QUeries) helps you to serve database query results as JSON via HTTP.

NOTE: the project is under development, so some breaking changes may happen

## Usage
### Docker
`docker run -v $(pwd)/config.yaml:/resqu/config.yaml pyrooka/resqu`

If you want to use docker compose check the examples directory.

### Helm
`helm upgrade -i resqu ./resqu-helm`

### Environment variables
- `SERVER_PORT`

### DB configs
E.g. `config.yaml`:
```yaml
sqlite:
  connection:
    path: /db/employees.sqlite3
  endpoints:
    - url: /employees
      query: SELECT * FROM emp LIMIT 10
    - url: /employees/{empNo}
      query: SELECT * FROM emp WHERE empno = {empNo}
```

#### SQLite3
```yaml
sqlite:
  connection:
    path: /db/employees.sqlite3 # path to the DB file
  endpoints: []
```

#### PostgreSQL
```yaml
postgresql:
  connection:
    connectionURL: "postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]"
  endpoints: []
```

#### BigQuery
```yaml
bigquery:
  connection:
    serviceAccPath: /google/sa.json # path to the service account json
    projectId: my-awesome-project # the project ID for the queries
  endpoints: []
```

### API responses
| Status code | Body            |
|:-----------:|-----------------|
|     200     | {"data": [...]} |
|     400     | error message   |
|     500     | error message   |