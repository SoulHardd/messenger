package dto

import (
	"D/Go/messenger/internal/chat/domain"
	"encoding/base64"
	"encoding/json"
)

func EncodeCursor(c *domain.Cursor) string {
	raw, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(raw)
}

func DecodeCursor(s string) (*domain.Cursor, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	var c domain.Cursor
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
