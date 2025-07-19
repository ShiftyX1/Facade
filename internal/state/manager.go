package state

import (
	"encoding/json"
	"log/slog"
	"sync"
)


type Manager struct {
	data   map[string]interface{}
	mutex  sync.RWMutex
	logger bool
}


func New(enableLogging bool) *Manager {
	return &Manager{
		data:   make(map[string]interface{}),
		logger: enableLogging,
	}
}


func (m *Manager) Set(key string, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = value

	if m.logger {
		slog.Debug("State updated", "key", key, "value", value)
	}
}


func (m *Manager) Get(key string) (interface{}, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	value, exists := m.data[key]
	return value, exists
}


func (m *Manager) GetString(key string) (string, bool) {
	value, exists := m.Get(key)
	if !exists {
		return "", false
	}

	str, ok := value.(string)
	return str, ok
}


func (m *Manager) GetInt(key string) (int, bool) {
	value, exists := m.Get(key)
	if !exists {
		return 0, false
	}

	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	default:
		return 0, false
	}
}


func (m *Manager) GetJSON(key string) (string, bool) {
	value, exists := m.Get(key)
	if !exists {
		return "", false
	}

	data, err := json.Marshal(value)
	if err != nil {
		return "", false
	}

	return string(data), true
}


func (m *Manager) SetFromJSON(key, jsonStr string) error {
	var value interface{}
	if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
		return err
	}

	m.Set(key, value)
	return nil
}


func (m *Manager) Delete(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.data, key)

	if m.logger {
		slog.Debug("State deleted", "key", key)
	}
}


func (m *Manager) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data = make(map[string]interface{})

	if m.logger {
		slog.Debug("State cleared")
	}
}


func (m *Manager) Keys() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	keys := make([]string, 0, len(m.data))
	for key := range m.data {
		keys = append(keys, key)
	}

	return keys
}


func (m *Manager) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.data)
}


func (m *Manager) Exists(key string) bool {
	_, exists := m.Get(key)
	return exists
}


func (m *Manager) GetAll() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]interface{})
	for key, value := range m.data {
		result[key] = value
	}

	return result
}
