package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type Cell struct {
	Value string `json:"value"`
	ID    int    `json:"id"`
}

type Sheet map[string]map[string]Cell

func main() {
	//inisialiasai Gin
	router := gin.Default()

	//membuat route dengan method GET
    router.POST("/generate", generate)
	router.GET("/test", func(c *gin.Context) {
		fmt.Println("File processed successfully!")
		// return response JSON
		c.JSON(200, gin.H{
			"result": "test Success!",
		})
	})

	//mulai server dengan port 3000
	router.Run(":3000")
}

func addFile(sheetData Sheet) error {
	// Open the existing spreadsheet file.
	f, err := excelize.OpenFile("source/Book1.xlsx")
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	// Defer close with error handling.
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("failed to close file:", err)
		}
	}()

	// Loop over each sheet and its cells in the sheetData parameter.
	for sheetName, cells := range sheetData {
		// Check if the sheet exists, create if it doesn't
		index, err := f.GetSheetIndex(sheetName)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		if index == -1 {
			// Create a new sheet if it doesn't exist
			f.NewSheet(sheetName)
		}

		// Loop through each cell and write the value to the Excel file
		for cell, data := range cells {
			if err := f.SetCellValue(sheetName, cell, data.Value); err != nil {
				return fmt.Errorf("failed to set cell value: %w", err)
			}
		}
	}

	// Save the updated spreadsheet to a new file
	if err := f.SaveAs("result/Book2.xlsx"); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func generate(c *gin.Context) {
	
	var sheetData Sheet
	
	// Get the JSON string from a form-encoded POST parameter called 'sheet'
	data := c.PostForm("data")
	
	// Unmarshal the JSON string into the 'sheetData' structure
	if err := json.Unmarshal([]byte(data), &sheetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	
	err := addFile(sheetData)
	if err != nil {
		// Handle the error appropriately, for example, logging it or exiting.
		fmt.Println("Error occurred:", err)
		return
	}

	// Return success response

	file := c.PostForm("file")
	name := c.PostForm("name")

	c.JSON(http.StatusOK, gin.H{
		"file": file,
		"name":  name,
		"data":  data,
	})
}