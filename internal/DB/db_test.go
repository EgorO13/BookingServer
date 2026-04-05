package db

import (
	"testing"

	"task/internal/Config"

	"github.com/stretchr/testify/assert"
)

func TestNewDB_InvalidConfig(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "nonexistent",
		DBPort:     "5432",
		DBUser:     "user",
		DBPassword: "1234",
		DBName:     "test_db",
	}
	_, err := NewDB(cfg)
	assert.Error(t, err)
}
