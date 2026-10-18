package di_test

import (
	"testing"

	"vyshkat/internal/di"
	"vyshkat/internal/domain"
	"vyshkat/internal/repository"
)

type StubInspector struct {
	verdict domain.InspectionVerdict
}

func (s *StubInspector) Inspect(domain.Transport) (domain.InspectionVerdict, error) {
	return s.verdict, nil
}

func TestDIContainer(t *testing.T) {
	t.Run("DefaultContainer собирает все зависимости", func(t *testing.T) {
		container, err := di.NewDefaultContainer()
		if err != nil {
			t.Fatalf("failed to build default container: %v", err)
		}

		if container.Repo == nil {
			t.Errorf("Repo is nil")
		}
		if container.Inspector == nil {
			t.Errorf("Inspector is nil")
		}
		if container.FleetService == nil {
			t.Errorf("FleetService is nil")
		}
	})

	t.Run("ContainerBuilder позволяет подменять зависимости", func(t *testing.T) {
		customRepo := repository.NewInMemoryFleetRepository()
		stubInspector := &StubInspector{
			verdict: domain.InspectionVerdict{Approved: true, Reason: "Stub approved"},
		}

		container, err := di.NewContainerBuilder().
			WithRepository(customRepo).
			WithInspector(stubInspector).
			WithMinConditionThreshold(80).
			Build()

		if err != nil {
			t.Fatalf("failed to build container: %v", err)
		}

		if container.Repo != customRepo {
			t.Errorf("custom repo was not injected")
		}
		if container.Inspector != stubInspector {
			t.Errorf("custom inspector was not injected")
		}
	})
}
