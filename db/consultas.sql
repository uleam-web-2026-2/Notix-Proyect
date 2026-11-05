-- name: BuscarTarea :one
SELECT id, titulo, contenido, estado, columna_id, creador_id, creado
FROM tareas
WHERE id = $1;

-- name: ActualizarEstadoTarea :execrows
UPDATE tareas
SET estado = $2
WHERE id = $1 AND estado = $3;
