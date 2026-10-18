package repository_test

import (
	"testing"

	"vyshkat/internal/domain"
	"vyshkat/internal/repository"
)

func TestInMemoryFleetRepository(t *testing.T) {
	repo := repository.NewInMemoryFleetRepository()

	s, _ := domain.NewKickScooter("S-100", "Ninebot", 7, 85, 1.2)
	h, _ := domain.NewHelmet("H-100", "CityRide", "M")

	t.Run("добавление и поиск транспорта", func(t *testing.T) {
		err := repo.AddTransport(s)
		if err != nil {
			t.Fatalf("unexpected error adding transport: %v", err)
		}

		found, ok := repo.FindTransportByInv("S-100")
		if !ok || found == nil {
			t.Fatalf("expected to find transport S-100")
		}

		if repo.Exists("S-100") != true {
			t.Errorf("expected Exists to be true")
		}
	})

	t.Run("дубликат транспорта вызывает ошибку", func(t *testing.T) {
		err := repo.AddTransport(s)
		if err == nil {
			t.Errorf("expected error when adding existing transport id")
		}
	})

	t.Run("добавление вещи и конфликт id", func(t *testing.T) {
		err := repo.AddItem(h)
		if err != nil {
			t.Fatalf("unexpected error adding item: %v", err)
		}

		// Попытка добавить вещь с тем же ID, что и транспорт
		conflictItem, _ := domain.NewHelmet("S-100", "Duplicate", "L")
		errConflict := repo.AddItem(conflictItem)
		if errConflict == nil {
			t.Errorf("expected conflict error between item and transport id")
		}
	})

	t.Run("GetAllInventory возвращает все сущности", func(t *testing.T) {
		all := repo.GetAllInventory()
		if len(all) != 2 {
			t.Errorf("expected 2 items, got %d", len(all))
		}
	})
}
