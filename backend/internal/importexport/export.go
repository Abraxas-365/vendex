package importexport

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/Abraxas-365/vendex/internal/customer"
	"github.com/Abraxas-365/vendex/internal/customer/customersrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/order"
	"github.com/Abraxas-365/vendex/internal/order/ordersrv"
	"github.com/Abraxas-365/vendex/internal/product"
	"github.com/Abraxas-365/vendex/internal/product/productsrv"
)

// ProductLister is the minimal interface needed to list products for export.
type ProductLister interface {
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error)
}

// OrderLister is the minimal interface needed to list orders for export.
type OrderLister interface {
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error)
}

// CustomerLister is the minimal interface needed to list customers for export.
type CustomerLister interface {
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[customer.Customer], error)
}

// Ensure the real services satisfy the interfaces at compile time.
var _ ProductLister = (*productsrv.Service)(nil)
var _ OrderLister = (*ordersrv.Service)(nil)
var _ CustomerLister = (*customersrv.Service)(nil)

// ExportService handles CSV export for products, orders, and customers.
type ExportService struct {
	products  ProductLister
	orders    OrderLister
	customers CustomerLister
}

// NewExportService creates a new ExportService.
func NewExportService(products ProductLister, orders OrderLister, customers CustomerLister) *ExportService {
	return &ExportService{
		products:  products,
		orders:    orders,
		customers: customers,
	}
}

const exportPageSize = 100

// ExportProducts writes all products for the tenant as CSV to w.
// Columns: id, name, description, sku, price_amount, price_currency, category_id, tags, status, stock
func (s *ExportService) ExportProducts(ctx context.Context, tenantID kernel.TenantID, w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"id", "name", "description", "sku", "price_amount", "price_currency", "category_id", "tags", "status", "stock"}
	if err := cw.Write(header); err != nil {
		return errx.Wrap(err, "writing CSV header", errx.TypeInternal)
	}

	page := 1
	for {
		pg := kernel.NewPaginationOptions(page, exportPageSize)
		result, err := s.products.List(ctx, tenantID, pg)
		if err != nil {
			return errx.Wrap(err, "listing products for export", errx.TypeInternal)
		}

		for _, p := range result.Items {
			row := []string{
				string(p.ID),
				p.Name,
				p.Description,
				p.SKU,
				fmt.Sprintf("%d", p.Price.Amount),
				p.Price.Currency,
				string(p.CategoryID),
				strings.Join(p.Tags, ";"),
				string(p.Status),
				fmt.Sprintf("%d", p.Stock),
			}
			if err := cw.Write(row); err != nil {
				return errx.Wrap(err, "writing product CSV row", errx.TypeInternal)
			}
		}

		if page >= result.TotalPages {
			break
		}
		page++
	}

	return nil
}

// ExportOrders writes all orders for the tenant as CSV to w.
// Columns: id, customer_id, status, total_amount, total_currency, created_at
func (s *ExportService) ExportOrders(ctx context.Context, tenantID kernel.TenantID, w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"id", "customer_id", "status", "total_amount", "total_currency", "created_at"}
	if err := cw.Write(header); err != nil {
		return errx.Wrap(err, "writing CSV header", errx.TypeInternal)
	}

	page := 1
	for {
		pg := kernel.NewPaginationOptions(page, exportPageSize)
		result, err := s.orders.List(ctx, tenantID, pg)
		if err != nil {
			return errx.Wrap(err, "listing orders for export", errx.TypeInternal)
		}

		for _, o := range result.Items {
			row := []string{
				string(o.ID),
				string(o.CustomerID),
				string(o.Status),
				fmt.Sprintf("%d", o.TotalAmount.Amount),
				o.TotalAmount.Currency,
				o.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			}
			if err := cw.Write(row); err != nil {
				return errx.Wrap(err, "writing order CSV row", errx.TypeInternal)
			}
		}

		if page >= result.TotalPages {
			break
		}
		page++
	}

	return nil
}

// ExportCustomers writes all customers for the tenant as CSV to w.
// Columns: id, email, name, phone, created_at
func (s *ExportService) ExportCustomers(ctx context.Context, tenantID kernel.TenantID, w io.Writer) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	header := []string{"id", "email", "name", "phone", "created_at"}
	if err := cw.Write(header); err != nil {
		return errx.Wrap(err, "writing CSV header", errx.TypeInternal)
	}

	page := 1
	for {
		pg := kernel.NewPaginationOptions(page, exportPageSize)
		result, err := s.customers.List(ctx, tenantID, pg)
		if err != nil {
			return errx.Wrap(err, "listing customers for export", errx.TypeInternal)
		}

		for _, cu := range result.Items {
			row := []string{
				string(cu.ID),
				string(cu.Email),
				cu.Name,
				cu.Phone,
				cu.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
			}
			if err := cw.Write(row); err != nil {
				return errx.Wrap(err, "writing customer CSV row", errx.TypeInternal)
			}
		}

		if page >= result.TotalPages {
			break
		}
		page++
	}

	return nil
}
