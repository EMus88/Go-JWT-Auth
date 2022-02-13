package repository

import (
	"JWT_auth/internal/model"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//this interface implements pgx.Conn, pgx.Pool and pgx.Mock
type DB interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Ping(ctx context.Context) error
}

func NewDB(ctx context.Context) (*pgx.Conn, error) {
	//db connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.username"),
		viper.GetString("db.dbname"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.sslmode"))
	//init connection
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	logrus.Println("DB connection success")
	return conn, nil
}

func AutoMigration(isAllowed bool) error {

	if !isAllowed {
		return nil
	}
	//db connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.username"),
		viper.GetString("db.dbname"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.sslmode"))
	//open connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	//run automigration
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return err
	}
	logrus.Println("Migration success")

	return nil
}
