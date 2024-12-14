package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Report struct {
	writer table.Writer
}

func NewReport(title string) *Report {
	report := Report{}

	report.writer = table.NewWriter()
	report.writer.SetOutputMirror(os.Stdout)
	report.writer.Style().Options.SeparateRows = false
	report.writer.Style().Options.SeparateColumns = true
	report.writer.Style().Box.MiddleVertical = "="
	report.writer.Style().Box.BottomSeparator = "="
	report.writer.Style().Box.MiddleSeparator = "="
	report.writer.Style().Box.MiddleHorizontal = "="

	report.writer.SetTitle(fmt.Sprintf("\n%v", title))
	report.writer.Style().Title.Align = text.AlignCenter

	report.writer.Style().Options.DrawBorder = false
	report.writer.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:    1,
			WidthMin:  30,
			WidthMax:  64,
			AutoMerge: true,
			Align:     text.AlignRight,
		},
		{
			Number:    2,
			WidthMin:  30,
			AutoMerge: true,
			Align:     text.AlignLeft,
		},
	})
	return &report
}

func (report *Report) AddEntry(key string, val interface{}) {
	report.writer.AppendRow(table.Row{key, val})
}

func (report *Report) AddSection(name string) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true, AutoMergeAlign: text.AlignCenter}
	report.writer.AppendRow(table.Row{"", ""}, rowConfigAutoMerge)
	report.writer.AppendSeparator()
	report.writer.AppendRow(table.Row{name, name}, rowConfigAutoMerge)
	report.writer.AppendSeparator()
}

func (report *Report) Render() {
	report.writer.Render()
}
