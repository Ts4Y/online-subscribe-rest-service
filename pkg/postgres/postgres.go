package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectToPosgres(ctx context.Context, dsn string) (*pgx.Conn, error) {

	conn, err := pgx.Connect(ctx, dsn)

	if err != nil {
		return nil, fmt.Errorf("pkg/postgres: ConnectToPosgres: %w", err)
	}

	err = conn.Ping(ctx)

	if err != nil {
		return nil, fmt.Errorf("pkg/postgres: conn.Ping: %w", err)
	}

	return conn, nil
}
