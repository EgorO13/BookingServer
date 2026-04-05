package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, "5432", cfg.DBPort)
	assert.Equal(t, "booking", cfg.DBUser)
	assert.Equal(t, "password", cfg.DBPassword)
	assert.Equal(t, "bookingdb", cfg.DBName)
	assert.Equal(t, "default-secret", cfg.JWTSecret)
}

func TestLoad_FromEnv(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_HOST", "postgres")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "secret")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("JWT_SECRET", "mysecret")

	cfg := Load()
	assert.Equal(t, "postgres", cfg.DBHost)
	assert.Equal(t, "5433", cfg.DBPort)
	assert.Equal(t, "testuser", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
	assert.Equal(t, "testdb", cfg.DBName)
	assert.Equal(t, "mysecret", cfg.JWTSecret)
}

func TestGetEnv(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, "default", getEnv("NONEXISTENT", "default"))
	os.Setenv("EXISTENT", "value")
	assert.Equal(t, "value", getEnv("EXISTENT", "default"))
}
