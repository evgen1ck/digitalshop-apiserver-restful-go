package handlers_v1

import (
	"log"
	"net/http"
	"strings"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
)

type Variant struct {
	VariantName     string  `json:"variant_name"`
	VariantId       string  `json:"variant_id"`
	Service         string  `json:"service"`
	ServiceSvgUrl   string  `json:"service_svg_url"`
	State           string  `json:"state"`
	Item            string  `json:"item"`
	Mask            string  `json:"mask"`
	TextQuantity    string  `json:"text_quantity"`
	Price           float64 `json:"price"`
	DiscountMoney   float64 `json:"discount_money"`
	DiscountPercent int     `json:"discount_percent"`
	FinalPrice      float64 `json:"final_price"`
}

type Subtype struct {
	SubtypeName string    `json:"subtype_name"`
	Type        string    `json:"type"`
	Variants    []Variant `json:"variants"`
}

type Product struct {
	ProductName     string    `json:"product_name"`
	ProductId       string    `json:"product_id"`
	ProductImageUrl string    `json:"product_image_url"`
	Description     string    `json:"description"`
	Subtypes        []Subtype `json:"subtypes"`
}

func (rs *Resolver) ProductsData(w http.ResponseWriter, r *http.Request) {
	fullUrl := tl.GetFullURL(r, rs.App.Config.App.Debug)

	value, err := tl.UrlGetParam(fullUrl, "mainpage")
	if err != nil {
		rs.App.Logger.NewWarn("error in get url param", err)
		api_v1.RespondWithInternalServerError(w)
		return
	} else if value != "" {
		rows, err := storage.GetProducts(r.Context(), rs.App.Postgres)
		if err != nil {
			rs.App.Logger.NewWarn("error in get products", err)
			api_v1.RespondWithInternalServerError(w)
			return
		}

		productsMap := make(map[string]*Product)

		for rows.Next() {
			var v Variant
			var s Subtype
			var p Product

			err := rows.Scan(
				&s.Type,
				&s.SubtypeName,
				&v.Service,
				&p.ProductName,
				&v.VariantName,
				&v.State,
				&v.Price,
				&v.DiscountMoney,
				&v.DiscountPercent,
				&v.FinalPrice,
				&v.Item,
				&v.Mask,
				&v.TextQuantity,
				&p.Description,
				&p.ProductId,
				&v.VariantId,
			)

			if err != nil {
				log.Fatal(err)
			}

			if _, ok := productsMap[p.ProductName]; !ok {
				productsMap[p.ProductName] = &Product{
					ProductName:     p.ProductName,
					ProductId:       p.ProductId,
					ProductImageUrl: rs.App.Config.App.Service.Url.Api + storage.ResourcesProductImagePath + p.ProductId,
					Description:     p.Description,
					Subtypes:        []Subtype{},
				}
			}

			subtypeExists := false

			// Проверяем, существует ли подтип
			for _, existingSubtype := range productsMap[p.ProductName].Subtypes {
				if existingSubtype.SubtypeName == s.SubtypeName && existingSubtype.Type == s.Type {
					subtypeExists = true
					break
				}
			}

			// Если подтип не существует, добавляем его
			if !subtypeExists {
				productsMap[p.ProductName].Subtypes = append(productsMap[p.ProductName].Subtypes, Subtype{
					SubtypeName: s.SubtypeName,
					Type:        s.Type,
					Variants:    []Variant{},
				})
			}

			// Добавляем вариант в подтип
			for i, st := range productsMap[p.ProductName].Subtypes {
				if st.SubtypeName == s.SubtypeName && st.Type == s.Type {
					v.ServiceSvgUrl = rs.App.Config.App.Service.Url.Api + storage.ResourcesSvgFilePath + strings.ReplaceAll(v.Service, " ", "-")
					productsMap[p.ProductName].Subtypes[i].Variants = append(productsMap[p.ProductName].Subtypes[i].Variants, v)
					break
				}
			}
		}

		productsList := make([]*Product, 0, len(productsMap))
		for _, product := range productsMap {
			productsList = append(productsList, product)
		}

		// Вывод результатов
		//for _, product := range productsList {
		//	fmt.Printf("Product: %s\n", product.ProductName)
		//	for _, subtype := range product.Subtypes {
		//		fmt.Printf("  Subtype: %s\n", subtype.SubtypeName)
		//		for _, variant := range subtype.Variants {
		//			fmt.Printf("    Variant: %s\n", variant.VariantName)
		//		}
		//	}
		//}

		api_v1.RespondWithCreated(w, productsList)
	}
}
