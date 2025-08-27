package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type Cell struct {
	Value string `json:"value"`
	ID    string `json:"id"`
}

type Sheet map[string]map[string]Cell

type Font struct {
	Bold   bool   `json:"bold"`
	Italic bool   `json:"italic"`
	Size   int    `json:"size"`
	Color  string `json:"color"`
}

type Fill struct {
	Type    string   `json:"type"`
	Color   []string `json:"color"`
	Pattern int      `json:"pattern"`
}

type Border struct {
	Type  string `json:"type"`
	Color string `json:"color"`
	Style int    `json:"style"`
}

type Alignment struct {
	Horizontal string `json:"horizontal"`
	Vertical   string `json:"vertical"`
}

type StyleOptions struct {
	Font      *Font      `json:"font"`
	Fill      *Fill      `json:"fill"`
	Border    []Border   `json:"border"`
	Alignment *Alignment `json:"alignment"`
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	router := gin.Default()

	router.GET("/templates", template)
	router.POST("/upload", upload)
	router.POST("/generate", generate)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": "test Success!",
		})
	})

	router.Run(":6969")
}

// parseNumeric tries to coerce a string into float64 (handles "1,234.56" etc)
func parseNumeric(s string) (float64, bool) {
	u := strings.TrimSpace(s)
	if u == "" {
		return 0, false
	}
	// very basic normalization: remove thousands separators
	u = strings.ReplaceAll(u, ",", "")
	// optional: handle trailing/leading spaces already trimmed
	if n, err := strconv.ParseFloat(u, 64); err == nil {
		return n, true
	}
	return 0, false
}

func addFile(sheetData Sheet, file string, styling map[string]StyleOptions) (string, error) {
	name := generateRandomString(32)

	// Open the existing spreadsheet file (template file).
	f, err := excelize.OpenFile("source/" + file)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", file, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("failed to close file:", err)
		}
	}()

	// Write data into sheets
	for sheetName, cells := range sheetData {
		// create sheet if missing
		index, err := f.GetSheetIndex(sheetName)
		if err != nil {
			return "", fmt.Errorf("failed to get sheet index for %s: %w", sheetName, err)
		}
		if index == -1 {
			f.NewSheet(sheetName)
		}

		for addr, data := range cells {
			// Write numbers as numbers, else as strings
			if n, ok := parseNumeric(data.Value); ok {
				if err := f.SetCellValue(sheetName, addr, n); err != nil {
					return "", fmt.Errorf("failed to set numeric value for %s in %s: %w", addr, sheetName, err)
				}
			} else {
				if err := f.SetCellValue(sheetName, addr, data.Value); err != nil {
					return "", fmt.Errorf("failed to set string value for %s in %s: %w", addr, sheetName, err)
				}
			}

			// Apply style if provided
			if styleOptions, exists := styling[data.ID]; exists {
				excelStyle, err := CreateExcelStyle(f, &styleOptions)
				if err == nil {
					_ = f.SetCellStyle(sheetName, addr, addr, excelStyle)
				}
			}
		}
	}

	// ⚙️ Ensure Excel recalculates everything on open
	_ = f.SetCalcProps(&excelize.CalcPropsOptions{
		CalcMode:       excelize.StringPtr("auto"), // "manual", "auto", "autoNoTable"
		FullCalcOnLoad: excelize.BoolPtr(true),     // force full rebuild on open
	})

	// 🔄 Drop stale calc chain so Excel rebuilds dependencies
	_ = f.DeleteCalcChain()

	// Ensure output dir exists
	dirPath := "result"
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Save the updated spreadsheet
	if err := f.SaveAs("result/" + name + ".xlsx"); err != nil {
		return "", fmt.Errorf("failed to save file as %s: %w", name, err)
	}

	return name, nil
}

func generate(c *gin.Context) {
	var sheetData Sheet
	var styles map[string]StyleOptions

	// Get the JSON string from form fields
	data := c.PostForm("data")
	style := c.PostForm("style")

	// Unmarshal inputs
	if err := json.Unmarshal([]byte(data), &sheetData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	if err := json.Unmarshal([]byte(style), &styles); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	file := c.PostForm("file")
	name := c.PostForm("name")

	result, err := addFile(sheetData, file, styles)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed modify file"})
		return
	}

	downloadFile(c, "result/"+result+".xlsx", name+".xlsx")
}

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error occurred while retrieving file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file."})
		return
	}

	path := "source/" + file.Filename
	if err := c.SaveUploadedFile(file, path); err != nil {
		fmt.Println("Error occurred while saving file:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file."})
		return
	}

	if _, err := excelize.OpenFile(path); err != nil {
		fmt.Println("Invalid xlsx file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uploaded file is not a valid xlsx file."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and is a valid xlsx file!"})
}

func template(c *gin.Context) {
	folderPath := "./source"
	files, err := os.ReadDir(folderPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	c.JSON(http.StatusOK, gin.H{"files": fileNames})
}

func generateRandomString(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func downloadFile(c *gin.Context, sourceFile string, downloadFileName string) {
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Could not read file"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+downloadFileName)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))

	if _, err := c.Writer.Write(content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to write content"})
		return
	}

	if err := os.Remove(sourceFile); err != nil {
		fmt.Println("Invalid xlsx file:", err)
	}
}

func CreateExcelStyle(f *excelize.File, styleOptions *StyleOptions) (int, error) {
	style := &excelize.Style{}

	if styleOptions != nil && styleOptions.Font != nil {
		style.Font = &excelize.Font{
			Bold:   styleOptions.Font.Bold,
			Italic: styleOptions.Font.Italic,
			Size:   float64(styleOptions.Font.Size),
			Color:  styleOptions.Font.Color,
		}
	}

	if styleOptions != nil && styleOptions.Fill != nil {
		style.Fill = excelize.Fill{
			Type:    styleOptions.Fill.Type,
			Color:   styleOptions.Fill.Color,
			Pattern: styleOptions.Fill.Pattern,
		}
	}

	if styleOptions != nil && len(styleOptions.Border) > 0 {
		var borders []excelize.Border
		for _, b := range styleOptions.Border {
			borders = append(borders, excelize.Border{
				Type:  b.Type,
				Color: b.Color,
				Style: b.Style,
			})
		}
		style.Border = borders
	}

	if styleOptions != nil && styleOptions.Alignment != nil {
		style.Alignment = &excelize.Alignment{
			Horizontal: styleOptions.Alignment.Horizontal,
			Vertical:   styleOptions.Alignment.Vertical,
		}
	}

	return f.NewStyle(style)
}
