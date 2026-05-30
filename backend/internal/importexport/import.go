package importexport

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/product"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
)

// ProductCreator is the minimal interface needed to create products during import.
type ProductCreator interface {
	Create(ctx context.Context, tenantID kernel.TenantID, in productsrv.CreateInput) (*product.Product, error)
}

// ImportResult holds the outcome of a CSV import operation.
type ImportResult struct {
	Total    int           `json:"total"`
	Imported int           `json:"imported"`
	Errors   []ImportError `json:"errors"`
}

// ImportError describes a single row failure during import.
type ImportError struct {
	Row   int    `json:"row"`
	Error string `json:"error"`
}

// ImportService handles CSV import for products.
type ImportService struct {
	products ProductCreator
}

// NewImportService creates a new ImportService.
func NewImportService(products ProductCreator) *ImportService {
	return &ImportService{products: products}
}

// productCSVColumns defines the expected column order for product import.
// id, name, description, sku, price_amount, price_currency, category_id, tags, status, stock
const (
	colName         = 0
	colDescription  = 1
	colSKU          = 2
	colPriceAmount  = 3
	colCurrency     = 4
	colCategoryID   = 5
	colTags         = 6
	colStatus       = 7
	colStock        = 8
	productColCount = 9
)

// ImportProducts reads CSV from r and creates products for the tenant.
// The CSV must have a header row followed by data rows with columns:
// name, description, sku, price_amount, price_currency, category_id, tags, status, stock
func (s *ImportService) ImportProducts(ctx context.Context, tenantID kernel.TenantID, r io.Reader) (*ImportResult, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true

	// Skip header row
	if _, err := cr.Read(); err != nil {
		if err == io.EOF {
			return &ImportResult{Errors: []ImportError{}}, nil
		}
		return nil, errx.Wrap(err, "reading CSV header", errx.TypeValidation)
	}

	result := &ImportResult{
		Errors: []ImportError{},
	}

	rowNum := 1 // data rows start at row 2 (1-based, header is row 1)
	for {
		rowNum++
		record, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result.Total++
			result.Errors = append(result.Errors, ImportError{
				Row:   rowNum,
				Error: fmt.Sprintf("CSV parse error: %v", err),
			})
			continue
		}

		result.Total++

		if len(record) < productColCount {
			result.Errors = append(result.Errors, ImportError{
				Row:   rowNum,
				Error: fmt.Sprintf("expected %d columns, got %d", productColCount, len(record)),
			})
			continue
		}

		// Parse price_amount
		priceAmount, err := strconv.ParseInt(strings.TrimSpace(record[colPriceAmount]), 10, 64)
		if err != nil {
			result.Errors = append(result.Errors, ImportError{
				Row:   rowNum,
				Error: fmt.Sprintf("invalid price_amount %q: %v", record[colPriceAmount], err),
			})
			continue
		}

		// Parse stock
		stock, err := strconv.Atoi(strings.TrimSpace(record[colStock]))
		if err != nil {
			result.Errors = append(result.Errors, ImportError{
				Row:   rowNum,
				Error: fmt.Sprintf("invalid stock %q: %v", record[colStock], err),
			})
			continue
		}

		// Parse tags (semicolon-separated)
		var tags []string
		if raw := strings.TrimSpace(record[colTags]); raw != "" {
			for _, t := range strings.Split(raw, ";") {
				if trimmed := strings.TrimSpace(t); trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		in := productsrv.CreateInput{
			Name:        strings.TrimSpace(record[colName]),
			Description: strings.TrimSpace(record[colDescription]),
			SKU:         strings.TrimSpace(record[colSKU]),
			Price: kernel.Money{
				Amount:   priceAmount,
				Currency: strings.TrimSpace(record[colCurrency]),
			},
			CategoryID: kernel.CategoryID(strings.TrimSpace(record[colCategoryID])),
			Tags:       tags,
			Stock:      stock,
		}

		if _, err := s.products.Create(ctx, tenantID, in); err != nil {
			result.Errors = append(result.Errors, ImportError{
				Row:   rowNum,
				Error: err.Error(),
			})
			continue
		}

		result.Imported++
	}

	return result, nil
}
