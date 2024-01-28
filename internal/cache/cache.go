package cache

import (
	"WB_Tech_level_0/internal/models"
	"errors"
)

type Cache struct {
	Data map[string]models.Data `json:"data"`
}

// Инициализация КЭШа
func New() *Cache {
	return &Cache{
		Data: make(map[string]models.Data),
	}
}

// Запись в КЭШ
func (c *Cache) InCache(data models.Data) error {
	_, exist := c.Data[data.Order.Order_uid]
	if !exist {
		c.Data[data.Order.Order_uid] = data
		return nil
	}
	return errors.New("element is already exists")
}

// Получение данных из КЭШа
func (c *Cache) GetCache(id string) (models.Data, bool) {
	data, exist := c.Data[id]
	return data, exist
}
