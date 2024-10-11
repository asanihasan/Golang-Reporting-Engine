package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
    router.GET("/templates", template)
    router.POST("/upload", upload)
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

func addFile(sheetData Sheet, file string, name string) error {
	// Open the existing spreadsheet file (template file).
	f, err := excelize.OpenFile("source/" + file) // Use the provided template file
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file, err)
	}

	// Defer close with error handling.
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("failed to close file:", err)
		}
	}()

	// Loop over each sheet and its cells in the sheetData parameter.
	for sheetName, cells := range sheetData {
		// Check if the sheet exists, create it if it doesn't
		index, err := f.GetSheetIndex(sheetName)
		if err != nil {
			return fmt.Errorf("failed to get sheet index for %s: %w", sheetName, err)
		}
		if index == -1 {
			// Create a new sheet if it doesn't exist
			f.NewSheet(sheetName)
		}

		// Loop through each cell and write the value to the Excel file
		for cell, data := range cells {
			if err := f.SetCellValue(sheetName, cell, data.Value); err != nil {
				return fmt.Errorf("failed to set cell value for %s in %s: %w", cell, sheetName, err)
			}
		}
	}

	// Save the updated spreadsheet to the new file (with the provided name).
	if err := f.SaveAs("result/" + name + ".xlsx"); err != nil {
		return fmt.Errorf("failed to save file as %s: %w", name, err)
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

	file := c.PostForm("file")
	name := c.PostForm("name")
	
	err := addFile(sheetData, file, name)
	if err != nil {
		// Handle the error appropriately, for example, logging it or exiting.
		fmt.Println("Error occurred:", err)
		return
	}

	// Return success response


	c.JSON(http.StatusOK, gin.H{
		"file": file,
		"name":  name,
		"data":  data,
	})
}

func upload(c *gin.Context) {
	// Get the uploaded file from the form-data POST request
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error occurred while retrieving file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file."})
		return
	}

	// Save the uploaded file to a specific path
	path := "source/" + file.Filename
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		fmt.Println("Error occurred while saving file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file."})
		return
	}

	// Try to open the file as an xlsx to validate if it's a valid Excel file
	_, err = excelize.OpenFile(path)
	if err != nil {
		fmt.Println("Invalid xlsx file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not a valid xlsx file."})
		return
	}

	// If everything is valid, return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded and is a valid xlsx file!",
	})
}

func template(c *gin.Context) {
	folderPath := "./source"

	// Get the list of files in the directory
	files, err := os.ReadDir(folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read directory",
		})
		return
	}

	// Create a slice to store the file names
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() { // Only include files, not directories
			fileNames = append(fileNames, file.Name())
		}
	}

	// Return the list of file names as a JSON array
	c.JSON(http.StatusOK, gin.H{
		"files": fileNames,
	})
}