package order_pgx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/playture/backend/internal/entity"
	"github.com/playture/backend/internal/infrastructure/postgresql"
	orderRepository "github.com/playture/backend/internal/repository/order_repository"
	"github.com/playture/backend/utils"
)

const (
	createOrder = `
		INSERT INTO orders (
			id, job_id, user_email, user_name, stripe_payment_intent_id, stripe_customer_id,
			amount, currency, payment_status, paid_at, order_type, requirements,
			production_job_id, production_status, delivery_method, delivered_at,
			customer_notes, support_ticket_id, ip_address, user_agent, expires_at,
			created_at, updated_at
		) VALUES (
			$1,$2,$3,$4,$5,$6,
			$7,$8,$9,$10,$11,$12,
			$13,$14,$15,$16,
			$17,$18,$19,$20,$21,
			$22,$23
		) RETURNING id`

	deleteOrder = "DELETE FROM orders WHERE id = $1"

	updateQuery = `
		UPDATE orders
		SET job_id=$2, user_email=$3, user_name=$4, stripe_payment_intent_id=$5, stripe_customer_id=$6,
			amount=$7, currency=$8, payment_status=$9, paid_at=$10, order_type=$11, requirements=$12,
			production_job_id=$13, production_status=$14, delivery_method=$15, delivered_at=$16,
			customer_notes=$17, support_ticket_id=$18, ip_address=$19, user_agent=$20, expires_at=$21,
			updated_at=$22
		WHERE id=$1`
)

type OrderPgx struct {
	logger   *slog.Logger
	postgres *postgresql.Postgres
}

func NewOrderPgx(
	logger *slog.Logger,
	postgres *postgresql.Postgres,
) *OrderPgx {
	return &OrderPgx{
		logger:   logger.With("layer", "OrderRepository"),
		postgres: postgres,
	}
}

func (o *OrderPgx) Create(
	ctx context.Context,
	order *entity.Order,
	tx pgx.Tx,
) (string, error) {
	lg := o.logger.With("method", "Create")

	var id string
	err := tx.QueryRow(ctx, createOrder,
		order.ID, order.JobID, order.UserEmail, order.UserName, order.StripePaymentIntentID, order.StripeCustomerID,
		order.Amount, order.Currency, order.PaymentStatus, order.PaidAt, order.OrderType, order.Requirements,
		order.ProductionJobID, order.ProductionStatus, order.DeliveryMethod, order.DeliveredAt,
		order.CustomerNotes, order.SupportTicketID, order.IPAddress, order.UserAgent, order.ExpiresAt,
		order.CreatedAt, order.UpdatedAt,
	).Scan(&id)

	if err != nil {
		lg.Error("failed to insert order", "err", err)
		return "", utils.WrapError("insert order", err)
	}

	return id, nil
}

func (o *OrderPgx) FindByField(ctx context.Context, field string, value interface{}, tx pgx.Tx) (*entity.Order, error) {
	lg := o.logger.With("method", "FindByField")

	query := fmt.Sprintf("SELECT * FROM orders WHERE %s = $1 LIMIT 1", field)

	row := tx.QueryRow(ctx, query, value)
	var order entity.Order
	if err := row.Scan(
		&order.ID, &order.JobID, &order.UserEmail, &order.UserName, &order.StripePaymentIntentID, &order.StripeCustomerID,
		&order.Amount, &order.Currency, &order.PaymentStatus, &order.PaidAt, &order.OrderType, &order.Requirements,
		&order.ProductionJobID, &order.ProductionStatus, &order.DeliveryMethod, &order.DeliveredAt,
		&order.CustomerNotes, &order.SupportTicketID, &order.IPAddress, &order.UserAgent, &order.ExpiresAt,
		&order.CreatedAt, &order.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, orderRepository.ErrOrderNotFound
		}
		lg.Error("failed to scan order", "err", err)
		return nil, utils.WrapError("find order by field", err)
	}

	return &order, nil
}

func (o *OrderPgx) List(ctx context.Context, paymentStatus *entity.PaymentStatus, productionStatus *entity.ProductionStatus, limit, offset int, tx pgx.Tx) ([]*entity.Order, error) {
	lg := o.logger.With("method", "List")

	query := "SELECT * FROM orders WHERE 1=1"
	args := []interface{}{}
	argIdx := 1

	if paymentStatus != nil {
		query += fmt.Sprintf(" AND payment_status = $%d", argIdx)
		args = append(args, *paymentStatus)
		argIdx++
	}

	if productionStatus != nil {
		query += fmt.Sprintf(" AND production_status = $%d", argIdx)
		args = append(args, *productionStatus)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		lg.Error("failed to list orders", "err", err)
		return nil, utils.WrapError("list orders", err)
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(
			&order.ID, &order.JobID, &order.UserEmail, &order.UserName, &order.StripePaymentIntentID, &order.StripeCustomerID,
			&order.Amount, &order.Currency, &order.PaymentStatus, &order.PaidAt, &order.OrderType, &order.Requirements,
			&order.ProductionJobID, &order.ProductionStatus, &order.DeliveryMethod, &order.DeliveredAt,
			&order.CustomerNotes, &order.SupportTicketID, &order.IPAddress, &order.UserAgent, &order.ExpiresAt,
			&order.CreatedAt, &order.UpdatedAt,
		); err != nil {
			lg.Error("failed to scan order", "err", err)
			return nil, utils.WrapError("scan order row", err)
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (o *OrderPgx) Delete(ctx context.Context, id string, tx pgx.Tx) error {
	lg := o.logger.With("method", "Delete")

	cmd, err := tx.Exec(ctx, deleteOrder, id)
	if err != nil {
		lg.Error("failed to delete order", "id", id, "err", err)
		return utils.WrapError("delete order", err)
	}
	if cmd.RowsAffected() == 0 {
		return orderRepository.ErrOrderNotFound
	}
	return nil
}

func (o *OrderPgx) Update(ctx context.Context, order *entity.Order, tx pgx.Tx) error {
	lg := o.logger.With("method", "Update")

	query := updateQuery

	cmd, err := tx.Exec(ctx, query,
		order.ID, order.JobID, order.UserEmail, order.UserName, order.StripePaymentIntentID, order.StripeCustomerID,
		order.Amount, order.Currency, order.PaymentStatus, order.PaidAt, order.OrderType, order.Requirements,
		order.ProductionJobID, order.ProductionStatus, order.DeliveryMethod, order.DeliveredAt,
		order.CustomerNotes, order.SupportTicketID, order.IPAddress, order.UserAgent, order.ExpiresAt,
		order.UpdatedAt,
	)
	if err != nil {
		lg.Error("failed to update order", "id", order.ID, "err", err)
		return utils.WrapError("update order", err)
	}
	if cmd.RowsAffected() == 0 {
		return orderRepository.ErrOrderNotFound
	}

	return nil
}
