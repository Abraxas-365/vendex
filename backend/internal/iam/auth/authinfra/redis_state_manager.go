package authinfra

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/iam/auth"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisStateManager implementación en Redis del StateManager
type RedisStateManager struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisStateManager crea un nuevo state manager con Redis
func NewRedisStateManager(client *redis.Client, ttl time.Duration) auth.StateManager {
	return &RedisStateManager{
		client: client,
		ttl:    ttl,
	}
}

// GenerateState genera un nuevo estado OAuth
func (sm *RedisStateManager) GenerateState() string {
	return uuid.NewString()
}

// StoreState almacena un estado con sus datos asociados
func (sm *RedisStateManager) StoreState(ctx context.Context, state string, data map[string]any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errx.Wrap(err, "failed to marshal state data", errx.TypeInternal)
	}

	key := fmt.Sprintf("oauth_state:%s", state)
	err = sm.client.Set(ctx, key, jsonData, sm.ttl).Err()
	if err != nil {
		return errx.Wrap(err, "failed to store state in Redis", errx.TypeExternal)
	}

	return nil
}

// ValidateState valida si un estado es válido
func (sm *RedisStateManager) ValidateState(state string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("oauth_state:%s", state)

	exists, err := sm.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return exists == 1
}

// GetStateData obtiene los datos asociados a un estado
func (sm *RedisStateManager) GetStateData(ctx context.Context, state string) (map[string]any, error) {
	key := fmt.Sprintf("oauth_state:%s", state)

	// Obtener y eliminar el estado (one-time use)
	jsonData, err := sm.client.GetDel(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, auth.ErrInvalidState()
		}
		return nil, errx.Wrap(err, "failed to get state from Redis", errx.TypeExternal)
	}

	var data map[string]any
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, errx.Wrap(err, "failed to unmarshal state data", errx.TypeInternal)
	}

	return data, nil
}
