package service_test

import (
	"testing"

	"vyshkat/internal/domain"
	"vyshkat/internal/service"
)

func TestServiceCenter_Inspect(t *testing.T) {
	inspector := service.NewServiceCenter(70)

	t.Run("транспорт в отличном состоянии одобряется", func(t *testing.T) {
		scooter, err := domain.NewKickScooter("TST-01", "Scooter Good", 8, 90, 1.2)
		if err != nil {
			t.Fatalf("error creating scooter: %v", err)
		}

		verdict, err := inspector.Inspect(scooter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !verdict.Approved {
			t.Errorf("expected approved=true, got false: %s", verdict.Reason)
		}
	})

	t.Run("транспорт с состоянием ниже 70% отклоняется", func(t *testing.T) {
		scooter, err := domain.NewKickScooter("TST-02", "Scooter Damaged", 8, 65, 1.2)
		if err != nil {
			t.Fatalf("error creating scooter: %v", err)
		}

		verdict, err := inspector.Inspect(scooter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if verdict.Approved {
			t.Errorf("expected approved=false for condition 65%%, got approved=true")
		}
	})

	t.Run("пограничное значение 70% одобряется", func(t *testing.T) {
		bike, err := domain.NewBicycle("TST-03", "Bike Edge", 7, 70)
		if err != nil {
			t.Fatalf("error creating bike: %v", err)
		}

		verdict, err := inspector.Inspect(bike)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !verdict.Approved {
			t.Errorf("condition 70%% must be approved, got reason: %s", verdict.Reason)
		}
	})

	t.Run("электросамокат с некорректным расходом энергии отклоняется", func(t *testing.T) {
		brokenEnergyScooter, err := domain.NewKickScooter("TST-04", "Scooter Broken Battery", 8, 85, 0.0)
		if err != nil {
			t.Fatalf("error creating scooter: %v", err)
		}

		verdict, err := inspector.Inspect(brokenEnergyScooter)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if verdict.Approved {
			t.Errorf("device with 0 daily energy should be rejected due to battery sensor fault")
		}
	})

	t.Run("nil transport returns error", func(t *testing.T) {
		_, err := inspector.Inspect(nil)
		if err == nil {
			t.Errorf("expected error when passing nil transport")
		}
	})
}
