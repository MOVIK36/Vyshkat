package repository

import (
	"fmt"
	"sync"

	"vyshkat/internal/domain"
)

// FleetRepository описывает контракт хранилища техники и сопутствующего инвентаря
type FleetRepository interface {
	AddTransport(t domain.Transport) error
	GetTransports() []domain.Transport
	FindTransportByInv(invNum string) (domain.Transport, bool)

	AddItem(item domain.InventoryItem) error
	GetItems() []domain.InventoryItem

	GetAllInventory() []domain.InventoryItem
	Exists(invNum string) bool
}

// InMemoryFleetRepository - потокобезопасная реализация хранилища в оперативной памяти
type InMemoryFleetRepository struct {
	mu         sync.RWMutex
	transports map[string]domain.Transport
	items      map[string]domain.InventoryItem
}

// NewInMemoryFleetRepository создает новый экземпляр репозитория
func NewInMemoryFleetRepository() *InMemoryFleetRepository {
	return &InMemoryFleetRepository{
		transports: make(map[string]domain.Transport),
		items:      make(map[string]domain.InventoryItem),
	}
}

func (r *InMemoryFleetRepository) Exists(invNum string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.transports[invNum]; ok {
		return true
	}
	if _, ok := r.items[invNum]; ok {
		return true
	}
	return false
}

func (r *InMemoryFleetRepository) AddTransport(t domain.Transport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inv := t.InventoryNumber()
	if _, exists := r.transports[inv]; exists {
		return fmt.Errorf("объект с инвентарным номером '%s' уже зарегистрирован в парке", inv)
	}
	if _, exists := r.items[inv]; exists {
		return fmt.Errorf("инвентарный номер '%s' уже занят другим инвентарем", inv)
	}

	r.transports[inv] = t
	return nil
}

func (r *InMemoryFleetRepository) GetTransports() []domain.Transport {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.Transport, 0, len(r.transports))
	for _, t := range r.transports {
		list = append(list, t)
	}
	return list
}

func (r *InMemoryFleetRepository) FindTransportByInv(invNum string) (domain.Transport, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.transports[invNum]
	return t, ok
}

func (r *InMemoryFleetRepository) AddItem(item domain.InventoryItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	inv := item.InventoryNumber()
	if _, exists := r.transports[inv]; exists {
		return fmt.Errorf("инвентарный номер '%s' уже занят транспортом", inv)
	}
	if _, exists := r.items[inv]; exists {
		return fmt.Errorf("вещь с инвентарным номером '%s' уже существует", inv)
	}

	r.items[inv] = item
	return nil
}

func (r *InMemoryFleetRepository) GetItems() []domain.InventoryItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.InventoryItem, 0, len(r.items))
	for _, it := range r.items {
		list = append(list, it)
	}
	return list
}

func (r *InMemoryFleetRepository) GetAllInventory() []domain.InventoryItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	total := make([]domain.InventoryItem, 0, len(r.transports)+len(r.items))
	for _, t := range r.transports {
		total = append(total, t)
	}
	for _, it := range r.items {
		total = append(total, it)
	}
	return total
}
