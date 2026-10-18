package di

import (
	"fmt"

	"vyshkat/internal/repository"
	"vyshkat/internal/service"
)

// Container инкапсулирует собранные зависимости приложения "ВышКат"
type Container struct {
	Repo         repository.FleetRepository
	Inspector    service.Inspector
	FleetService *service.FleetService
}

// ContainerBuilder предоставляет расширяемый способ настройки зависимостей
type ContainerBuilder struct {
	repo                  repository.FleetRepository
	inspector             service.Inspector
	minConditionThreshold int
}

// NewContainerBuilder инициализирует новый билдер контейнера
func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		minConditionThreshold: 70, // стандартный порог безопасности
	}
}

// WithRepository позволяет переопределить репозиторий
func (b *ContainerBuilder) WithRepository(repo repository.FleetRepository) *ContainerBuilder {
	b.repo = repo
	return b
}

// WithInspector позволяет подменить сервис техосмотра
func (b *ContainerBuilder) WithInspector(inspector service.Inspector) *ContainerBuilder {
	b.inspector = inspector
	return b
}

// WithMinConditionThreshold настраивает порог прохождения техосмотра для стандартного ServiceCenter
func (b *ContainerBuilder) WithMinConditionThreshold(threshold int) *ContainerBuilder {
	b.minConditionThreshold = threshold
	return b
}

// Build собирает граф зависимостей и возвращает готовый контейнер
func (b *ContainerBuilder) Build() (*Container, error) {
	repo := b.repo
	if repo == nil {
		repo = repository.NewInMemoryFleetRepository()
	}

	inspector := b.inspector
	if inspector == nil {
		inspector = service.NewServiceCenter(b.minConditionThreshold)
	}

	if repo == nil || inspector == nil {
		return nil, fmt.Errorf("не удалось инициализировать зависимости контейнера")
	}

	fleetSvc := service.NewFleetService(repo, inspector)

	return &Container{
		Repo:         repo,
		Inspector:    inspector,
		FleetService: fleetSvc,
	}, nil
}

// NewDefaultContainer создает контейнер со стандартными рабочими зависимостями
func NewDefaultContainer() (*Container, error) {
	return NewContainerBuilder().Build()
}
