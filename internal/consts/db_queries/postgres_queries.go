package db_queries

var (
    CREATE_EXPRESSION_TABLE = "CREATE TABLE IF NOT EXISTS expression (expression_id uuid primary key, status varchar(30), result numeric);"
    INSERT_EXPRESSION = "INSERT INTO expression (expression_id, status, result) VALUES ('%s', '%s', %f);"
    UPDATE_EXPRESSION = "UPDATE expression SET status = '%s', result = %f WHERE expression_id = '%s';"
    GET_EXPRESSIONS = "SELECT * FROM expression;"
    GET_EXPRESSION_BY_ID = "SELECT expression_id, status, result FROM expression WHERE expression_id = '%s';"
)
