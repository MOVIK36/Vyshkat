package service_test

import (
	"testing"

	"vyshkat/internal/di"
	"vyshkat/internal/domain"
	"vyshkat/internal/repository"
	"vyshkat/internal/service"
)

// MockInspector реализует service.Inspector для изоляции и проверки поведения.
type MockInspector struct {
	ShouldApprove bool
	Reason        string
	Err           error
	CallCount     int
	LastInspected domain.Transport
}

func (m *MockInspector) Inspect(t domain.Transport) (domain.InspectionVerdict, error) {
	m.CallCount++
	m.LastInspected = t
	if m.Err != nil {
		return domain.InspectionVerdict{}, m.Err
	}
	return domain.InspectionVerdict{
		Approved: m.ShouldApprove,
		Reason:   m.Reason,
	}, nil
}

// TestFleetService_RegisterTransport_Success проверяет успешное добавление транспорта
// после прохождения техосмотра.
func TestFleetService_RegisterTransport_Success(t *testing.T) {
	repo := repository.NewInMemoryFleetRepository()
	mockInspector := &MockInspector{
		ShouldApprove: true,
		Reason:        "Все системы исправны",
	}

	fleetService := service.NewFleetService(repo, mockInspector)

	scooter, err := domain.NewKickScooter("SCOOT-001", "Ninebot Max G30", 8, 95, 1.2)
	if err != nil {
		t.Fatalf("неожиданная ошибка создания самоката: %v", err)
	}

	verdict, err := fleetService.RegisterTransport(scooter)
	if err != nil {
		t.Fatalf("ожидалась успешная регистрация, получена ошибка: %v", err)
	}

	if !verdict.Approved {
		t.Errorf("ожидалось одобрение техосмотра, получено отклонение")
	}

	if mockInspector.CallCount != 1 {
		t.Errorf("ожидался 1 вызов техосмотра, получено: %d", mockInspector.CallCount)
	}

	// Проверяем, что транспорт действительно появился в парке
	transports := fleetService.GetTransports()
	if len(transports) != 1 {
		t.Fatalf("ожидался 1 транспорт в парке, найдено: %d", len(transports))
	}
	if transports[0].InventoryNumber() != "SCOOT-001" {
		t.Errorf("неверный инвентарный номер: %s", transports[0].InventoryNumber())
	}
}

// TestFleetService_RegisterTransport_Rejected_UnderCondition фиксирует поведение,
// когда техосмотр отклоняет транспорт (техническое состояние неудовлетворительное).
// В этом случае транспорт НЕ должен попасть в парк!
func TestFleetService_RegisterTransport_Rejected_UnderCondition(t *testing.T) {
	repo := repository.NewInMemoryFleetRepository()
	mockInspector := &MockInspector{
		ShouldApprove: false,
		Reason:        "Критический износ тормозов: состояние 40%",
	}

	fleetService := service.NewFleetService(repo, mockInspector)

	scooter, err := domain.NewKickScooter("SCOOT-BROKEN", "Xiaomi Pro 2", 7, 40, 1.0)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	verdict, err := fleetService.RegisterTransport(scooter)
	if err != nil {
		t.Fatalf("вызов метода не должен возвращать системную ошибку при отказе: %v", err)
	}

	if verdict.Approved {
		t.Errorf("ожидался отказ техосмотра, но получено одобрение")
	}

	// Проверяем ключевое бизнес-требование ТЗ: транспорт НЕ попал в активный парк!
	transports := fleetService.GetTransports()
	if len(transports) != 0 {
		t.Errorf("парк должен быть пуст, но обнаружено устройств: %d", len(transports))
	}
}

// TestFleetService_DI_MockSubstitution демонстрирует работу DI-контейнера
// и подмену зависимости сервиса техосмотра через DI Builder.
func TestFleetService_DI_MockSubstitution(t *testing.T) {
	mockInspector := &MockInspector{
		ShouldApprove: true,
		Reason:        "Тестовый вердикт через DI контейнер",
	}

	// Сборка графа зависимостей через DI Container Builder с подменой инспектора
	container, err := di.NewContainerBuilder().
		WithInspector(mockInspector).
		Build()
	if err != nil {
		t.Fatalf("ошибка сборки DI контейнера: %v", err)
	}

	scooter, _ := domain.NewKickScooter("DI-001", "Whoosh Custom", 9, 100, 1.5)
	verdict, err := container.FleetService.RegisterTransport(scooter)
	if err != nil || !verdict.Approved {
		t.Fatalf("ожидалось успешное добавление через DI контейнер, verdict: %+v, err: %v", verdict, err)
	}

	if mockInspector.CallCount != 1 {
		t.Errorf("моковый инспектор в DI не был вызван")
	}
}

// TestFleetService_GetBeginnerFriendlyTransports фиксирует правило ТЗ:
// устройство подходит для новичка при характеристике "простота для новичка" >= 6.
func TestFleetService_GetBeginnerFriendlyTransports(t *testing.T) {
	mockInspector := &MockInspector{ShouldApprove: true}
	fleetService := service.NewFleetService(repository.NewInMemoryFleetRepository(), mockInspector)

	// Добавим 4 устройства с разными оценками простоты:
	// 1. Простой самокат (оценка 8) -> подходит
	// 2. Сложный спортивный электробайк (оценка 4) -> НЕ подходит
	// 3. Пограничный самокат (оценка 6) -> подходит
	// 4. Очень сложный моноцикл/прототип (оценка 3) -> НЕ подходит
	s1, _ := domain.NewKickScooter("S-1", "SimpleScooter", 8, 90, 1.0)
	s2, _ := domain.NewElectricBike("B-2", "SportSpeedEbike", 4, 90, 2.5)
	s3, _ := domain.NewKickScooter("S-3", "StandardScooter", 6, 85, 1.1)
	s4, _ := domain.NewBicycle("B-4", "ProTrackBike", 3, 90)

	_, _ = fleetService.RegisterTransport(s1)
	_, _ = fleetService.RegisterTransport(s2)
	_, _ = fleetService.RegisterTransport(s3)
	_, _ = fleetService.RegisterTransport(s4)

	beginnerList := fleetService.GetBeginnerFriendlyTransports()

	// Должны попасть только S-1 (8) и S-3 (6)
	if len(beginnerList) != 2 {
		t.Fatalf("ожидалось 2 устройства для новичков (оценка >= 6), получено: %d", len(beginnerList))
	}

	for _, item := range beginnerList {
		if item.BeginnerScore() < 6 {
			t.Errorf("устройство %s с оценкой %d ошибочно включено в список для новичков",
				item.InventoryNumber(), item.BeginnerScore())
		}
	}
}

// TestFleetService_GetTotalDailyEnergyKWh проверяет корректность подсчета суточного расхода энергии:
// учитывает электросамокаты, электробайки, зарядные шкафы, и корректно пропускает
// объекты без расхода энергии (обычные велосипеды и шлемы).
func TestFleetService_GetTotalDailyEnergyKWh(t *testing.T) {
	mockInspector := &MockInspector{ShouldApprove: true}
	fleetService := service.NewFleetService(repository.NewInMemoryFleetRepository(), mockInspector)

	// Электросамокат: 1.5 кВт·ч
	scooter, _ := domain.NewKickScooter("ES-1", "Ninebot", 7, 90, 1.5)
	_ = register(fleetService, scooter)

	// Электровелосипед: 2.3 кВт·ч
	ebike, _ := domain.NewElectricBike("EB-1", "Electra", 6, 88, 2.3)
	_ = register(fleetService, ebike)

	// Обычный велосипед: НЕ потребляет энергию (0 кВт·ч)
	bike, _ := domain.NewBicycle("BK-1", "Stels", 8, 95)
	_ = register(fleetService, bike)

	// Зарядный шкаф: 10.0 кВт·ч (вещь на балансе)
	station, _ := domain.NewChargingStation("CH-1", "Station Pokrovka", 8, 10.0)
	_ = fleetService.RegisterItem(station)

	// Шлем: вещь на балансе, 0 кВт·ч
	helmet, _ := domain.NewHelmet("HLM-1", "SafetyPro", "L")
	_ = fleetService.RegisterItem(helmet)

	totalEnergy := fleetService.GetTotalDailyEnergyKWh()
	expectedEnergy := 1.5 + 2.3 + 10.0 // 13.8 кВт·ч

	diff := totalEnergy - expectedEnergy
	if diff < -0.001 || diff > 0.001 {
		t.Errorf("ожидалось потребление %.2f кВт·ч, получено: %.2f кВт·ч", expectedEnergy, totalEnergy)
	}
}

// TestFleetService_GetAllInventory_IncludesTransportAndItems проверяет,
// что в сводный список инвентаря попадают и транспортные средства, и вещи.
func TestFleetService_GetAllInventory_IncludesTransportAndItems(t *testing.T) {
	mockInspector := &MockInspector{ShouldApprove: true}
	fleetService := service.NewFleetService(repository.NewInMemoryFleetRepository(), mockInspector)

	scooter, _ := domain.NewKickScooter("T-01", "KickScooter", 7, 90, 1.2)
	helmet, _ := domain.NewHelmet("I-01", "SafeRide", "M")
	dock, _ := domain.NewDockStation("I-02", "Корпус Покровка", 20)

	_ = register(fleetService, scooter)
	_ = fleetService.RegisterItem(helmet)
	_ = fleetService.RegisterItem(dock)

	all := fleetService.GetAllInventory()
	if len(all) != 3 {
		t.Fatalf("ожидалось 3 инвентарных объекта на балансе, получено: %d", len(all))
	}
}

func register(svc *service.FleetService, t domain.Transport) error {
	_, err := svc.RegisterTransport(t)
	return err
}

func TestFleetService_DuplicateInventoryNumber(t *testing.T) {
	mockInspector := &MockInspector{ShouldApprove: true}
	fleetService := service.NewFleetService(repository.NewInMemoryFleetRepository(), mockInspector)

	s1, _ := domain.NewKickScooter("DUP-1", "Ninebot", 7, 90, 1.2)
	s2, _ := domain.NewKickScooter("DUP-1", "Xiaomi", 8, 90, 1.0)

	_, err1 := fleetService.RegisterTransport(s1)
	if err1 != nil {
		t.Fatalf("первая регистрация должна пройти успешно: %v", err1)
	}

	_, err2 := fleetService.RegisterTransport(s2)
	if err2 == nil {
		t.Errorf("ожидалась ошибка дублирования инвентарного номера, но регистрация прошла успешно")
	}
}

func TestFleetService_StatsAndEdgeCases(t *testing.T) {
	mockInspector := &MockInspector{ShouldApprove: true}
	fleetService := service.NewFleetService(repository.NewInMemoryFleetRepository(), mockInspector)

	// nil checks
	if _, err := fleetService.RegisterTransport(nil); err == nil {
		t.Errorf("expected error registering nil transport")
	}
	if err := fleetService.RegisterItem(nil); err == nil {
		t.Errorf("expected error registering nil item")
	}

	scooter, _ := domain.NewKickScooter("ST-1", "Scoot", 7, 85, 1.0)
	helmet, _ := domain.NewHelmet("HL-1", "Helmet", "M")

	_, _ = fleetService.RegisterTransport(scooter)
	_ = fleetService.RegisterItem(helmet)

	items := fleetService.GetItems()
	if len(items) != 1 {
		t.Errorf("expected 1 item, got %d", len(items))
	}

	stats := fleetService.GetFleetStats()
	if stats.TotalTransports != 1 || stats.TotalItems != 1 || stats.TotalInventory != 2 {
		t.Errorf("unexpected stats: %+v", stats)
	}
	if stats.BeginnerFriendly != 1 {
		t.Errorf("expected 1 beginner friendly, got %d", stats.BeginnerFriendly)
	}
}
