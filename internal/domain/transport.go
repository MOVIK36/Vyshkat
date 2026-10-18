package domain

import (
	"fmt"
)

// BaseTransport содержит общие поля для всех видов транспорта
type BaseTransport struct {
	invNum    string
	title     string
	typeName  string
	score     int
	condition int // 0 - 100%
}

func (b *BaseTransport) InventoryNumber() string {
	return b.invNum
}

func (b *BaseTransport) Title() string {
	return b.title
}

func (b *BaseTransport) Category() string {
	return "Транспорт"
}

func (b *BaseTransport) TransportType() string {
	return b.typeName
}

func (b *BaseTransport) BeginnerScore() int {
	return b.score
}

// IsSuitableForBeginner реализует требование ТЗ: при значении >= 6 можно выдавать новичкам
func (b *BaseTransport) IsSuitableForBeginner() bool {
	return b.score >= 6
}

func (b *BaseTransport) TechnicalCondition() int {
	return b.condition
}

// KickScooter описывает электросамокат
type KickScooter struct {
	BaseTransport
	dailyEnergy float64 // кВт·ч в сутки
}

func NewKickScooter(invNum, model string, score int, condition int, dailyEnergy float64) (*KickScooter, error) {
	if score < 1 || score > 10 {
		return nil, fmt.Errorf("оценка простоты должна быть в диапазоне [1, 10], получено: %d", score)
	}
	if condition < 0 || condition > 100 {
		return nil, fmt.Errorf("техническое состояние должно быть в диапазоне [0, 100]%%, получено: %d", condition)
	}
	return &KickScooter{
		BaseTransport: BaseTransport{
			invNum:    invNum,
			title:     model,
			typeName:  "Электросамокат",
			score:     score,
			condition: condition,
		},
		dailyEnergy: dailyEnergy,
	}, nil
}

func (s *KickScooter) DailyEnergyKWh() float64 {
	return s.dailyEnergy
}

// ElectricBike реализует: Transport, InventoryItem, BeginnerFriendly, EnergyConsumer
type ElectricBike struct {
	BaseTransport
	dailyEnergy float64
}

func NewElectricBike(invNum, model string, score int, condition int, dailyEnergy float64) (*ElectricBike, error) {
	if score < 1 || score > 10 {
		return nil, fmt.Errorf("оценка простоты должна быть в диапазоне [1, 10], получено: %d", score)
	}
	if condition < 0 || condition > 100 {
		return nil, fmt.Errorf("техническое состояние должно быть в диапазоне [0, 100]%%, получено: %d", condition)
	}
	return &ElectricBike{
		BaseTransport: BaseTransport{
			invNum:    invNum,
			title:     model,
			typeName:  "Электровелосипед",
			score:     score,
			condition: condition,
		},
		dailyEnergy: dailyEnergy,
	}, nil
}

func (e *ElectricBike) DailyEnergyKWh() float64 {
	return e.dailyEnergy
}

// Bicycle описывает механический велосипед без электродвигателя
type Bicycle struct {
	BaseTransport
}

func NewBicycle(invNum, model string, score int, condition int) (*Bicycle, error) {
	if score < 1 || score > 10 {
		return nil, fmt.Errorf("оценка простоты должна быть в диапазоне [1, 10], получено: %d", score)
	}
	if condition < 0 || condition > 100 {
		return nil, fmt.Errorf("техническое состояние должно быть в диапазоне [0, 100]%%, получено: %d", condition)
	}
	return &Bicycle{
		BaseTransport: BaseTransport{
			invNum:    invNum,
			title:     model,
			typeName:  "Классический велосипед",
			score:     score,
			condition: condition,
		},
	}, nil
}
