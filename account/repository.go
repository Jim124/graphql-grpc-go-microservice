package account

import (
	"context"
	"database/sql"
)

type Account struct {
	ID   string
	Name string
}
type Repository interface {
	Close()
	PutAccount(ctx context.Context, a Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccounts(ctx context.Context, skip, take uint64) ([]Account, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close()
}

func (r *postgresRepository) Ping() error {
	return r.db.Ping()
}

func (r *postgresRepository) PutAccount(ctx context.Context, a Account) error {
	_, err := r.db.ExecContext(ctx, "insert into accounts(id,name) values($1,$2)", a.ID, a.Name)
	return err
}

func (r *postgresRepository) GetAccountById(ctx context.Context, id string) (*Account, error) {
	row := r.db.QueryRowContext(ctx, "select id,name from accounts where id=$1", id)
	a := &Account{}
	if err := row.Scan(&a.ID, &a.Name); err != nil {
		return nil, err
	}
	return a, nil
}

func (r *postgresRepository) ListAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	rows, err := r.db.QueryContext(ctx, "select id,name from accounts order by id desc offset $1 limit $2", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []Account
	for rows.Next() {
		a := &Account{}
		if err = rows.Scan(&a.ID, &a.Name); err == nil {
			accounts = append(accounts, *a)
		}

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}
