package database

import (
	"fmt"
	"net"
	"net/url"

	"notix/config"
	"notix/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DemoData struct {
	BoardID         uint
	OwnerColumnID   uint
	ForbiddenTaskID uint
}

func Connect(cfg config.Config) (*gorm.DB, error) {
	dsn := url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(cfg.DBHost, cfg.DBPort),
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Path:   cfg.DBName,
	}
	query := dsn.Query()
	query.Set("sslmode", cfg.DBSSLMode)
	query.Set("TimeZone", "UTC")
	dsn.RawQuery = query.Encode()

	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to PostgreSQL at %s:%s: %w", cfg.DBHost, cfg.DBPort, err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get PostgreSQL connection pool: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping PostgreSQL at %s:%s: %w", cfg.DBHost, cfg.DBPort, err)
	}
	return db, nil
}

func MigrateAndSeed(db *gorm.DB) (DemoData, error) {
	if err := db.AutoMigrate(
		&models.Usuario{},
		&models.Tablero{},
		&models.Columna{},
		&models.Tarea{},
		&models.TableroCompartido{},
		&models.Comentario{},
		&models.Cuenta{},
	); err != nil {
		return DemoData{}, fmt.Errorf("migrate database tables: %w", err)
	}
	demo, err := SeedDemoData(db)
	if err != nil {
		return DemoData{}, fmt.Errorf("seed demo data: %w", err)
	}
	return demo, nil
}

func SeedDemoData(db *gorm.DB) (DemoData, error) {
	var demo DemoData
	err := db.Transaction(func(tx *gorm.DB) error {
		owner, err := seedUserAndAccount(tx, "Propietaria Demo", "propietaria@notix.test", "Prop-demo-2026", models.RolPropietario)
		if err != nil {
			return err
		}
		collaborator, err := seedUserAndAccount(tx, "Colaborador Demo", "colaborador@notix.test", "Colab-demo-2026", models.RolColaborador)
		if err != nil {
			return err
		}

		board := models.Tablero{Titulo: "Notix - tablero de demostración", CreadorID: owner.ID}
		if err := tx.Where("creador_id = ? AND titulo = ?", owner.ID, board.Titulo).FirstOrCreate(&board).Error; err != nil {
			return fmt.Errorf("create demo board: %w", err)
		}
		demo.BoardID = board.ID

		var columns []models.Columna
		for i, name := range []string{"Por hacer", "En proceso", "Completado"} {
			column := models.Columna{Nombre: name, Orden: i + 1, TableroID: board.ID}
			if err := tx.Where("tablero_id = ? AND orden = ?", board.ID, i+1).FirstOrCreate(&column).Error; err != nil {
				return fmt.Errorf("create demo column %q: %w", name, err)
			}
			columns = append(columns, column)
		}
		demo.OwnerColumnID = columns[0].ID

		share := models.TableroCompartido{
			TableroID:        board.ID,
			UsuarioID:        collaborator.ID,
			RolEnTablero:     models.RolColaborador,
			EstadoInvitacion: models.InvitacionActiva,
		}
		if err := tx.Where("tablero_id = ? AND usuario_id = ?", board.ID, collaborator.ID).Assign(
			models.TableroCompartido{RolEnTablero: models.RolColaborador, EstadoInvitacion: models.InvitacionActiva},
		).FirstOrCreate(&share).Error; err != nil {
			return fmt.Errorf("share demo board with collaborator: %w", err)
		}

		forbiddenTask := models.Tarea{
			Titulo:    "DEMO: transición prohibida",
			Contenido: "Debe permanecer en pendiente durante toda la demostración.",
			Estado:    models.EstadoPendiente,
			ColumnaID: columns[0].ID,
			CreadorID: owner.ID,
		}
		if err := tx.Where("columna_id = ? AND titulo = ?", columns[0].ID, forbiddenTask.Titulo).FirstOrCreate(&forbiddenTask).Error; err != nil {
			return fmt.Errorf("create forbidden-transition demo task: %w", err)
		}
		demo.ForbiddenTaskID = forbiddenTask.ID
		return nil
	})
	return demo, err
}

func seedUserAndAccount(tx *gorm.DB, name, email, password, role string) (models.Usuario, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, fmt.Errorf("hash demo password for %s: %w", email, err)
	}
	user := models.Usuario{Nombre: name, Correo: email, Clave: string(hash)}
	if err := tx.Where("correo = ?", email).FirstOrCreate(&user).Error; err != nil {
		return models.Usuario{}, fmt.Errorf("create demo user %s: %w", email, err)
	}

	account := models.Cuenta{
		UsuarioID: user.ID,
		Correo:    email,
		ClaveHash: string(hash),
		Rol:       role,
	}
	if err := tx.Where("correo = ?", email).FirstOrCreate(&account).Error; err != nil {
		return models.Usuario{}, fmt.Errorf("create demo account %s: %w", email, err)
	}
	return user, nil
}
