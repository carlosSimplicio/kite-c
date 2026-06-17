-- name: CreateCommand :one
INSERT INTO command
(name, description, command_query, type_id)
VALUES
(?, ?, ?, ?)
RETURNING *;


-- name: UpdateCommandById :one
UPDATE command SET 
    name = ?,
    description = ?,
    command_query = ?
WHERE id = ?
RETURNING *;

-- name: GetCommandById :one
SELECT id, name, description, command_query, type_id FROM command WHERE id = ?;



