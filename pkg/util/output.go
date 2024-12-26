package util

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

func DefaultTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(w)
	table.SetCenterSeparator("  ")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetTablePadding("  ")
	table.SetAutoFormatHeaders(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetFooterAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetNoWhiteSpace(true)
	return table
}
