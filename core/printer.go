package core

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Printer struct {
	clearFns map[string]func()
	t        table.Writer
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

	t := table.NewWriter()
	t.SetStyle(table.StyleColoredBright)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Address", "TCP", "HTTP status", "HTTP Latency"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:     "HTTP status",
			WidthMin: 50,
			WidthMax: 50,
		},
	})

	return Printer{clearFns: clear, t: t}
}

func (p *Printer) ToTable(results *Store) {
	p.Clear()
	p.t.ResetRows()
	results.ForEach(func(r TestResult) bool {
		p.t.AppendRow(table.Row{
			r.Id,
			r.Tcp,
			r.HttpResponse,
			r.Duration,
		})
		return true
	})
	p.t.SortBy([]table.SortBy{{Name: "Address"}})
	p.t.AppendSeparator()
	p.t.Render()
}

func (p *Printer) Clear() {
	value, ok := p.clearFns[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
