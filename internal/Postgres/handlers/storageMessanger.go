package handlers

import (
	"VK_posts/internal/models"
	"VK_posts/pkg/Postgres"
	"context"
	"fmt"
)

type MessengerHandler struct {
	Storage *Postgres.Storage
}

func (m *MessengerHandler) MessageSave(message models.Message) error {
	stmt, err := m.Storage.Db.Prepare(`INSERT INTO (from,to,content,created_at) VALUES ($1,$2,$3,$4)`)
	if err != nil {
		return fmt.Errorf("error while preparing statement: %w", err)
	}
	ctx := context.Background()
	result, err := stmt.ExecContext(ctx, message.From, message.To, message.Content, message.Timestamp)
	if err != nil {
		return fmt.Errorf("error while executing statement: %w", err)
	}
	id, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error while getting rows affected: %w", err)
	}
	if id == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (m *MessengerHandler) GetUserServer(userid string) (string, error) {
	stmt, err := m.Storage.Db.Prepare("SELECT server FROM server_log WHERE userid = $1 ")
	if err != nil {
		return "", fmt.Errorf("error while preparing statement: %w", err)
	}
	ctx := context.Background()
	var server string
	err = stmt.QueryRowContext(ctx, userid).Scan(&server)
	if err != nil {
		return "", fmt.Errorf("error while getting user server: %w", err)
	}
	return server, nil
}
