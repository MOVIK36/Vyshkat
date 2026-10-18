package domain

import "fmt"

// Helmet описывает защитный шлем
type Helmet struct {
	invNum string
	model  string
	size   string // S, M, L
}

func NewHelmet(invNum, model, size string) (*Helmet, error) {
	if invNum == "" {
		return nil, fmt.Errorf("инвентарный номер шлема не может быть пустым")
	}
	return &Helmet{
		invNum: invNum,
		model:  model,
		size:   size,
	}, nil
}

func (h *Helmet) InventoryNumber() string {
	return h.invNum
}

func (h *Helmet) Title() string {
	return fmt.Sprintf("Шлем %s (Размер %s)", h.model, h.size)
}

func (h *Helmet) Category() string {
	return "Инвентарь / Вещь"
}

// ChargingStation описывает стационарный шкаф быстрой зарядки
type ChargingStation struct {
	invNum      string
	name        string
	portsCount  int
	dailyEnergy float64 // кВт·ч в сутки
}

func NewChargingStation(invNum, name string, ports int, dailyEnergy float64) (*ChargingStation, error) {
	if invNum == "" {
		return nil, fmt.Errorf("инвентарный номер зарядной станции не может быть пустым")
	}
	if ports <= 0 {
		return nil, fmt.Errorf("количество портов должно быть > 0")
	}
	return &ChargingStation{
		invNum:      invNum,
		name:        name,
		portsCount:  ports,
		dailyEnergy: dailyEnergy,
	}, nil
}

func (c *ChargingStation) InventoryNumber() string {
	return c.invNum
}

func (c *ChargingStation) Title() string {
	return fmt.Sprintf("Зарядный шкаф %s (%d портов)", c.name, c.portsCount)
}

func (c *ChargingStation) Category() string {
	return "Инвентарь / Зарядная инфраструктура"
}

func (c *ChargingStation) DailyEnergyKWh() float64 {
	return c.dailyEnergy
}

// DockStation док-станция для парковки
type DockStation struct {
	invNum   string
	location string
	capacity int
}

func NewDockStation(invNum, location string, capacity int) (*DockStation, error) {
	if invNum == "" {
		return nil, fmt.Errorf("инвентарный номер док-станции не может быть пустым")
	}
	return &DockStation{
		invNum:   invNum,
		location: location,
		capacity: capacity,
	}, nil
}

func (d *DockStation) InventoryNumber() string {
	return d.invNum
}

func (d *DockStation) Title() string {
	return fmt.Sprintf("Док-станция парковки: %s (вместимость %d)", d.location, d.capacity)
}

func (d *DockStation) Category() string {
	return "Инвентарь / Инфраструктура"
}
