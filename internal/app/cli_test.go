package app_test

import (
	"bytes"
	"strings"
	"testing"

	"vyshkat/internal/app"
	"vyshkat/internal/di"
)

func TestCLI_RunDemo(t *testing.T) {
	container, err := di.NewDefaultContainer()
	if err != nil {
		t.Fatalf("failed to init container: %v", err)
	}

	var in bytes.Buffer
	var out bytes.Buffer

	cli := app.NewCLI(container.FleetService, &in, &out)
	cli.RunDemo()

	output := out.String()

	// Проверяем наличие ключевых секций ТЗ в выводе
	if !strings.Contains(output, "ВЫШКАТ") {
		t.Errorf("output should contain ВЫШКАТ")
	}
	if !strings.Contains(output, "ДЕМО-СЦЕНАРИЙ 1") {
		t.Errorf("output should contain scenario 1")
	}
	if !strings.Contains(output, "ДЕМО-СЦЕНАРИЙ 2") {
		t.Errorf("output should contain scenario 2")
	}
	if !strings.Contains(output, "ТЕХОСМОТР ОТКЛОНЕН") {
		t.Errorf("output should show rejected inspection")
	}
	if !strings.Contains(output, "СУММАРНОЕ ПОТРЕБЛЕНИЕ ФЛОТА") {
		t.Errorf("output should show energy report")
	}
	if !strings.Contains(output, "Техника для новичков") {
		t.Errorf("output should show beginner friendly report")
	}
}

func TestCLI_InteractiveQuit(t *testing.T) {
	container, _ := di.NewDefaultContainer()

	in := bytes.NewBufferString("0\n")
	var out bytes.Buffer

	cli := app.NewCLI(container.FleetService, in, &out)
	cli.RunInteractiveLoop()

	if !strings.Contains(out.String(), "До новых встреч") {
		t.Errorf("expected exit message")
	}
}
