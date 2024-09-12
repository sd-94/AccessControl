package repositories

import (
	"context"
	"prac/types"
)

const (
	GET_ACC      string = "SELECT acc_id, first_name, last_name, email, password, rights FROM accounts WHERE acc_id = $1 AND rights != $2"
	GET_ACCS     string = "SELECT acc_id, first_name, last_name, email, password, rights FROM accounts WHERE rights != $1"
	CREATE_ACC   string = "INSERT INTO accounts (first_name, last_name, email, password, rights) VALUES ($1, $2, $3, $4, $5) RETURNING acc_id;"
	UPDATE_ACC   string = "UPDATE accounts SET first_name = $1, last_name = $2, email = $3, password = $4, rights = $5 WHERE acc_id = $6 AND rights != $7"
	DELETE_ACC   string = "DELETE FROM accounts WHERE acc_id = $1 AND rights != $2"
	GET_ROLE     string = "SELECT rights FROM accounts WHERE email = $1"
	GET_ACC_AUTH string = "SELECT EXISTS (SELECT 1 FROM accounts WHERE email = $1 AND password = $2);"
)

func (repo *Repository) GetAccount(ctx context.Context, id string) (*types.Account, error) {
	var acc types.Account

	err := repo.DB.QueryRowContext(ctx, GET_ACC, id, types.Tier0).Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Email, &acc.Password, &acc.Rights)
	if err != nil {
		return nil, types.SQLExecutionError.Wrap(err, "failed to retrieve account")
	}

	return &acc, nil
}

func (repo *Repository) GetAccounts(ctx context.Context) ([]types.Account, error) {
	var accs []types.Account

	rows, err := repo.DB.QueryContext(ctx, GET_ACCS, types.Tier0)
	if err != nil {
		return nil, types.SQLExecutionError.Wrap(err, "failed to retrieve accounts")
	}
	defer rows.Close()

	for rows.Next() {
		var acc types.Account
		if err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Email, &acc.Password, &acc.Rights); err != nil {
			return nil, types.SQLExecutionError.Wrap(err, "failed to scan account")
		}
		accs = append(accs, acc)
	}

	return accs, nil
}

func (repo *Repository) CreateAccount(ctx context.Context, acc *types.Account) (string, error) {
	var id string
	err := repo.DB.QueryRowContext(ctx, CREATE_ACC, acc.FirstName, acc.LastName, acc.Email, acc.Password, acc.Rights).Scan(&id)
	if err != nil {
		return "", types.SQLExecutionError.Wrap(err, "failed to create account")
	}

	return id, nil
}

func (repo *Repository) UpdateAccount(ctx context.Context, acc *types.Account, id string) error {
	_, err := repo.DB.ExecContext(ctx, UPDATE_ACC, acc.FirstName, acc.LastName, acc.Email, acc.Password, acc.Rights, id, types.Tier0)
	if err != nil {
		return types.SQLExecutionError.Wrap(err, "failed to update account")
	}

	return nil
}

func (repo *Repository) DeleteAccount(ctx context.Context, id string) error {
	_, err := repo.DB.ExecContext(ctx, DELETE_ACC, id, types.Tier0)
	if err != nil {
		return types.SQLExecutionError.Wrap(err, "failed to delete account")
	}

	return nil
}

func (repo *Repository) FindRoleByAccountEmail(email string) (role string, err error) {
	err = repo.DB.QueryRow(GET_ROLE, email).Scan(&role)
	if err != nil {
		return "", types.SQLExecutionError.Wrap(err, "failed to find role for account")
	}

	return role, nil
}

func (repo *Repository) SignIn(body *types.SignIn) (bool, error) {
	var exists bool
	err := repo.DB.QueryRow(GET_ACC_AUTH, body.Email, body.Password).Scan(&exists)
	if err != nil {
		return false, types.SQLExecutionError.Wrap(err, "failed to authenticate user")
	}
	return exists, nil
}
