package domain

// InventoryItem представляет любой объект, находящийся на балансе
// (как транспортные средства, так и инвентарь: шлемы, зарядные станции)
type InventoryItem interface {
	// InventoryNumber возвращает уникальный инвентарный номер объекта
	InventoryNumber() string
	// Title возвращает понятное наименование объекта (например, "Электросамокат Ninebot Max G30")
	Title() string
	// Category возвращает категорию объекта ("Транспорт" или "Инвентарь/Вещь")
	Category() string
}

// EnergyConsumer описывает устройства, потребляющие электроэнергию (суточный расход в кВт·ч)
type EnergyConsumer interface {
	DailyEnergyKWh() float64
}

// BeginnerFriendly описывает технику с оценкой доступности для новичков (1–10)
type BeginnerFriendly interface {
	// BeginnerScore возвращает оценку простоты управления от 1 (сложно) до 10 (очень легко)
	BeginnerScore() int
	// IsSuitableForBeginner определяет, подходит ли устройство для новичка (оценка >= 6)
	IsSuitableForBeginner() bool
}

// Transport объединяет прокатные транспортные средства
// Любой транспорт является инвентарным объектом и оценивается по простоте для новичка
type Transport interface {
	InventoryItem
	BeginnerFriendly
	// TransportType возвращает тип транспорта (например, "Электросамокат", "Велосипед")
	TransportType() string
	// TechnicalCondition возвращает процент технического состояния (0-100%)
	TechnicalCondition() int
}

// InspectionVerdict представляет результат прохождения техосмотра
type InspectionVerdict struct {
	Approved bool   // Одобрен ли допуск в активный парк
	Reason   string // Обоснование вердикта или обнаруженные дефекты
}
