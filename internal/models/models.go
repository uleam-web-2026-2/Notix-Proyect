package models

import "time"

const (
	EstadoPendiente  = "pendiente"
	EstadoEnProceso  = "en_proceso"
	EstadoCompletado = "completado"

	RolPropietario      = "propietario"
	RolColaborador      = "colaborador"
	InvitacionPendiente = "pendiente"
	InvitacionActiva    = "activo"
	InvitacionRevocada  = "revocado"
)

type Usuario struct {
	ID     uint      `gorm:"primaryKey" json:"id"`
	Nombre string    `gorm:"not null" json:"nombre"`
	Correo string    `gorm:"uniqueIndex;not null" json:"correo"`
	Clave  string    `gorm:"not null" json:"-"`
	Creado time.Time `gorm:"autoCreateTime" json:"creado"`
}

func (Usuario) TableName() string { return "usuarios" }

type Tablero struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Titulo      string    `gorm:"not null" json:"titulo"`
	Descripcion string    `json:"descripcion"`
	CreadorID   uint      `gorm:"not null;index" json:"creador_id"`
	Creador     Usuario   `gorm:"foreignKey:CreadorID" json:"creador,omitempty"`
	Creado      time.Time `gorm:"autoCreateTime" json:"creado"`
}

func (Tablero) TableName() string { return "tableros" }

type Columna struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Nombre    string  `gorm:"not null" json:"nombre"`
	Orden     int     `gorm:"not null" json:"orden"`
	TableroID uint    `gorm:"not null;index" json:"tablero_id"`
	Tablero   Tablero `gorm:"foreignKey:TableroID" json:"-"`
}

func (Columna) TableName() string { return "columnas" }

type Tarea struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Titulo    string    `gorm:"not null" json:"titulo"`
	Contenido string    `json:"contenido"`
	Estado    string    `gorm:"type:varchar(20);not null;default:'pendiente';check:chk_tarea_estado,estado IN ('pendiente','en_proceso','completado')" json:"estado"`
	ColumnaID uint      `gorm:"not null;index" json:"columna_id"`
	Columna   Columna   `gorm:"foreignKey:ColumnaID" json:"columna,omitempty"`
	CreadorID uint      `gorm:"not null;index" json:"creador_id"`
	Creador   Usuario   `gorm:"foreignKey:CreadorID" json:"creador,omitempty"`
	Creado    time.Time `gorm:"autoCreateTime" json:"creado"`
}

func (Tarea) TableName() string { return "tareas" }

type TableroCompartido struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TableroID        uint      `gorm:"not null;uniqueIndex:idx_tablero_usuario" json:"tablero_id"`
	Tablero          Tablero   `gorm:"foreignKey:TableroID" json:"tablero,omitempty"`
	UsuarioID        uint      `gorm:"not null;uniqueIndex:idx_tablero_usuario" json:"usuario_id"`
	Usuario          Usuario   `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	RolEnTablero     string    `gorm:"type:varchar(20);not null;default:'colaborador';check:chk_rol_tablero,rol_en_tablero IN ('propietario','colaborador')" json:"rol_en_tablero"`
	EstadoInvitacion string    `gorm:"type:varchar(20);not null;default:'pendiente';check:chk_estado_invitacion,estado_invitacion IN ('pendiente','activo','revocado')" json:"estado_invitacion"`
	Agregado         time.Time `gorm:"autoCreateTime" json:"agregado"`
}

func (TableroCompartido) TableName() string { return "tableros_compartidos" }

type Comentario struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	Texto   string    `gorm:"not null" json:"texto"`
	TareaID uint      `gorm:"not null;index" json:"tarea_id"`
	Tarea   Tarea     `gorm:"foreignKey:TareaID" json:"-"`
	AutorID uint      `gorm:"not null;index" json:"autor_id"`
	Autor   Usuario   `gorm:"foreignKey:AutorID" json:"autor,omitempty"`
	Creado  time.Time `gorm:"autoCreateTime" json:"creado"`
}

func (Comentario) TableName() string { return "comentarios" }

type Cuenta struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	UsuarioID uint    `gorm:"not null;uniqueIndex" json:"usuario_id"`
	Usuario   Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Correo    string  `gorm:"uniqueIndex;not null" json:"correo"`
	ClaveHash string  `gorm:"not null" json:"-"`
	Rol       string  `gorm:"type:varchar(20);not null;check:chk_cuenta_rol,rol IN ('propietario','colaborador')" json:"rol"`
}

func (Cuenta) TableName() string { return "cuentas" }
