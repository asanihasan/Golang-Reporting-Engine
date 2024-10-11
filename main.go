package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)


func main() {
	//inisialiasai Gin
	router := gin.Default()

	//membuat route dengan method GET
	router.GET("/", func(c *gin.Context) {

		err := addFile()
		if err != nil {
			// Handle the error appropriately, for example, logging it or exiting.
			fmt.Println("Error occurred:", err)
			return
		}
		
		fmt.Println("File processed successfully!")
		// return response JSON
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	//mulai server dengan port 3000
	router.Run(":3000")
	

}

func addFile() error {
	// Open the spreadsheet file.
	f, err := excelize.OpenFile("source/Book1.xlsx")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	
	// Defer close with better error handling.
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("failed to close file:", err)
		}
	}()

	// Set a value in the spreadsheet.
	if err := f.SetCellValue("Sheet1", "A2", "Hello world."); err != nil {
		return fmt.Errorf("failed to set cell value: %w", err)
	}

	// Save the spreadsheet to a new path.
	if err := f.SaveAs("result/Book2.xlsx"); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}
