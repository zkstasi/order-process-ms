package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"order-ms/internal/model"
)

type Repo struct {
	db  *sql.DB
	ctx context.Context
}

func NewPostgresRepo(db *sql.DB) *Repo {
	return &Repo{db: db, ctx: context.Background()}
}

// Общий Save
func (r *Repo) Save(s model.Storable) error {
	switch v := s.(type) {
	case *model.Order:
		return r.SaveOrder(v)
	case *model.User:
		return r.SaveUser(v)
	case *model.Delivery:
		// опционально: INSERT в deliveries, если нужно
		return fmt.Errorf("Save Delivery: not implemented yet")
	case *model.Warehouse:
		// опционально: INSERT в warehouses, если нужно
		return fmt.Errorf("Save Warehouse: not implemented yet")
	default:
		return fmt.Errorf("unsupported type %T", s)
	}
}

// Заказы
func (r *Repo) SaveOrder(o *model.Order) error {
	_, err := r.db.ExecContext(r.ctx,
		`INSERT INTO orders (id, user_id, status, created_at)
		 VALUES ($1, $2, $3, $4)`,
		o.Id, o.UserID, int(o.Status), o.CreatedAt)
	return err
}

func (r *Repo) GetOrders() ([]*model.Order, error) {
	rows, err := r.db.QueryContext(r.ctx,
		`SELECT id, user_id, status, created_at
		   FROM orders
		   ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.Order
	for rows.Next() {
		var o model.Order
		var st int
		if err := rows.Scan(&o.Id, &o.UserID, &st, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.Status = model.OrderStatus(st)
		out = append(out, &o)
	}
	return out, rows.Err()
}

func (r *Repo) GetOrderByID(id string) (*model.Order, error) {
	var o model.Order
	var st int
	err := r.db.QueryRowContext(r.ctx,
		`SELECT id, user_id, status, created_at FROM orders WHERE id=$1`, id).
		Scan(&o.Id, &o.UserID, &st, &o.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	o.Status = model.OrderStatus(st)
	return &o, nil
}

func (r *Repo) DeleteOrder(id string) (bool, error) {
	res, err := r.db.ExecContext(r.ctx, `DELETE FROM orders WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repo) ConfirmOrder(orderId string) (bool, error) {
	// меняем 0->1
	res, err := r.db.ExecContext(r.ctx,
		`UPDATE orders SET status=$1 WHERE id=$2 AND status=$3`,
		int(model.OrderConfirmed), orderId, int(model.OrderCreated))
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repo) DeliverOrder(orderId string) (bool, error) {
	// меняем 1->2
	res, err := r.db.ExecContext(r.ctx,
		`UPDATE orders SET status=$1 WHERE id=$2 AND status=$3`,
		int(model.OrderDelivered), orderId, int(model.OrderConfirmed))
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repo) CancelOrder(orderId string) (bool, error) {
	// меняем 0/1 -> 3
	res, err := r.db.ExecContext(r.ctx,
		`UPDATE orders SET status=$1 WHERE id=$2 AND status IN ($3,$4)`,
		int(model.OrderCancelled), orderId,
		int(model.OrderCreated), int(model.OrderConfirmed))
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// Пользователи
func (r *Repo) SaveUser(u *model.User) error {
	_, err := r.db.ExecContext(r.ctx,
		`INSERT INTO users (id, name) VALUES ($1,$2)`,
		u.Id, u.Name)
	return err
}

func (r *Repo) GetUsers() ([]*model.User, error) {
	rows, err := r.db.QueryContext(r.ctx, `SELECT id, name FROM users ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.Id, &u.Name); err != nil {
			return nil, err
		}
		out = append(out, &u)
	}
	return out, rows.Err()
}

func (r *Repo) GetUserByID(id string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRowContext(r.ctx, `SELECT id, name FROM users WHERE id=$1`, id).
		Scan(&u.Id, &u.Name)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) UpdateUserName(id, name string) (bool, error) {
	res, err := r.db.ExecContext(r.ctx, `UPDATE users SET name=$1 WHERE id=$2`, name, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (r *Repo) DeleteUser(id string) (bool, error) {
	res, err := r.db.ExecContext(r.ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// Доставки и склады (минимальные заглушки, чтобы удовлетворить интерфейсу)
func (r *Repo) GetDeliveries() ([]*model.Delivery, error)  { return []*model.Delivery{}, nil }
func (r *Repo) GetWarehouses() ([]*model.Warehouse, error) { return []*model.Warehouse{}, nil }
