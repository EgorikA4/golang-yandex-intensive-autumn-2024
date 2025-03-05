package db_queries

var (
    CREATE_NODE = "CREATE (n:Node {id: '%s', operation: '%s', arg1: '%s', arg2: '%s', expression_id: '%s'});"
    CREATE_RIGHT_REL = "MATCH (child:Node {id: '%s'}), (root:Node {id: '%s'}) CREATE (child)-[r:RIGHT]->(root);"
    CREATE_LEFT_REL = "MATCH (child:Node {id: '%s'}), (root:Node {id: '%s'}) CREATE (child)-[r:LEFT]->(root);"
    UPDATE_NODE = `MATCH (child:Node {id: '%s'})-[r]->(parent:Node)
    SET parent.arg1 = CASE WHEN type(r) = 'LEFT' THEN '%f' ELSE parent.arg1 END,
        parent.arg2 = CASE WHEN type(r) = 'RIGHT' THEN '%f' ELSE parent.arg2 END
    DELETE r, child
    RETURN parent;`
    IS_ROOT = `MATCH (n:Node {id: '%s'})
    OPTIONAL MATCH (n)-[left:LEFT]->()
    OPTIONAL MATCH (n)-[right:RIGHT]->()
    RETURN left IS NULL AND right IS NULL AS is_root, n.expression_id as expression_id;`
    DELETE_NODE_BY_EXPRESSION_ID = `MATCH (n:Node) WHERE n.expression_id = '%s' DETACH DELETE n;`
    GET_NODE_BY_ID = "MATCH (n:Node {id: '%s'}) RETURN n;"
)
