package export

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/rs/zerolog/log"
	"github.com/tealeg/xlsx"
)

func ExportToPDF(c *gin.Context, tableOr string, data [][]string, headers []string, widths []float64, fileName string, title string, fontPath string) {
	// Initialize gofpdf
	pdf := gofpdf.New(tableOr, "mm", "A4", "")
	pdf.AddUTF8Font("Cairo", "R", fontPath)

	// Validate widths length
	if len(widths) != len(headers) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Widths array length must match headers length",
		})
		return
	}

	// Add a page
	pdf.AddPage()
	fmt.Print(len(data))
	// Add title
	pdf.SetFont("Times", "B", 14)
	pdf.Cell(0, 10, title)
	pdf.Ln(20)

	// Add table headers
	pdf.SetFillColor(200, 200, 200)
	pdf.SetFont("Times", "B", 9)

	for i, header := range headers {
		pdf.CellFormat(widths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Add table content
	pdf.SetFont("Cairo", "R", 9)
	if len(data) > 0 {
		for _, row := range data {
			for i, field := range row {
				pdf.CellFormat(widths[i], 10, field, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
		}
	} else {
		totalWidth := 0.0
		for _, w := range widths {
			totalWidth += w
		}
		pdf.CellFormat(totalWidth, 10, "No data available", "1", 1, "C", false, 0, "")
	}

	// Add footer
	currentTime := time.Now().Format("15:04:05 02-01-2006")
	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 8)
	pdf.CellFormat(0, 10, fmt.Sprintf("Generated on: %s", currentTime), "", 1, "C", false, 0, "")

	// Set response headers and output PDF
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileName))
	err := pdf.Output(c.Writer)
	if err != nil {
		log.Debug().Msgf("Failed to generate PDF: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
	}
}

func ExportToExcel(c *gin.Context, data [][]string, headers []string, fileName string) {
	// Create a new Excel file
	xlFile := xlsx.NewFile()

	// Add a sheet
	sheet, err := xlFile.AddSheet(fmt.Sprintf("%s Data", fileName))
	if err != nil {
		log.Debug().Msgf("Failed to create Excel sheet: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet"})
		return
	}

	// Create header row
	headerRow := sheet.AddRow()
	for _, header := range headers {
		headerRow.AddCell().Value = header
	}

	// Add data rows
	for _, row := range data {
		dataRow := sheet.AddRow()
		for _, field := range row {
			dataRow.AddCell().Value = field
		}
	}

	// Set response headers and output Excel
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", fileName))
	err = xlFile.Write(c.Writer)
	if err != nil {
		log.Debug().Msgf("Failed to write Excel file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file"})
	}
}
