package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"strconv"
	"strings"
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

func GetProductsForMainpage(ctx context.Context, pdb *Postgres, apiUrl, id, searchText, sort, sortType string) ([]Product, error) {
	productsMap := make(map[string]*Product)
	products := make([]Product, 0, len(productsMap))

	query := "SELECT type_name, subtype_name, service_name, product_name, variant_name, state_name, price, discount_money, discount_percent, final_price, item_name, mask, text_quantity, description, product_id, variant_id FROM product.product_variants_summary_all_data WHERE CONCAT(product_name, variant_name, tags, description) ILIKE ANY (ARRAY[$1])"
	if id != "" {
		query += " AND variant_id = '" + strings.ToLower(id) + "'"
	}
	query += getSort(0, sort, sortType, []string{"type_name", "subtype_name", "product_name", "variant_name", "price", "final_price", "discount_money", "discount_percent"})

	rows, err := pdb.Pool.Query(context.Background(), query, getTextWithPercents(searchText))
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
					continue
				} else if v.State == ProductStateUnavailableWithoutPrice {
					v.Price = 0
					v.DiscountPercent = 0
					v.DiscountMoney = 0
					v.FinalPrice = 0
					v.TextQuantity = ""
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

	for productName, product := range productsMap {
		for i := len(product.Subtypes) - 1; i >= 0; i-- {
			if len(product.Subtypes[i].Variants) == 0 {
				productsMap[productName].Subtypes = append(productsMap[productName].Subtypes[:i], productsMap[productName].Subtypes[i+1:]...)
			}
		}
	}

	var deleteKeys []string
	for productName, product := range productsMap {
		// If there are no subtypes left in the product, add the product's key to the deleteKeys slice
		if len(product.Subtypes) == 0 {
			deleteKeys = append(deleteKeys, productName)
		}
	}
	// Delete products with no subtypes from productsMap
	for _, key := range deleteKeys {
		delete(productsMap, key)
	}

	// Create the products slice from the remaining items in productsMap
	for _, product := range productsMap {
		products = append(products, *product)
	}

	return products, err
}

type AdminProducts struct {
	ProductId       string  `json:"product_id"`
	ProductImageUrl string  `json:"product_image_url"`
	ProductName     string  `json:"product_name"`
	Description     string  `json:"description"`
	Type            string  `json:"type_name"`
	SubtypeName     string  `json:"subtype_name"`
	VariantId       string  `json:"variant_id"`
	ServiceSvgUrl   string  `json:"service_svg_url"`
	VariantName     string  `json:"variant_name"`
	Service         string  `json:"service_name"`
	State           string  `json:"state_name"`
	Item            string  `json:"item_name"`
	Mask            string  `json:"mask"`
	TextQuantity    string  `json:"text_quantity"`
	QuantityCurrent int     `json:"quantity_current"`
	QuantitySold    int     `json:"quantity_sold"`
	Price           float64 `json:"price"`
	DiscountMoney   float64 `json:"discount_money"`
	DiscountPercent int     `json:"discount_percent"`
	FinalPrice      float64 `json:"final_price"`
}

func GetAdminVariants(ctx context.Context, pdb *Postgres, apiUrl, id, searchText, sort, sortType string) ([]AdminProducts, error) {
	var products []AdminProducts

	query := "SELECT product_id, product_name, description, type_name, subtype_name, variant_id, variant_name, service_name, state_name, item_name, mask, text_quantity, quantity_current, quantity_sold, price, discount_money, discount_percent, final_price FROM product.product_variants_summary_all_data WHERE CONCAT(product_name, variant_name, tags, description) ILIKE ANY (ARRAY[$1])"
	if id != "" {
		query += " AND variant_id = '" + strings.ToLower(id) + "'"
	}
	query += getSort(1, sort, sortType, []string{"CASE WHEN state_name = 'active' THEN 0 ELSE 1 END", "type_name", "subtype_name", "product_name", "variant_name", "price", "final_price", "discount_money", "discount_percent"})

	rows, err := pdb.Pool.Query(context.Background(), query, getTextWithPercents(searchText))
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var p AdminProducts
		err = rows.Scan(
			&p.ProductId,
			&p.ProductName,
			&p.Description,
			&p.Type,
			&p.SubtypeName,
			&p.VariantId,
			&p.VariantName,
			&p.Service,
			&p.State,
			&p.Item,
			&p.Mask,
			&p.TextQuantity,
			&p.QuantityCurrent,
			&p.QuantitySold,
			&p.Price,
			&p.DiscountMoney,
			&p.DiscountPercent,
			&p.FinalPrice,
		)
		if err != nil {
			return products, err
		}

		p.ProductImageUrl = GetProductImageUrl(apiUrl, p.ProductId)
		p.ServiceSvgUrl = GetSvgFileUrl(apiUrl, p.Service)

		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}

	return products, nil
}

type ProductItem struct {
	ItemNo     int     `json:"item_no"`
	ItemName   string  `json:"item_name"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
	Commentary *string `json:"commentary"`
}

func AdminGetItems(ctx context.Context, pdb *Postgres) ([]ProductItem, error) {
	var items []ProductItem

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT item_no, item_name, created_at, modified_at, commentary FROM product.item")
	if err != nil {
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item ProductItem
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&item.ItemNo,
			&item.ItemName,
			&createdAt,
			&modifiedAt,
			&item.Commentary,
		); err != nil {
			return items, err
		}
		item.CreatedAt = createdAt.Format(time.DateTime)
		item.ModifiedAt = modifiedAt.Format(time.DateTime)

		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return items, err
	}

	return items, err
}

type ProductService struct {
	ServiceNo   int     `json:"service_no"`
	ServiceName string  `json:"service_name"`
	ServiceUrl  string  `json:"service_url"`
	CreatedAt   string  `json:"created_at"`
	ModifiedAt  string  `json:"modified_at"`
	Commentary  *string `json:"commentary"`
}

func AdminGetServices(ctx context.Context, pdb *Postgres, apiUrl string) ([]ProductService, error) {
	var services []ProductService

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT service_no, service_name, created_at, modified_at, commentary FROM product.service ORDER BY service_name")
	if err != nil {
		return services, err
	}
	defer rows.Close()

	for rows.Next() {
		var service ProductService
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&service.ServiceNo,
			&service.ServiceName,
			&createdAt,
			&modifiedAt,
			&service.Commentary,
		); err != nil {
			return services, err
		}
		service.ServiceUrl = GetSvgFileUrl(apiUrl, service.ServiceName)
		service.CreatedAt = createdAt.Format(time.DateTime)
		service.ModifiedAt = modifiedAt.Format(time.DateTime)

		services = append(services, service)
	}
	if err = rows.Err(); err != nil {
		return services, err
	}

	return services, err
}

type ProductState struct {
	StateNo    int     `json:"state_no"`
	StateName  string  `json:"state_name"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
	Commentary *string `json:"commentary"`
}

func AdminGetStates(ctx context.Context, pdb *Postgres) ([]ProductState, error) {
	var states []ProductState

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT state_no, state_name, created_at, modified_at, commentary FROM product.state ORDER BY state_name")
	if err != nil {
		return states, err
	}
	defer rows.Close()

	for rows.Next() {
		var state ProductState
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&state.StateNo,
			&state.StateName,
			&createdAt,
			&modifiedAt,
			&state.Commentary,
		); err != nil {
			return states, err
		}

		state.CreatedAt = createdAt.Format(time.DateTime)
		state.ModifiedAt = modifiedAt.Format(time.DateTime)

		states = append(states, state)
	}
	if err = rows.Err(); err != nil {
		return states, err
	}

	return states, err
}

type ProductType struct {
	TypeNo     int     `json:"type_no"`
	TypeName   string  `json:"type_name"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
	Commentary *string `json:"commentary"`
}

func AdminGetTypes(ctx context.Context, pdb *Postgres) ([]ProductType, error) {
	var types []ProductType

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT type_no, type_name, created_at, modified_at, commentary FROM product.type ORDER BY type_name")
	if err != nil {
		return types, err
	}
	defer rows.Close()

	for rows.Next() {
		var typ ProductType
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&typ.TypeNo,
			&typ.TypeName,
			&createdAt,
			&modifiedAt,
			&typ.Commentary,
		); err != nil {
			return types, err
		}
		typ.CreatedAt = createdAt.Format(time.DateTime)
		typ.ModifiedAt = modifiedAt.Format(time.DateTime)

		types = append(types, typ)
	}
	if err = rows.Err(); err != nil {
		return types, err
	}

	return types, err
}

type ProductSubtype struct {
	SubtypeNo   int     `json:"subtype_no"`
	SubtypeName string  `json:"subtype_name"`
	CreatedAt   string  `json:"created_at"`
	ModifiedAt  string  `json:"modified_at"`
	Commentary  *string `json:"commentary"`
}

func AdminGetSubtypes(ctx context.Context, pdb *Postgres, typeName string) ([]ProductSubtype, error) {
	var subtypes []ProductSubtype

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT subtype_no, subtype_name, st.created_at, st.modified_at, st.commentary FROM product.subtype st join product.type t on st.type_no = t.type_no WHERE type_name  = $1 ORDER BY subtype_name",
		typeName)
	if err != nil {
		return subtypes, err
	}
	defer rows.Close()

	for rows.Next() {
		var subtype ProductSubtype
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&subtype.SubtypeNo,
			&subtype.SubtypeName,
			&createdAt,
			&modifiedAt,
			&subtype.Commentary,
		); err != nil {
			return subtypes, err
		}
		subtype.CreatedAt = createdAt.Format(time.DateTime)
		subtype.ModifiedAt = modifiedAt.Format(time.DateTime)

		subtypes = append(subtypes, subtype)
	}
	if err = rows.Err(); err != nil {
		return subtypes, err
	}

	return subtypes, err
}

type Product2 struct {
	ProductId   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Description string  `json:"description"`
	Tags        *string `json:"tags"`
	CreatedAt   string  `json:"created_at"`
	ModifiedAt  string  `json:"modified_at"`
	Commentary  *string `json:"commentary"`
}

func AdminGetProducts(ctx context.Context, pdb *Postgres) ([]Product2, error) {
	var products []Product2

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT product_id, product_name, description, tags, created_at, modified_at, commentary FROM product.product ORDER BY product_name")
	if err != nil {
		return products, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product2
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&product.ProductId,
			&product.ProductName,
			&product.Description,
			&product.Tags,
			&createdAt,
			&modifiedAt,
			&product.Commentary,
		); err != nil {
			return products, err
		}
		product.CreatedAt = createdAt.Format(time.DateTime)
		product.ModifiedAt = modifiedAt.Format(time.DateTime)

		products = append(products, product)
	}
	if err = rows.Err(); err != nil {
		return products, err
	}

	return products, err
}

func CreateAdminVariant(ctx context.Context, pdb *Postgres, productName, variantName, serviceName, stateName, subtypeName, itemName, mask, price, discountMoney, discountPercent, accountId string) error {
	var productId string
	var serviceId, stateId, subtypeId, itemId int
	var localDiscountMoney float64
	var localDiscountPercent int

	localDiscountMoney, _ = strconv.ParseFloat(discountMoney, 64)
	localDiscountPercent, _ = strconv.Atoi(discountPercent)

	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT product_id FROM product.product WHERE product_name = $1",
		productName).Scan(&productId); err != nil {
		return err
	}
	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT service_no FROM product.service WHERE service_name = $1",
		serviceName).Scan(&serviceId); err != nil {
		return err
	}
	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT state_no FROM product.state WHERE state_name = $1",
		stateName).Scan(&stateId); err != nil {
		return err
	}
	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT subtype_no FROM product.subtype WHERE subtype_name = $1",
		subtypeName).Scan(&subtypeId); err != nil {
		return err
	}
	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT item_no FROM product.item WHERE item_name = $1",
		itemName).Scan(&itemId); err != nil {
		return err
	}

	result, err := pdb.Pool.Exec(context.Background(),
		"INSERT INTO product.variant(product_id, variant_name, variant_service, variant_state, variant_subtype, variant_item, mask, price, discount_money, discount_percent, variant_account) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		productId, variantName, serviceId, stateId, subtypeId, itemId, mask, price, localDiscountMoney, localDiscountPercent, accountId)
	if err != nil {
		return err
	} else if result.RowsAffected() < 1 {
		return FailedInsert
	}

	UpdateData(ctx, pdb)
	return err
}

func UpdateData(ctx context.Context, pdb *Postgres) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"REFRESH MATERIALIZED VIEW product.product_variants_summary_all_data")

	return err
}

func UpdateAdminVariant(ctx context.Context, pdb *Postgres, id string, updateData map[string]interface{}) error {
	if len(updateData) == 0 {
		return errors.New("no data provided for update")
	}

	// Build SQL query
	query := "UPDATE product.variant SET"
	var args []interface{}
	i := 1
	for key, val := range updateData {
		query += fmt.Sprintf(" %s = $%d,", key, i)
		args = append(args, val)
		i++
	}

	query = query[:len(query)-1] + " WHERE variant_id = $" + strconv.Itoa(i)
	args = append(args, id)

	result, err := pdb.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	} else if result.RowsAffected() < 1 {
		return FailedUpdate
	}

	UpdateData(ctx, pdb)
	return nil
}

func AdminDeleteVariant(ctx context.Context, pdb *Postgres, variantId string) (bool, error) {
	var inUsage bool

	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM product.content WHERE content_variant = $1 AND content_order IS NOT NULL)",
		variantId).Scan(&inUsage); err != nil {
		return false, err
	}
	if inUsage {
		return true, nil
	}

	result, err := pdb.Pool.Exec(context.Background(),
		"DELETE FROM product.variant WHERE variant_id = $1",
		variantId)
	if err != nil {
		return false, err
	} else if result.RowsAffected() < 1 {
		return false, FailedDelete
	}

	UpdateData(ctx, pdb)
	return false, nil
}

func GetItemNo(ctx context.Context, pdb *Postgres, itemName string) (int, error) {
	var itemNo int

	err := pdb.Pool.QueryRow(ctx,
		"SELECT item_no FROM product.item WHERE item_name = $1",
		itemName).Scan(&itemNo)

	return itemNo, err
}

func GetStateNo(ctx context.Context, pdb *Postgres, stateName string) (int, error) {
	var stateNo int

	err := pdb.Pool.QueryRow(ctx,
		"SELECT state_no FROM product.state WHERE state_name = $1",
		stateName).Scan(&stateNo)

	return stateNo, err
}

func CreateAdminContent(ctx context.Context, pdb *Postgres, variantId string, data []string) error {
	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"UPDATE product.variant SET quantity_current = quantity_current + $1 WHERE variant_id = $2",
			len(data), variantId)
		if err != nil {
			return err
		} else if res.RowsAffected() < 1 {
			return FailedUpdate
		}

		for _, val := range data {
			res, err = tx.Exec(ctx,
				"INSERT INTO product.content(content_variant, data) VALUES ($1, $2)",
				variantId, val)
			if err != nil {
				return err
			} else if res.RowsAffected() < 1 {
				return FailedInsert
			}
		}

		return err
	})

	UpdateData(ctx, pdb)
	return err
}

func CreateOrder(ctx context.Context, pdb *Postgres, accountId, variantId string) (string, string, float64, error) {
	var orderId, variantName string
	var finalPrice float64

	err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx,
			"UPDATE product.variant SET quantity_current = quantity_current - 1 WHERE variant_id = $1 AND quantity_current > 0 AND variant_state = (SELECT state_no FROM product.state WHERE state_name = 'active') RETURNING variant_name",
			variantId).Scan(&variantName); err != nil {
			return err
		}

		if err := tx.QueryRow(ctx,
			"INSERT INTO product.order (order_account, price) SELECT $1, pp.final_price FROM product.product_variants_summary_all_data pp JOIN product.variant pv ON pv.variant_id = pp.variant_id WHERE pv.variant_id = $2 RETURNING order_id, price",
			accountId, variantId).Scan(&orderId, &finalPrice); err != nil {
			return err
		}

		result, err := tx.Exec(ctx,
			"UPDATE product.content SET content_order = $1 WHERE content_variant = $2 AND content_order IS NULL AND content_id = (SELECT content_id FROM product.content WHERE content_variant = $2 AND content_order IS NULL LIMIT 1)",
			orderId, variantId)
		if err != nil {
			return err
		} else if result.RowsAffected() < 1 {
			return FailedUpdate
		}

		return err
	})

	UpdateData(ctx, pdb)
	return orderId, variantName, finalPrice, err
}
