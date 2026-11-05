CREATE TABLE tareas (
    id         bigserial PRIMARY KEY,
    titulo     text NOT NULL,
    contenido  text,
    estado     varchar(20) NOT NULL DEFAULT 'pendiente',
    columna_id bigint NOT NULL,
    creador_id bigint NOT NULL,
    creado     timestamptz,
    CHECK (estado IN ('pendiente', 'en_proceso', 'completado'))
);
