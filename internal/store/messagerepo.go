// Инкапсуляция БД
package store

import (
	"WB_Tech_level_0/internal/models"
	"database/sql"
	"fmt"
)

type MessageRepository struct {
	store *Store
}

// Метод Create - создание записи в БД
func (r *MessageRepository) Create(model *models.Data) error {
	if err := r.insertOrders(&model.Order); err != nil {
		return err
	}
	if err := r.insertPayment(&model.Payment); err != nil {
		return err
	}
	if err := r.insertDelivery(&model.Delivery, &model.Order.Order_uid); err != nil {
		return err
	}
	if err := r.insertItems(model.Items); err != nil {
		return err
	}
	return nil
}

// Вставка в таблицу
func (r *MessageRepository) insertDelivery(model *models.Delivery, key *string) error {
	_, err := r.store.db.Exec(`
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, key, model.Name, model.Phone, model.Zip, model.City,
		model.Address, model.Region, model.Email,
	)
	if err != nil {
		return err
	}
	return nil
}

// Вставка в таблицу
func (r *MessageRepository) insertPayment(model *models.Payment) error {
	_, err := r.store.db.Exec(`
		INSERT INTO payment (transaction, request_id, currency, provider, 
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`, model.Transaction, model.Request_id, model.Currency, model.Provider, model.Amount,
		model.Payment_dt, model.Bank, model.Delivery_cost, model.Goods_total, model.Custom_fee,
	)
	if err != nil {
		return err
	}
	return nil
}

// Вставка в таблицу
func (r *MessageRepository) insertItems(model []models.Items) error {
	for _, item := range model {
		_, err := r.store.db.Exec(`
		INSERT INTO items (chrt_id, track_number, price, rid, 
			name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`, item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.Total_price, item.Nm_id, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Вставка в таблицу
func (r *MessageRepository) insertOrders(model *models.Orders) error {
	_, err := r.store.db.Exec(`
	INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, model.Order_uid, model.Track_number, model.Entry, model.Locale, model.Internal_signature,
		model.Customer_id, model.Delivery_service, model.Shardkey, model.Sm_id, model.Date_created, model.Oof_shard,
	)
	if err != nil {
		return err
	}
	return nil
}

// Метод Read_id - поиск записи по order_uid
func (r *MessageRepository) Read_id(id string) (models.Data, error) {
	var data models.Data
	var err error
	data.Order, err = r.readOrders(id)
	if err == sql.ErrNoRows {
		return models.Data{}, fmt.Errorf("запись не сущесвтует")
	} else if err != nil {
		return models.Data{}, err
	}
	data.Delivery, err = r.readDelivery(data.Order.Order_uid)
	if err == sql.ErrNoRows {
		return models.Data{}, fmt.Errorf("запись не сущесвтует")
	} else if err != nil {
		return models.Data{}, err
	}
	data.Payment, err = r.readPayment(data.Order.Order_uid)
	if err == sql.ErrNoRows {
		return models.Data{}, fmt.Errorf("запись не сущесвтует")
	} else if err != nil {
		return models.Data{}, err
	}
	data.Items, err = r.readItems(data.Order.Track_number)
	if err != nil {
		return models.Data{}, err
	}
	return data, nil
}

// Поиск в таблице
func (r *MessageRepository) readOrders(id string) (models.Orders, error) {
	row := r.store.db.QueryRow("SELECT * FROM orders WHERE order_uid = $1", id)
	var order models.Orders
	err := row.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale,
		&order.Internal_signature, &order.Customer_id, &order.Delivery_service, &order.Shardkey,
		&order.Sm_id, &order.Date_created, &order.Oof_shard)
	if err != nil {
		return models.Orders{}, err
	}
	return order, nil
}

// Поиск в таблице
func (r *MessageRepository) readDelivery(id string) (models.Delivery, error) {
	row := r.store.db.QueryRow("SELECT * FROM delivery WHERE order_uid = $1", id)
	var order models.Delivery
	err := row.Scan(&order.Order_uid, &order.Name, &order.Phone, &order.Zip,
		&order.City, &order.Address, &order.Region, &order.Email)
	if err != nil {
		return models.Delivery{}, err
	}
	return order, nil
}

// Поиск в таблице
func (r *MessageRepository) readItems(id string) ([]models.Items, error) {
	var order []models.Items

	rows, err := r.store.db.Query("SELECT * FROM items WHERE track_number = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.Items
		var id int
		err := rows.Scan(&id, &item.Chrt_id, &item.Track_number, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.Total_price, &item.Nm_id, &item.Brand, &item.Status)
		if err != nil {
			return nil, nil
		}
		order = append(order, item)
	}
	if err := rows.Err(); err != nil {
		return nil, nil
	}
	if len(order) == 0 {
		return nil, fmt.Errorf("записи %s не найдены", id)
	}
	return order, nil
}

// Поиск в таблице
func (r *MessageRepository) readPayment(id string) (models.Payment, error) {
	row := r.store.db.QueryRow("SELECT * FROM payment WHERE transaction = $1", id)
	var order models.Payment
	err := row.Scan(&order.Transaction, &order.Request_id, &order.Currency, &order.Provider,
		&order.Amount, &order.Payment_dt, &order.Bank, &order.Delivery_cost,
		&order.Goods_total, &order.Custom_fee)
	if err != nil {
		return models.Payment{}, err
	}
	return order, nil
}

// Метод Read_ALL - получение всех данных из БД, используется для восстановления КЭШа
func (r *MessageRepository) Read_ALL() ([]models.Data, error) {
	var data []models.Data
	orders, err := r.readOrder()
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		var data_buf models.Data
		data_buf.Order = order
		data_buf.Delivery, err = r.readDelivery(order.Order_uid)
		if err != nil {
			return nil, err
		}
		data_buf.Payment, err = r.readPayment(order.Order_uid)
		if err != nil {
			return nil, err
		}
		data_buf.Items, err = r.readItems(order.Track_number)
		if err != nil {
			return nil, err
		}
		data = append(data, data_buf)
	}
	return data, nil
}

// Поиск в таблице
func (r *MessageRepository) readOrder() ([]models.Orders, error) {
	var orders []models.Orders
	rows, err := r.store.db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Orders
		err := rows.Scan(&order.Order_uid, &order.Track_number, &order.Entry, &order.Locale, &order.Internal_signature,
			&order.Customer_id, &order.Delivery_service, &order.Shardkey, &order.Sm_id, &order.Date_created, &order.Oof_shard)
		if err != nil {
			return nil, nil
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, nil
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("записи не найдены")
	}
	return orders, nil
}
