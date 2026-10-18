package main

import (
	"flag"
	"fmt"
	"os"

	"vyshkat/internal/app"
	"vyshkat/internal/di"
)

func main() {
	interactiveFlag := flag.Bool("interactive", false, "Запустить в интерактивном режиме с меню")
	demoFlag := flag.Bool("demo", false, "Запустить демонстрационный прогон всех сценариев ТЗ")
	thresholdFlag := flag.Int("threshold", 70, "Минимальный порог технического состояния для техосмотра (по умолчанию 70%)")
	flag.Parse()

	// Инициализация графа зависимостей через DI-контейнер
	container, err := di.NewContainerBuilder().
		WithMinConditionThreshold(*thresholdFlag).
		Build()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Критическая ошибка инициализации зависимостей: %v\n", err)
		os.Exit(1)
	}

	cli := app.NewCLI(container.FleetService, os.Stdin, os.Stdout)

	if *demoFlag {
		cli.RunDemo()
		return
	}

	if *interactiveFlag {
		cli.RunInteractiveLoop()
		return
	}

	// Поведение по умолчанию: сначала наглядно выполняем все демо сценарии, демонстрируя работу всех бизнес-требований ТЗ
	cli.RunDemo()

	fmt.Println("\nДля запуска с интерактивным меню:")
	fmt.Println("  go run ./cmd/vyshkat --interactive")
}
