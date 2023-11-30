package infrastructure

import (
	"context"
	"database/sql"
	config "erp/config"
	"erp/internal/models"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseTransaction interface {
	WithTransaction(txFunc func(tx *Database) error) (err error)
}

type Database struct {
	RDBMS   *gorm.DB
	Mongo   *mongo.Database
	Context context.Context // use for mongo transaction
	logger  *zap.Logger
}

func NewDatabase(config *config.Config, logger *zap.Logger) Database {
	var err error
	var sqlDB *sql.DB

	logger.Info("Connecting to database...")
	rdbms, nosql, err := getDatabaseInstance(config)
	if err != nil {
		for i := 0; i < 5; i++ {
			rdbms, nosql, err = getDatabaseInstance(config)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		logger.Fatal("Database connection error", zap.Error(err))
	} else {
		logger.Info("Database connected")
	}

	db := Database{RDBMS: rdbms, Mongo: nosql, logger: logger}

	if config.Database.Driver == "mysql" || config.Database.Driver == "postgres" {
		db.RegisterTables()

		if err != nil {
			logger.Fatal("Database connection error", zap.Error(err))
		}
		sqlDB, err = db.RDBMS.DB()
		if err != nil {
			logger.Fatal("sqlDB connection error", zap.Error(err))
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
	}
	return db
}

func getDatabaseInstance(config *config.Config) (rdbms *gorm.DB, nosql *mongo.Database, err error) {
	switch config.Database.Driver {
	//case "mysql":
	//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	//		config.Database.Username,
	//		config.Database.Password,
	//		config.Database.Host,
	//		config.Database.Port,
	//		config.Database.Name,
	//	)
	//	rdbms, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//	if err != nil {
	//		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	//	}
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			config.Database.Postgres.Host, config.Database.Postgres.Username, config.Database.Postgres.Password, config.Database.Postgres.Name,
			config.Database.Postgres.Port, config.Database.Postgres.SSLMode, config.Database.Postgres.TimeZone)
		rdbms, err = gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})

		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect database: %w", err)
		}

		rdbms.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	case "mongo":
		ctx := context.TODO()
		uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/?replicaSet=rs0", config.Database.Mongo.Username, config.Database.Mongo.Password, config.Database.Mongo.Host, config.Database.Mongo.Port)
		clientOptions := options.Client().SetDirect(true).ApplyURI(uri)
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		nosql = client.Database(config.Database.Mongo.Name)
	}

	return rdbms, nosql, nil
}

func (d Database) RegisterTables() {
	err := d.RDBMS.AutoMigrate(
		models.User{},
		models.Token{},
	)

	if err != nil {
		d.logger.Fatal("Database migration error", zap.Error(err))
		os.Exit(0)
	}
	d.logger.Info("Database migration success")
}
