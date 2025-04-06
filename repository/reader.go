package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/graduate-work-mirea/data-collector-service/entity"
	"go.uber.org/zap"
)

// DataReader is responsible for reading data from the dataset file
type DataReader struct {
	datasetPath string
	logger      *zap.SugaredLogger
}

// NewDataReader creates a new instance of DataReader
func NewDataReader(datasetPath string, logger *zap.SugaredLogger) *DataReader {
	return &DataReader{
		datasetPath: datasetPath,
		logger:      logger,
	}
}

// ReadData reads all data from the dataset file and returns it as a slice of ProductData
func (r *DataReader) ReadData() ([]entity.ProductData, error) {
	r.logger.Infof("Reading data from file: %s", r.datasetPath)

	file, err := os.Open(r.datasetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dataset file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var products []entity.ProductData

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			r.logger.Warnf("Error reading CSV record: %v, skipping", err)
			continue
		}

		product, err := r.parseRecord(record)
		if err != nil {
			r.logger.Warnf("Error parsing record: %v, skipping", err)
			continue
		}

		products = append(products, product)
	}

	r.logger.Infof("Successfully read %d products from dataset", len(products))
	return products, nil
}

// parseRecord converts a CSV record to a ProductData struct
func (r *DataReader) parseRecord(record []string) (entity.ProductData, error) {
	if len(record) < 16 {
		return entity.ProductData{}, fmt.Errorf("record has insufficient fields: %d", len(record))
	}

	salesQty, _ := strconv.Atoi(record[3])
	price, _ := strconv.ParseFloat(record[4], 64)

	// Handle empty fields for original price and discount
	var originalPrice float64
	if record[5] != "" {
		originalPrice, _ = strconv.ParseFloat(record[5], 64)
	}

	var discountPercentage float64
	if record[6] != "" {
		discountPercentage, _ = strconv.ParseFloat(record[6], 64)
	}

	stockLevel, _ := strconv.Atoi(record[7])
	customerRating, _ := strconv.ParseFloat(record[10], 64)
	reviewCount, _ := strconv.Atoi(record[11])
	deliveryDays, _ := strconv.Atoi(record[12])

	isWeekend := false
	if strings.ToLower(record[14]) == "1" || strings.ToLower(record[14]) == "true" {
		isWeekend = true
	}

	isHoliday := false
	if strings.ToLower(record[15]) == "1" || strings.ToLower(record[15]) == "true" {
		isHoliday = true
	}

	return entity.ProductData{
		ProductName:        record[0],
		Brand:              record[1],
		Date:               record[2],
		SalesQuantity:      salesQty,
		Price:              price,
		OriginalPrice:      originalPrice,
		DiscountPercentage: discountPercentage,
		StockLevel:         stockLevel,
		Region:             record[8],
		Category:           record[9],
		CustomerRating:     customerRating,
		ReviewCount:        reviewCount,
		DeliveryDays:       deliveryDays,
		Seller:             record[13],
		IsWeekend:          isWeekend,
		IsHoliday:          isHoliday,
	}, nil
}
