package app

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"vyshkat/internal/domain"
	"vyshkat/internal/service"
)

// CLI реализует консольный интерфейс приложения
type CLI struct {
	fleetSvc *service.FleetService
	reader   *bufio.Reader
	writer   io.Writer
}

func NewCLI(fleetSvc *service.FleetService, in io.Reader, out io.Writer) *CLI {
	return &CLI{
		fleetSvc: fleetSvc,
		reader:   bufio.NewReader(in),
		writer:   out,
	}
}

// RunDemo выполняет прогон основных сценариев работы парка
func (c *CLI) RunDemo() {
	c.printHeader()

	c.println("--- [ДЕМО-СЦЕНАРИЙ 1] Регистрация пилотного парка транспорта ---")
	s1, _ := domain.NewKickScooter("HSE-SCOOT-001", "Ninebot Max G30 (Покровка)", 8, 98, 1.25)
	s2, _ := domain.NewKickScooter("HSE-SCOOT-002", "Xiaomi 4 Pro (Шаболовка)", 7, 92, 1.10)
	s3, _ := domain.NewKickScooter("HSE-SCOOT-003", "Ninebot F2 Pro (Мясницкая)", 6, 88, 1.30)
	ebike, _ := domain.NewElectricBike("HSE-EBIKE-001", "Electra Townie Go (ФКН)", 5, 85, 2.20)
	bike, _ := domain.NewBicycle("HSE-BIKE-001", "Forward Quadro Sport (Кирпич)", 9, 95)

	c.acceptTransportWithLog(s1)
	c.acceptTransportWithLog(s2)
	c.acceptTransportWithLog(s3)
	c.acceptTransportWithLog(ebike)
	c.acceptTransportWithLog(bike)

	c.println("\n--- [ДЕМО-СЦЕНАРИЙ 2] Отклонение техники с дефектами на техосмотре ---")
	brokenScooter, _ := domain.NewKickScooter("HSE-SCOOT-999", "Ninebot ES4 (износ 45%)", 8, 45, 1.20)
	c.acceptTransportWithLog(brokenScooter)

	c.println("\n--- [ДЕМО-СЦЕНАРИЙ 3] Постановка на баланс вещей и инфраструктуры ---")
	helmet1, _ := domain.NewHelmet("HSE-HELMET-01", "Abus Urban-I 3.0", "M")
	helmet2, _ := domain.NewHelmet("HSE-HELMET-02", "Livall Smart Shield", "L")
	chargingStation, _ := domain.NewChargingStation("HSE-CHARGER-01", "Зарядная станция Покровка", 12, 14.50)
	dockStation, _ := domain.NewDockStation("HSE-DOCK-01", "Парковочный хаб ФКН", 15)

	_ = c.fleetSvc.RegisterItem(helmet1)
	c.println(fmt.Sprintf("  [+] Вещь на балансе: %-32s | Инв.№: %s", helmet1.Title(), helmet1.InventoryNumber()))
	_ = c.fleetSvc.RegisterItem(helmet2)
	c.println(fmt.Sprintf("  [+] Вещь на балансе: %-32s | Инв.№: %s", helmet2.Title(), helmet2.InventoryNumber()))
	_ = c.fleetSvc.RegisterItem(chargingStation)
	c.println(fmt.Sprintf("  [+] Вещь на балансе: %-32s | Инв.№: %s | Расход: %.2f кВт*ч",
		chargingStation.Title(), chargingStation.InventoryNumber(), chargingStation.DailyEnergyKWh()))
	_ = c.fleetSvc.RegisterItem(dockStation)
	c.println(fmt.Sprintf("  [+] Вещь на балансе: %-32s | Инв.№: %s", dockStation.Title(), dockStation.InventoryNumber()))

	c.println("\n--- [ДЕМО-СЦЕНАРИЙ 4] Отчет по суточному энергопотреблению флота ---")
	c.PrintEnergyReport()

	c.println("\n--- [ДЕМО-СЦЕНАРИЙ 5] Техника для новичков (критерий простоты >= 6) ---")
	c.PrintBeginnerReport()

	c.println("\n--- [ДЕМО-СЦЕНАРИЙ 6] Реестр всех инвентарных объектов парка ---")
	c.PrintFullInventory()
}

func (c *CLI) acceptTransportWithLog(t domain.Transport) {
	c.println(fmt.Sprintf("  * Диагностика: %s [%s, состояние: %d%%, простота: %d/10]",
		t.Title(), t.InventoryNumber(), t.TechnicalCondition(), t.BeginnerScore()))

	verdict, err := c.fleetSvc.RegisterTransport(t)
	if err != nil {
		c.println(fmt.Sprintf("    [ОШИБКА]: %v", err))
		return
	}

	if verdict.Approved {
		c.println(fmt.Sprintf("    [ТЕХОСМОТР: ОК] %s -> принят в парк", verdict.Reason))
	} else {
		c.println(fmt.Sprintf("    [ТЕХОСМОТР ОТКЛОНЕН] %s -> отклонен", verdict.Reason))
	}
}

// PrintEnergyReport выводит энергопотребление всех подключенных устройств.
func (c *CLI) PrintEnergyReport() {
	c.println(fmt.Sprintf("%-16s | %-38s | %-12s", "Инв. номер", "Наименование", "Расход (кВт*ч/сут)"))
	c.println("-----------------+----------------------------------------+------------------")

	var total float64
	for _, item := range c.fleetSvc.GetAllInventory() {
		if ec, ok := item.(domain.EnergyConsumer); ok {
			kwh := ec.DailyEnergyKWh()
			total += kwh
			c.println(fmt.Sprintf("%-16s | %-38s | %10.2f кВт*ч", item.InventoryNumber(), item.Title(), kwh))
		}
	}
	c.println("-----------------+----------------------------------------+------------------")
	c.println(fmt.Sprintf("СУММАРНОЕ ПОТРЕБЛЕНИЕ ФЛОТА:                               %10.2f кВт*ч\n", total))
}

// PrintBeginnerReport выводит список устройств, подходящих новичкам.
func (c *CLI) PrintBeginnerReport() {
	beginners := c.fleetSvc.GetBeginnerFriendlyTransports()
	c.println(fmt.Sprintf("Техника для новичков (найдено единиц: %d):", len(beginners)))
	c.println(fmt.Sprintf("%-16s | %-32s | %-18s | %-10s", "Инв. номер", "Модель", "Тип транспорта", "Простота (1-10)"))
	c.println("-----------------+----------------------------------+--------------------+----------------")
	for _, t := range beginners {
		c.println(fmt.Sprintf("%-16s | %-32s | %-18s | %d/10",
			t.InventoryNumber(), t.Title(), t.TransportType(), t.BeginnerScore()))
	}
	c.println("")
}

// PrintFullInventory выводит все материальные ценности на балансе.
func (c *CLI) PrintFullInventory() {
	all := c.fleetSvc.GetAllInventory()
	c.println(fmt.Sprintf("Всего объектов на балансе: %d", len(all)))
	c.println(fmt.Sprintf("%-16s | %-32s | %-24s", "Инв. номер", "Наименование / Модель", "Категория учёта"))
	c.println("-----------------+----------------------------------+-------------------------")
	for _, it := range all {
		c.println(fmt.Sprintf("%-16s | %-32s | %-24s", it.InventoryNumber(), it.Title(), it.Category()))
	}
	c.println("")
}

// RunInteractiveLoop запускает меню консоли.
func (c *CLI) RunInteractiveLoop() {
	c.printHeader()

	for {
		c.println("\nМеню управления парком:")
		c.println("1. Принять новый транспорт (через техосмотр)")
		c.println("2. Поставить на баланс вещь/инвентарь")
		c.println("3. Отчет по суточному расходу энергии")
		c.println("4. Список техники для новичков")
		c.println("5. Реестр всех объектов и инвентарных номеров")
		c.println("6. Запустить демонстрационный сценарий")
		c.println("0. Завершить работу")
		c.print("Выберите действие [0-6]: ")

		line, err := c.readLine()
		if err != nil {
			return
		}

		choice := strings.TrimSpace(line)
		switch choice {
		case "1":
			c.handleRegisterTransport()
		case "2":
			c.handleRegisterItem()
		case "3":
			c.PrintEnergyReport()
		case "4":
			c.PrintBeginnerReport()
		case "5":
			c.PrintFullInventory()
		case "6":
			c.RunDemo()
		case "0", "q", "exit":
			c.println("Работа завершена. До новых встреч.")
			return
		default:
			c.println("Неизвестная команда.")
		}
	}
}

func (c *CLI) handleRegisterTransport() {
	c.println("\nРегистрация транспорта:")
	c.println("1 - Электросамокат, 2 - Электровелосипед, 3 - Обычный велосипед")
	c.print("Тип: ")
	tType, _ := c.readLine()

	c.print("Инвентарный номер: ")
	invNum, _ := c.readLine()
	invNum = strings.TrimSpace(invNum)

	c.print("Модель/Наименование: ")
	model, _ := c.readLine()
	model = strings.TrimSpace(model)

	c.print("Простота для новичка (1-10): ")
	scoreStr, _ := c.readLine()
	score, _ := strconv.Atoi(strings.TrimSpace(scoreStr))

	c.print("Техническое состояние в % (порог допуска 70%): ")
	condStr, _ := c.readLine()
	condition, _ := strconv.Atoi(strings.TrimSpace(condStr))

	var t domain.Transport
	var err error

	switch strings.TrimSpace(tType) {
	case "1":
		c.print("Суточный расход энергии (кВт*ч): ")
		eStr, _ := c.readLine()
		energy, _ := strconv.ParseFloat(strings.TrimSpace(eStr), 64)
		t, err = domain.NewKickScooter(invNum, model, score, condition, energy)
	case "2":
		c.print("Суточный расход энергии (кВт*ч): ")
		eStr, _ := c.readLine()
		energy, _ := strconv.ParseFloat(strings.TrimSpace(eStr), 64)
		t, err = domain.NewElectricBike(invNum, model, score, condition, energy)
	case "3":
		t, err = domain.NewBicycle(invNum, model, score, condition)
	default:
		c.println("Неизвестный тип транспорта.")
		return
	}

	if err != nil {
		c.println(fmt.Sprintf("Ошибка валидации: %v", err))
		return
	}

	c.acceptTransportWithLog(t)
}

func (c *CLI) handleRegisterItem() {
	c.println("\nРегистрация вещи/инвентаря:")
	c.println("1 - Шлем, 2 - Зарядная станция, 3 - Док-станция")
	c.print("Тип: ")
	iType, _ := c.readLine()

	c.print("Инвентарный номер: ")
	invNum, _ := c.readLine()
	invNum = strings.TrimSpace(invNum)

	switch strings.TrimSpace(iType) {
	case "1":
		c.print("Модель шлема: ")
		model, _ := c.readLine()
		c.print("Размер (S/M/L): ")
		size, _ := c.readLine()
		h, err := domain.NewHelmet(invNum, strings.TrimSpace(model), strings.TrimSpace(size))
		if err == nil {
			err = c.fleetSvc.RegisterItem(h)
		}
		if err != nil {
			c.println(fmt.Sprintf("Ошибка: %v", err))
		} else {
			c.println("Шлем добавлен на баланс.")
		}
	case "2":
		c.print("Название станции: ")
		name, _ := c.readLine()
		c.print("Число портов: ")
		pStr, _ := c.readLine()
		ports, _ := strconv.Atoi(strings.TrimSpace(pStr))
		c.print("Суточный расход энергии (кВт*ч): ")
		eStr, _ := c.readLine()
		energy, _ := strconv.ParseFloat(strings.TrimSpace(eStr), 64)

		st, err := domain.NewChargingStation(invNum, strings.TrimSpace(name), ports, energy)
		if err == nil {
			err = c.fleetSvc.RegisterItem(st)
		}
		if err != nil {
			c.println(fmt.Sprintf("Ошибка: %v", err))
		} else {
			c.println("Зарядная станция добавлена на баланс.")
		}
	case "3":
		c.print("Локация док-станции: ")
		loc, _ := c.readLine()
		c.print("Вместимость: ")
		capStr, _ := c.readLine()
		capVal, _ := strconv.Atoi(strings.TrimSpace(capStr))

		dock, err := domain.NewDockStation(invNum, strings.TrimSpace(loc), capVal)
		if err == nil {
			err = c.fleetSvc.RegisterItem(dock)
		}
		if err != nil {
			c.println(fmt.Sprintf("Ошибка: %v", err))
		} else {
			c.println("Док-станция добавлена на баланс.")
		}
	default:
		c.println("Неизвестный тип.")
	}
}

func (c *CLI) printHeader() {
	c.println("=== Сервис кикшеринга \"ВЫШКАТ\" (НИУ ВШЭ) ===")
}

func (c *CLI) print(s string) {
	_, _ = fmt.Fprint(c.writer, s)
}

func (c *CLI) println(s string) {
	_, _ = fmt.Fprintln(c.writer, s)
}

func (c *CLI) readLine() (string, error) {
	return c.reader.ReadString('\n')
}
