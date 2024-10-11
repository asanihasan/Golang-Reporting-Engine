package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)


func main() {
    f, err := excelize.OpenFile("Book1.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer func() {
        // Close the spreadsheet.
        if err := f.Close(); err != nil {
            fmt.Println(err)
        }
    }()

	f.SetCellValue("Sheet1", "A2", "Hello world.")
    
    // Save the spreadsheet with the origin path.
    if err := f.SaveAs("Book2.xlsx"); err != nil {
        fmt.Println(err)
    }
}