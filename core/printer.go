package core

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Printer struct {
	clearFns map[string]func()
}

func NewPrinter() Printer {
	clear := make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	return Printer{clearFns: clear}
}

func (p *Printer) ToTable(results *[]TestResult) {
	p.Clear()
	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Address", "Connection", "Status"})
	for _, r := range *results {
		t.AppendRow(table.Row{r.url.Host, strconv.FormatInt(r.duration.Milliseconds(), 10) + "ms", r.status})
	}

	t.AppendSeparator()
	t.Render()
}

func (p *Printer) Clear() {
	value, ok := p.clearFns[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
