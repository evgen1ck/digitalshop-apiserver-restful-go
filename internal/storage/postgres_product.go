package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"time"
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

func GetProductsForMainpage(ctx context.Context, pdb *Postgres, apiUrl string) ([]Product, error) {
	productsMap := make(map[string]*Product)
	products := make([]Product, 0, len(productsMap))

	rows, err := pdb.Pool.Query(ctx,
		"SELECT type_name, subtype_name, service_name, product_name, variant_name, state_name, price, discount_money, discount_percent, final_price, item_name, mask, text_quantity, description, product_id, variant_id FROM product.product_variants_summary_for_mainpage")
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var v Variant
		var s Subtype
		var p Product

		if err = rows.Scan(
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
		); err != nil {
			return products, err
		}

		if _, ok := productsMap[p.ProductName]; !ok {
			productsMap[p.ProductName] = &Product{
				ProductName:     p.ProductName,
				ProductId:       p.ProductId,
				ProductImageUrl: GetProductImageUrl(apiUrl, p.ProductId),
				Description:     p.Description,
				Subtypes:        []Subtype{},
			}
		}

		var subtypeExists bool

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
				if v.State == ProductStateInvisible || v.State == ProductStateDeleted {
					break
				} else if v.State == ProductStateUnavailableWithoutPrice {
					v.Price = 0
					v.DiscountPercent = 0
					v.DiscountMoney = 0
					v.FinalPrice = 0
				}
				v.ServiceSvgUrl = GetSvgFileUrl(apiUrl, v.Service)
				productsMap[p.ProductName].Subtypes[i].Variants = append(productsMap[p.ProductName].Subtypes[i].Variants, v)
				break
			}
		}
	}
	if err = rows.Err(); err != nil {
		return products, err
	}

	for _, product := range productsMap {
		products = append(products, *product)
	}

	return products, err
}

func GetProductsWithParams(ctx context.Context, pdb *Postgres, query1, query2, query3 string) (pgx.Rows, error) {
	rows, err := pdb.Pool.Query(ctx,
		"SELECT type_name, subtype_name, service_name, product_name, variant_name, state_name, price, discount_money, discount_percent, final_price, item_name, mask, text_quantity, description, product_id, variant_id FROM product.product_variants_summary_all_data WHERE concat(product_name, variant_name, description, tags) ILIKE ANY (ARRAY[$1, $2, $3])",
		"%"+query1+"%", "%"+query2+"%", "%"+query3+"%")

	return rows, err
}

type ProductItem struct {
	ItemNo     int       `json:"item_no"`
	ItemName   string    `json:"item_name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Commentary string    `json:"commentary"`
}

func GetProductItems(ctx context.Context, pdb *Postgres) ([]ProductItem, error) {
	var items []ProductItem

	rows, err := pdb.Pool.Query(ctx,
		"SELECT item_no, item_name, created_at, modified_at, commentary FROM product.item")
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item ProductItem

		if err = rows.Scan(
			&item.ItemNo,
			&item.ItemName,
			&item.CreatedAt,
			&item.ModifiedAt,
			&item.Commentary,
		); err != nil {
			return items, err
		}

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return items, err
	}

	return items, err
}

type ProductService struct {
	ServiceNo   int       `json:"service_no"`
	ServiceName string    `json:"service_name"`
	ServiceUrl  string    `json:"service_url"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	Commentary  *string   `json:"commentary"`
}

func GetProductServices(ctx context.Context, pdb *Postgres, apiUrl string) ([]ProductService, error) {
	var services []ProductService

	rows, err := pdb.Pool.Query(ctx,
		"SELECT service_no, service_name, created_at, modified_at, commentary FROM product.service")
	if err != nil {
		return services, err
	}
	defer rows.Close()

	for rows.Next() {
		var service ProductService

		if err = rows.Scan(
			&service.ServiceNo,
			&service.ServiceName,
			&service.CreatedAt,
			&service.ModifiedAt,
			&service.Commentary,
		); err != nil {
			return services, err
		}
		service.ServiceUrl = GetSvgFileUrl(apiUrl, service.ServiceName)

		services = append(services, service)
	}
	if err = rows.Err(); err != nil {
		return services, err
	}

	return services, err
}

type ProductState struct {
	StateNo    int       `json:"state_no"`
	StateName  string    `json:"state_name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Commentary string    `json:"commentary"`
}

func GetProductStates(ctx context.Context, pdb *Postgres) ([]ProductState, error) {
	var states []ProductState

	rows, err := pdb.Pool.Query(ctx,
		"SELECT state_no, state_name, created_at, modified_at, commentary FROM product.state")
	if err != nil {
		return states, err
	}
	defer rows.Close()

	for rows.Next() {
		var state ProductState

		if err = rows.Scan(
			&state.StateNo,
			&state.StateName,
			&state.CreatedAt,
			&state.ModifiedAt,
			&state.Commentary,
		); err != nil {
			return states, err
		}

		states = append(states, state)
	}
	if err = rows.Err(); err != nil {
		return states, err
	}

	return states, err
}

type ProductType struct {
	TypeNo     int       `json:"type_no"`
	TypeName   string    `json:"type_name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Commentary string    `json:"commentary"`
}

func GetProductTypes(ctx context.Context, pdb *Postgres) ([]ProductType, error) {
	var types []ProductType

	rows, err := pdb.Pool.Query(ctx,
		"SELECT state_no, state_name, created_at, modified_at, commentary FROM product.state")
	if err != nil {
		return types, err
	}
	defer rows.Close()

	for rows.Next() {
		var typ ProductType

		if err = rows.Scan(
			&typ.TypeNo,
			&typ.TypeName,
			&typ.CreatedAt,
			&typ.ModifiedAt,
			&typ.Commentary,
		); err != nil {
			return types, err
		}

		types = append(types, typ)
	}
	if err = rows.Err(); err != nil {
		return types, err
	}

	return types, err
}

type ProductSubtype struct {
	SubtypeNo   int       `json:"subtype_no"`
	SubtypeName string    `json:"subtype_name"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	Commentary  string    `json:"commentary"`
}

func GetProductSubtypes(ctx context.Context, pdb *Postgres, typeNo string) ([]ProductSubtype, error) {
	var subtypes []ProductSubtype

	rows, err := pdb.Pool.Query(ctx,
		"SELECT subtype_no, subtype_name, created_at, modified_at, commentary FROM product.subtype WHERE type_no = $1",
		typeNo)
	if err != nil {
		return subtypes, err
	}
	defer rows.Close()

	for rows.Next() {
		var subtype ProductSubtype

		if err = rows.Scan(
			&subtype.SubtypeNo,
			&subtype.SubtypeName,
			&subtype.CreatedAt,
			&subtype.ModifiedAt,
			&subtype.Commentary,
		); err != nil {
			return subtypes, err
		}

		subtypes = append(subtypes, subtype)
	}
	if err = rows.Err(); err != nil {
		return subtypes, err
	}

	return subtypes, err
}

func GetProductVariantForPayment(ctx context.Context, pdb *Postgres, variantId string) (string, string, string, int, float64, error) {
	var variantName, variantState string
	var finalPrice float64
	var quantityCurrent int
	var productId uuid.UUID

	err := pdb.Pool.QueryRow(ctx,
		"SELECT product_id, variant_name, state_name, quantity_current, final_price FROM product.product_variants_summary_all_data WHERE variant_id = $1",
		variantId).Scan(&productId, &variantName, &variantState, &quantityCurrent, &finalPrice)
	if err != nil {
		return productId.String(), variantName, variantState, quantityCurrent, finalPrice, err
	}

	return productId.String(), variantName, variantState, quantityCurrent, finalPrice, err
}
