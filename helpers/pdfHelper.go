package helpers

import (
	"fmt"
	"restaurant-management/models"

	"github.com/johnfercher/maroto/v2"
	// "github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func GeneratePdfFromData(invoice models.Invoice, order models.Order, restaurant models.Restaurant) (core.Document, error) {
	config := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		WithTopMargin(15).
		Build()

	m := maroto.New(config)

	// Header: Restaurant Name & Invoice #
	m.AddRow(15,
		text.NewCol(12, restaurant.Name, props.Text{
			Align: align.Center, Size: 16, Style: fontstyle.Bold,
		}),
	)

	m.AddRow(10,
		text.NewCol(12, fmt.Sprintf("Invoice #%d", invoice.ID), props.Text{
			Align: align.Center, Style: fontstyle.Bold, Size: 12,
		}),
	)

	// Invoice Meta Info
	m.AddRow(10,
		text.NewCol(6, fmt.Sprintf("Date: %s", invoice.CreatedAt.Format("02 Jan 2006")), props.Text{Size: 10}),
		text.NewCol(6, fmt.Sprintf("Order ID: %d", invoice.OrderID), props.Text{Align: align.Right, Size: 10}),
	)

	// Table Header
	m.AddRow(10,
		text.NewCol(4, "Item", props.Text{Style: fontstyle.Bold, Left: 1}),
		text.NewCol(2, "Qty", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Unit Price", props.Text{Style: fontstyle.Bold, Align: align.Center}),
		text.NewCol(3, "Subtotal", props.Text{Style: fontstyle.Bold, Align: align.Right, Right: 1}),
	)

	// Items
	for i, item := range order.OrderItems {
		background := &props.Color{Red: 245, Green: 245, Blue: 245}
		style := &props.Cell{}
		if i%2 == 0 {
			style.BackgroundColor = background
		}

		m.AddRow(8,
			text.NewCol(4, item.FoodName, props.Text{Top: 2, VerticalPadding: 2, Left: 1}).WithStyle(style),
			text.NewCol(2, fmt.Sprintf("%d", item.Quantity), props.Text{Top: 2, Align: align.Center, VerticalPadding: 2}).WithStyle(style),
			text.NewCol(3, fmt.Sprintf("%.2f", item.UnitPrice), props.Text{Top: 2, Align: align.Center, VerticalPadding: 2}).WithStyle(style),
			text.NewCol(3, fmt.Sprintf("%.2f", item.SubTotal), props.Text{Top: 2, Align: align.Right, VerticalPadding: 2, Right: 1}).WithStyle(style),
		)
	}

	m.AddRow(10, line.NewCol(12))

	// Totals
	m.AddRow(8,
		text.NewCol(10, "Total:", props.Text{Align: align.Right, Style: fontstyle.Bold}),
		text.NewCol(2, fmt.Sprintf("%.2f", invoice.Total), props.Text{Align: align.Right, Style: fontstyle.Bold, Right: 1}),
	)

	// Add Footer
	m.AddRow(40,
		signature.NewCol(6, "Authorized Signature", props.Signature{FontFamily: fontfamily.Courier}),
		// code.NewQrCol(6, "https://codeheim.io", props.Rect{ // QR Code
		// 	Percent: 75,
		// 	Center:  true,
		// }),
	)

	return m.Generate()
}
