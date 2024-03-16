package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func NewExcel() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	for i, row := range [][]any{
		{"Nama", "Jenis Kelamin", "Tanggal lahir"},
		{"putu", "Laki-Laki", "1996-05-01"},
		{"yasa", "Laki-Laki", "1996-05-01"},
		{"putu", "Laki-Laki", "1996-05-01"},
		{"yasa", "Laki-Laki", "1996-05-01"},
		{"putu", "Laki-Laki", "1996-05-01"},
	} {
		cell, err := excelize.CoordinatesToCellName(1, i+1)
		if err != nil {
			return err
		}
		err = f.SetSheetRow("Sheet1", cell, &row)

		if err != nil {
			return err
		}
	}

	borderStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			excelize.Border{
				Type:  "left",
				Color: "000000",
				Style: 1,
			},
			excelize.Border{
				Type:  "right",
				Color: "000000",
				Style: 1,
			},
			excelize.Border{
				Type:  "top",
				Color: "000000",
				Style: 1,
			},
			excelize.Border{
				Type:  "bottom",
				Color: "000000",
				Style: 1,
			},
		},
	})
	if err != nil {
		return err
	}

	f.SetCellStyle("Sheet1", "A1", "C6", borderStyle)
	f.SetColWidth("Sheet1", "A", "C", 20)

	err = f.SaveAs("public/test.xlsx")
	if err != nil {
		return err
	}
	fmt.Println(f.Path)
	return nil
}
