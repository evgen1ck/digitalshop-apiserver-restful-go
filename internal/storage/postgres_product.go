package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
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

func GetAdminVariants(ctx context.Context, pdb *Postgres, apiUrl, id, searchText, sort, sortType, activeFirst string) ([]AdminProducts, error) {
	var products []AdminProducts

	query := "SELECT product_id, product_name, description, type_name, subtype_name, variant_id, variant_name, service_name, state_name, item_name, mask, text_quantity, quantity_current, quantity_sold, price, discount_money, discount_percent, final_price FROM product.product_variants_summary_all_data WHERE CONCAT(product_name, variant_name, tags, description) ILIKE ANY (ARRAY[$1])"
	if id != "" {
		query += " AND variant_id = '" + strings.ToLower(id) + "'"
	}

	var ind int
	if activeFirst == "false" {
		ind = 0
	} else {
		ind = 1
	}
	query += getSort(ind, sort, sortType, []string{"CASE WHEN state_name = 'active' THEN 0 ELSE 1 END", "type_name", "subtype_name", "product_name", "variant_name", "price", "final_price", "discount_money", "discount_percent", "quantity_current"})

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

func AdminGetServices(ctx context.Context, pdb *Postgres, apiUrl, serviceName string) ([]ProductService, error) {
	var services []ProductService

	query := "SELECT service_no, service_name, created_at, modified_at, commentary FROM product.service"
	if serviceName != "" {
		query += " WHERE service_name = '" + serviceName + "'"
	}
	query += " ORDER BY service_name"

	rows, err := pdb.Pool.Query(context.Background(),
		query)
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

func AdminGetTypes(ctx context.Context, pdb *Postgres, typeName string) ([]ProductType, error) {
	var types []ProductType

	query := "SELECT type_no, type_name, created_at, modified_at, commentary FROM product.type"
	if typeName != "" {
		query += " WHERE type_name = '" + typeName + "'"
	}
	query += " ORDER BY type_name"

	rows, err := pdb.Pool.Query(context.Background(),
		query)
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

func AdminGetSubtypes(ctx context.Context, pdb *Postgres, typeName, subtypeName string) ([]ProductSubtype, error) {
	var subtypes []ProductSubtype

	query := "SELECT subtype_no, subtype_name, st.created_at, st.modified_at, st.commentary FROM product.subtype st join product.type t on st.type_no = t.type_no  "
	if subtypeName != "" && typeName != "" {
		query += " WHERE type_name = '" + typeName + "' AND subtype_name = '" + subtypeName + "'"
	} else if subtypeName != "" {
		query += " WHERE subtype_name = '" + subtypeName + "'"
	} else if typeName != "" {
		query += " WHERE type_name = '" + typeName + "'"
	}
	query += " ORDER BY subtype_name"

	rows, err := pdb.Pool.Query(context.Background(),
		query)
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

func AdminDeleteVariant(ctx context.Context, pdb *Postgres, variantId string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"DELETE FROM product.variant WHERE variant_id = $1",
		variantId)

	UpdateData(ctx, pdb)
	return err
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

type GetAdminContentsData struct {
	ContentId  string  `json:"content_id"`
	Data       string  `json:"data"`
	CreatedAt  string  `json:"created_at"`
	ModifiedAt string  `json:"modified_at"`
	Commentary *string `json:"commentary"`
}

func GetAdminContents(ctx context.Context, pdb *Postgres, id string) ([]GetAdminContentsData, error) {
	var contents []GetAdminContentsData

	rows, err := pdb.Pool.Query(context.Background(),
		"SELECT content_id, data, created_at, modified_at, commentary FROM product.content WHERE content_variant = $1 ORDER BY created_at DESC",
		id)
	if err != nil {
		return contents, err
	}
	defer rows.Close()

	for rows.Next() {
		var content GetAdminContentsData
		var createdAt, modifiedAt time.Time

		if err = rows.Scan(
			&content.ContentId,
			&content.Data,
			&createdAt,
			&modifiedAt,
			&content.Commentary,
		); err != nil {
			return contents, err
		}

		content.CreatedAt = createdAt.Format(time.DateTime)
		content.ModifiedAt = modifiedAt.Format(time.DateTime)

		contents = append(contents, content)
	}
	if err = rows.Err(); err != nil {
		return contents, err
	}

	return contents, err
}

func DeleteAdminContent(ctx context.Context, pdb *Postgres, id string) error {
	if err := execInTx(ctx, pdb.Pool, func(tx pgx.Tx) error {
		res, err := tx.Exec(ctx,
			"UPDATE product.variant SET quantity_current = quantity_current - 1 WHERE variant_id = (SELECT content_variant FROM product.content WHERE content_id = $1)",
			id)
		if err != nil {
			return err
		} else if res.RowsAffected() < 1 {
			return FailedUpdate
		}

		res, err = tx.Exec(context.Background(),
			"DELETE FROM product.content WHERE content_id = $1 AND content_order IS NULL",
			id)
		if err != nil {
			return err
		} else if res.RowsAffected() < 1 {
			return FailedDelete
		}

		return err
	}); err != nil {
		return err
	}

	return UpdateData(ctx, pdb)
}

func DeleteAdminType(ctx context.Context, pdb *Postgres, id string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"DELETE FROM product.type WHERE type_name = $1",
		id)

	UpdateData(ctx, pdb)
	return err
}

func DeleteAdminSubtype(ctx context.Context, pdb *Postgres, id string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"DELETE FROM product.subtype WHERE subtype_name = $1",
		id)

	UpdateData(ctx, pdb)
	return err
}

func DeleteAdminService(ctx context.Context, pdb *Postgres, id string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"DELETE FROM product.service WHERE service_name = $1",
		id)

	UpdateData(ctx, pdb)
	return err
}

func DeleteAdminProduct(ctx context.Context, pdb *Postgres, id string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"DELETE FROM product.product WHERE product_name = $1",
		id)

	UpdateData(ctx, pdb)
	return err
}

func CreateAdminType(ctx context.Context, pdb *Postgres, name string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"INSERT INTO product.type(type_name) VALUES ($1)",
		name)

	UpdateData(ctx, pdb)
	return err
}

func CreateAdminSubtype(ctx context.Context, pdb *Postgres, name, name2 string) error {
	var typeNo int

	if err := pdb.Pool.QueryRow(context.Background(),
		"SELECT type_no FROM product.type WHERE type_name = $1",
		name2).Scan(&typeNo); err != nil {
		return err
	}

	_, err := pdb.Pool.Exec(context.Background(),
		"INSERT INTO product.subtype(type_no, subtype_name) VALUES ($1, $2)",
		typeNo, name)

	UpdateData(ctx, pdb)
	return err
}

func CreateAdminService(ctx context.Context, pdb *Postgres, name string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"INSERT INTO product.service(service_name) VALUES ($1)",
		name)

	UpdateData(ctx, pdb)
	return err
}

func CreateAdminProduct(ctx context.Context, pdb *Postgres, productName, tags, description string) (string, error) {
	var uuid string

	if err := pdb.Pool.QueryRow(context.Background(),
		"INSERT INTO product.product(product_name, tags, description) VALUES ($1, $2, $3) RETURNING product_id",
		productName, tags, description).Scan(&uuid); err != nil {
		return "", err
	}

	UpdateData(ctx, pdb)
	return uuid, nil
}

func EditAdminType(ctx context.Context, pdb *Postgres, name, newName string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"UPDATE product.type SET type_name = $1 WHERE type_name = $2",
		newName, name)

	UpdateData(ctx, pdb)
	return err
}

func EditAdminSubtype(ctx context.Context, pdb *Postgres, name, newName string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"UPDATE product.subtype SET subtype_name = $1 WHERE subtype_name = $2",
		newName, name)

	UpdateData(ctx, pdb)
	return err
}

func EditAdminService(ctx context.Context, pdb *Postgres, name, newName string) error {
	_, err := pdb.Pool.Exec(context.Background(),
		"UPDATE product.service SET service_name = $1 WHERE service_name = $2",
		newName, name)

	UpdateData(ctx, pdb)
	return err
}
