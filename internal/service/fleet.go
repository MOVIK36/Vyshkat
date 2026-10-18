package service

import (
	"fmt"
	"sort"

	"vyshkat/internal/domain"
	"vyshkat/internal/repository"
)

// FleetStats агрегирует ключевые показатели работы сервиса
type FleetStats struct {
	TotalTransports  int
	TotalItems       int
	TotalInventory   int
	TotalDailyEnergy float64
	BeginnerFriendly int
}

// FleetService инкапсулирует бизнес-логику управления парком
type FleetService struct {
	repo      repository.FleetRepository
	inspector Inspector
}

// NewFleetService создает экземпляр сервиса с внедренными зависимостями
func NewFleetService(repo repository.FleetRepository, inspector Inspector) *FleetService {
	return &FleetService{
		repo:      repo,
		inspector: inspector,
	}
}

// RegisterTransport проводит обязательный техосмотр перед приёмом в парк
func (s *FleetService) RegisterTransport(t domain.Transport) (domain.InspectionVerdict, error) {
	if t == nil {
		return domain.InspectionVerdict{}, fmt.Errorf("транспорт не указан")
	}

	if s.repo.Exists(t.InventoryNumber()) {
		return domain.InspectionVerdict{}, fmt.Errorf("инвентарный номер '%s' уже зарегистрирован", t.InventoryNumber())
	}

	// 1. Обязательный техосмотр
	verdict, err := s.inspector.Inspect(t)
	if err != nil {
		return domain.InspectionVerdict{}, fmt.Errorf("ошибка в процессе техосмотра: %w", err)
	}

	if !verdict.Approved {
		// Отклонено техосмотром: НЕ ставится на баланс
		return verdict, nil
	}

	// 2. Постановка на баланс
	if err := s.repo.AddTransport(t); err != nil {
		return verdict, fmt.Errorf("ошибка сохранения в парке: %w", err)
	}

	return verdict, nil
}

// RegisterItem регистрирует вещь (шлем, зарядный шкаф) на баланс парка
func (s *FleetService) RegisterItem(item domain.InventoryItem) error {
	if item == nil {
		return fmt.Errorf("инвентарный объект не указан")
	}
	return s.repo.AddItem(item)
}

// GetTotalDailyEnergyKWh вычисляет суммарное суточное потребление энергии всего парка
func (s *FleetService) GetTotalDailyEnergyKWh() float64 {
	var total float64

	// Сканируем весь инвентарь: если объект реализует EnergyConsumer, плюсуем энергию
	for _, item := range s.repo.GetAllInventory() {
		if ec, ok := item.(domain.EnergyConsumer); ok {
			total += ec.DailyEnergyKWh()
		}
	}

	return total
}

// GetBeginnerFriendlyTransports возвращает список транспорта, подходящего для новичков
func (s *FleetService) GetBeginnerFriendlyTransports() []domain.Transport {
	transports := s.repo.GetTransports()
	result := make([]domain.Transport, 0)

	for _, t := range transports {
		if t.IsSuitableForBeginner() {
			result = append(result, t)
		}
	}

	// Сортировка по убыванию простоты (самые простые сверху)
	sort.Slice(result, func(i, j int) bool {
		return result[i].BeginnerScore() > result[j].BeginnerScore()
	})

	return result
}

// GetAllInventory возвращает все объекты на балансе (технику и вещи)
func (s *FleetService) GetAllInventory() []domain.InventoryItem {
	items := s.repo.GetAllInventory()
	sort.Slice(items, func(i, j int) bool {
		return items[i].InventoryNumber() < items[j].InventoryNumber()
	})
	return items
}

// GetTransports возвращает весь транспорт в активном парке
func (s *FleetService) GetTransports() []domain.Transport {
	transports := s.repo.GetTransports()
	sort.Slice(transports, func(i, j int) bool {
		return transports[i].InventoryNumber() < transports[j].InventoryNumber()
	})
	return transports
}

// GetItems возвращает все вещи на балансе
func (s *FleetService) GetItems() []domain.InventoryItem {
	items := s.repo.GetItems()
	sort.Slice(items, func(i, j int) bool {
		return items[i].InventoryNumber() < items[j].InventoryNumber()
	})
	return items
}

// GetFleetStats возвращает сводные аналитические данные парка
func (s *FleetService) GetFleetStats() FleetStats {
	transports := s.repo.GetTransports()
	items := s.repo.GetItems()

	beginnerCount := 0
	for _, t := range transports {
		if t.IsSuitableForBeginner() {
			beginnerCount++
		}
	}

	return FleetStats{
		TotalTransports:  len(transports),
		TotalItems:       len(items),
		TotalInventory:   len(transports) + len(items),
		TotalDailyEnergy: s.GetTotalDailyEnergyKWh(),
		BeginnerFriendly: beginnerCount,
	}
}
