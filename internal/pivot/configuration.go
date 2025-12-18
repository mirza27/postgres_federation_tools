package pivot

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// GetConfigurationByName mengambil satu entri konfigurasi berdasarkan key.
// Mengembalikan nil jika tidak ditemukan.
func (r *Repo) GetConfigurationByName(name string) (*Configuration, error) {
	var cfg Configuration

	err := r.DB.QueryRow(context.Background(), `
        select config_key, config_value, updated_at from configuration
        where config_key = $1`,
		name).Scan(&cfg.ConfigKey, &cfg.ConfigValue, &cfg.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &cfg, nil
}

// InsertConfigurationByName menambah entri konfigurasi baru. Bila key sudah ada, tidak mengubah nilai.
func (r *Repo) InsertConfigurationByName(config *Configuration) error {
	var updatedAt time.Time

	err := r.DB.QueryRow(context.Background(), `
        insert into configuration (config_key, config_value)
        values ($1, $2) on conflict (config_key) do nothing
        returning updated_at`,
		config.ConfigKey, config.ConfigValue).Scan(&updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			// key sudah ada, treat as success
			return nil
		}
		return err
	}
	config.UpdatedAt = updatedAt
	return nil
}

// UpdateConfigurationByName memperbarui nilai konfigurasi yang ada.
func (r *Repo) UpdateConfigurationByName(config *Configuration) error {
	return r.DB.QueryRow(context.Background(), `
        update configuration
        set config_value = $2, updated_at = now()
        where config_key = $1
        returning updated_at`,
		config.ConfigKey, config.ConfigValue).Scan(&config.UpdatedAt)
}
