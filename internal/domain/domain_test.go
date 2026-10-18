package domain_test

import (
	"testing"

	"vyshkat/internal/domain"
)

func TestDomain_Validation(t *testing.T) {
	t.Run("валидация самоката", func(t *testing.T) {
		_, err := domain.NewKickScooter("S1", "Model", 11, 80, 1.0)
		if err == nil {
			t.Errorf("ожидалась ошибка для score > 10")
		}

		_, err2 := domain.NewKickScooter("S2", "Model", 5, 120, 1.0)
		if err2 == nil {
			t.Errorf("ожидалась ошибка для condition > 100")
		}

		s, errValid := domain.NewKickScooter("S3", "Model", 5, 80, 1.0)
		if errValid != nil {
			t.Fatalf("неожиданная ошибка: %v", errValid)
		}
		if s.IsSuitableForBeginner() != false {
			t.Errorf("оценка 5 не должна подходить для новичка")
		}
	})

	t.Run("валидация электровелосипеда", func(t *testing.T) {
		_, err := domain.NewElectricBike("EB1", "E-Bike", 0, 80, 2.0)
		if err == nil {
			t.Errorf("ожидалась ошибка для score < 1")
		}

		_, errCond := domain.NewElectricBike("EB1", "E-Bike", 7, -5, 2.0)
		if errCond == nil {
			t.Errorf("ожидалась ошибка для condition < 0")
		}

		eb, errValid := domain.NewElectricBike("EB2", "E-Bike", 7, 80, 2.5)
		if errValid != nil || eb.DailyEnergyKWh() != 2.5 {
			t.Errorf("ошибка создания валидного электровелосипеда")
		}
	})

	t.Run("валидация велосипеда", func(t *testing.T) {
		_, err := domain.NewBicycle("B1", "Bike", 12, 90)
		if err == nil {
			t.Errorf("ожидалась ошибка для score > 10")
		}

		b, errValid := domain.NewBicycle("B2", "Bike", 9, 90)
		if errValid != nil {
			t.Fatalf("ошибка: %v", errValid)
		}
		if !b.IsSuitableForBeginner() {
			t.Errorf("оценка 9 должна подходить для новичка")
		}
	})

	t.Run("валидация инвентаря и вещей", func(t *testing.T) {
		_, errH := domain.NewHelmet("", "Helmet", "M")
		if errH == nil {
			t.Errorf("ожидалась ошибка для пустого инвентарного номера шлема")
		}

		_, errCS := domain.NewChargingStation("CS1", "Station", 0, 5.0)
		if errCS == nil {
			t.Errorf("ожидалась ошибка для ports <= 0")
		}

		_, errCS2 := domain.NewChargingStation("", "Station", 5, 5.0)
		if errCS2 == nil {
			t.Errorf("ожидалась ошибка для пустого номера зарядной станции")
		}

		_, errD := domain.NewDockStation("", "Loc", 10)
		if errD == nil {
			t.Errorf("ожидалась ошибка для пустого номера док-станции")
		}
	})
}
