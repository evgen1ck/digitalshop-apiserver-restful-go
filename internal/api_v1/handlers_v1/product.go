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
	VariantId       string  `json:"variant_id"`
	ServiceSvgUrl   string  `json:"service_svg_url"`
	VariantName     string  `json:"variant_name"`
	Service         string  `json:"service"`
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
	Type        string    `json:"type"`
	SubtypeName string    `json:"subtype_name"`
	Variants    []Variant `json:"variants"`
}

type Product struct {
	ProductId       string    `json:"product_id"`
	ProductImageUrl string    `json:"product_image_url"`
	ProductName     string    `json:"product_name"`
	Description     string    `json:"description"`
	Subtypes        []Subtype `json:"subtypes"`
}

func (rs *Resolver) ProductsDataForMainpage(w http.ResponseWriter, r *http.Request) {
	rows, err := storage.GetProductsForMainpage(r.Context(), rs.App.Postgres)
	if err != nil {
		rs.App.Logger.NewWarn("error in get products for mainpage", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	defer rows.Close()

	productsMap := make(map[string]*Product)
	for rows.Next() {
		var v Variant
		var s Subtype
		var p Product

		err = rows.Scan(
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

		// Checking if a subtype exists
		for _, existingSubtype := range productsMap[p.ProductName].Subtypes {
			if existingSubtype.SubtypeName == s.SubtypeName && existingSubtype.Type == s.Type {
				subtypeExists = true
				break
			}
		}

		// If the subtype does not exist, add it
		if !subtypeExists {
			productsMap[p.ProductName].Subtypes = append(productsMap[p.ProductName].Subtypes, Subtype{
				SubtypeName: s.SubtypeName,
				Type:        s.Type,
				Variants:    []Variant{},
			})
		}

		// Adding a variant to the subtype
		for i, st := range productsMap[p.ProductName].Subtypes {
			if st.SubtypeName == s.SubtypeName && st.Type == s.Type {
				if v.State == "invisible" || v.State == "deleted" {
					break
				} else if v.State == "unavailable without price" {
					v.Price = 0
					v.DiscountPercent = 0
					v.DiscountMoney = 0
					v.FinalPrice = 0
				}
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

	api_v1.RespondWithCreated(w, productsList)
}

func (rs *Resolver) ProductsData(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	transliterate := tl.Transliterate(search)
	rusToEng := tl.RusToEng(search)

	rows, err := storage.GetProductsWithParams(r.Context(), rs.App.Postgres, search, transliterate, rusToEng)
	if err != nil {
		rs.App.Logger.NewWarn("error in get products with params", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	defer rows.Close()

}
