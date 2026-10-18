package service

import (
	"fmt"

	"vyshkat/internal/domain"
)

// Inspector — контракт сервиса технического контроля (ServiceCenter)
// Позволяет легко подменять реализацию в тестах
type Inspector interface {
	Inspect(t domain.Transport) (domain.InspectionVerdict, error)
}

// ServiceCenter — стандартная реализация сервиса техосмотра
// Проверяет техническое состояние устройства перед включением в парк
type ServiceCenter struct {
	minConditionThreshold int
}

// NewServiceCenter создает новый сервисный центр с заданным порогом допуска (по умолчанию 70%)
func NewServiceCenter(minConditionThreshold int) *ServiceCenter {
	if minConditionThreshold <= 0 {
		minConditionThreshold = 70
	}
	return &ServiceCenter{
		minConditionThreshold: minConditionThreshold,
	}
}

// Inspect проводит диагностику транспорта
func (s *ServiceCenter) Inspect(t domain.Transport) (domain.InspectionVerdict, error) {
	if t == nil {
		return domain.InspectionVerdict{}, fmt.Errorf("передан пустой объект транспорта")
	}

	cond := t.TechnicalCondition()
	if cond < s.minConditionThreshold {
		return domain.InspectionVerdict{
			Approved: false,
			Reason: fmt.Sprintf("Отказ: техническое состояние %d%% ниже минимального порога безопасности (%d%%)",
				cond, s.minConditionThreshold),
		}, nil
	}

	// Для энергопотребляющего транспорта проверяем энергопотребление
	if ec, ok := t.(domain.EnergyConsumer); ok {
		if ec.DailyEnergyKWh() <= 0 {
			return domain.InspectionVerdict{
				Approved: false,
				Reason:   "Отказ: обнаружена неисправность батареи или датчика (расход энергии <= 0 кВт·ч)",
			}, nil
		}
	}

	return domain.InspectionVerdict{
		Approved: true,
		Reason:   fmt.Sprintf("Одобрено: транспорт исправен, состояние %d%%", cond),
	}, nil
}
