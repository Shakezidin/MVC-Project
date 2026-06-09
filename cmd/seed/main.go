package main

import (
	"context"
	"fmt"
	"os"

	"github.com/banking/bank-server/internal/auth"
	"github.com/banking/bank-server/internal/config"
	"github.com/banking/bank-server/internal/repository"
	"github.com/jackc/pgx/v5"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "seed error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("seed complete")
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	db, err := repository.NewDatabase(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		return err
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queries := []struct {
		name  string
		query string
		args  []interface{}
	}{
		{
			"users",
			`INSERT INTO users (id, email, password_hash, full_name) VALUES
				('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'john.doe@example.com', $1, 'John Doe'),
				('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'jane.smith@example.com', $1, 'Jane Smith')
			ON CONFLICT (email) DO NOTHING`,
			[]interface{}{passwordHash},
		},
		{
			"bank_accounts",
			`INSERT INTO bank_accounts (id, user_id, account_type, branch_name, account_number, status) VALUES
				('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'SAVINGS', 'Mumbai Main Branch', '1234567890123456', 'ACTIVE'),
				('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'CURRENT', 'Mumbai Main Branch', '9876543210987654', 'ACTIVE'),
				('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'SALARY', 'Delhi Central Branch', '5555666677778888', 'ACTIVE')
			ON CONFLICT (account_number) DO NOTHING`,
			nil,
		},
		{
			"balances",
			`INSERT INTO balances (account_id, available_balance, current_balance, currency) VALUES
				('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 125000.50, 125000.50, 'INR'),
				('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 500000.00, 500000.00, 'INR'),
				('e4eebc99-9c0b-4ef8-bb6d-6bb9bd380a55', 75000.25, 75000.25, 'INR')
			ON CONFLICT (account_id) DO NOTHING`,
			nil,
		},
		{
			"beneficiaries",
			`INSERT INTO beneficiaries (user_id, beneficiary_name, bank_name, account_number, ifsc, nickname) VALUES
				('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Alice Johnson', 'HDFC Bank', '1111222233334444', 'HDFC0001234', 'Alice'),
				('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Bob Williams', 'ICICI Bank', '5555666677778888', 'ICIC0005678', 'Bob'),
				('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Charlie Brown', 'SBI', '9999888877776666', 'SBIN0009012', 'Charlie')`,
			nil,
		},
		{
			"transfer_modes",
			`INSERT INTO transfer_modes (code, name, description, is_active) VALUES
				('UPI', 'UPI', 'Unified Payments Interface - instant transfer up to ₹1 lakh', true),
				('NEFT', 'NEFT', 'National Electronic Funds Transfer - batch settlement', true),
				('RTGS', 'RTGS', 'Real Time Gross Settlement - high value instant transfer', true),
				('IMPS', 'IMPS', 'Immediate Payment Service - 24x7 instant transfer', true)
			ON CONFLICT (code) DO NOTHING`,
			nil,
		},
	}

	for _, q := range queries {
		if _, err := tx.Exec(ctx, q.query, q.args...); err != nil {
			return fmt.Errorf("seed %s: %w", q.name, err)
		}
		fmt.Printf("seeded %s\n", q.name)
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Verify connectivity
	var count int
	if err := db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return fmt.Errorf("verify seed: %w", err)
	}
	if count == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
