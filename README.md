# resqu
resqu (RESt from QUeries) helps you to serve database query results as JSON via HTTP.
## Config
See `config.yaml` for examples.
### SQLite3
`path`: path to the DB file  
### PostgreSQL
`connectionURL`: postgresql://[user[:password]@][netloc][:port][,...][/dbname][?param1=value1&...]  
### BigQuery
`serviceAccPath`: path to the service account json  
`projectId`: the project ID for the queries  