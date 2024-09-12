package handlers

import (
	"encoding/json"
	"net/http"
	"prac/types"

	"github.com/gorilla/mux"
)

func (handler *Handler) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	role, err := handler.Repository.FindRoleByAccountEmail(email)
	if err != nil {
		return types.AuthorizationError.Wrap(err, "couldn't find role for account")
	}

	if (role != types.Tier0) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission denied"))
		return types.AuthorizationError.New("Permission denied")
	}

	params := mux.Vars(r)

	acc, err := handler.Repository.GetAccount(ctx, params["acc_id"])
	if err != nil {
		return types.SQLExecutionError.Wrap(err, "couldn't retrieve account")
	}

	if err := json.NewEncoder(w).Encode(acc); err != nil {
		return types.JSONEncodingError.Wrap(err, "failed to encode account to JSON")
	}

	return nil
}

func (handler *Handler) HandleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	role, err := handler.Repository.FindRoleByAccountEmail(email)
	if err != nil {
		return types.AuthorizationError.Wrap(err, "couldn't find role for account")
	}

	if (role != types.Tier0) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission denied"))
		return types.AuthorizationError.New("Permission denied")
	}

	accs, err := handler.Repository.GetAccounts(ctx)
	if err != nil {
		return types.SQLExecutionError.Wrap(err, "couldn't retrieve accounts")
	}

	if err := json.NewEncoder(w).Encode(accs); err != nil {
		return types.JSONEncodingError.Wrap(err, "failed to encode accounts to JSON")
	}

	return nil
}

func (handler *Handler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	role, err := handler.Repository.FindRoleByAccountEmail(email)
	if err != nil {
		return types.AuthorizationError.Wrap(err, "couldn't find role for account")
	}

	if (role != types.Tier0) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission denied"))
		return types.AuthorizationError.New("Permission denied")
	}

	var acc types.Account

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		return types.JSONDecodingError.Wrap(err, "failed to decode account from request body")
	}

	id, err := handler.Repository.CreateAccount(ctx, &acc)
	if err != nil {
		return types.SQLExecutionError.Wrap(err, "failed to create account")
	}

	acc.ID = id

	if err := json.NewEncoder(w).Encode(&acc); err != nil {
		return types.JSONEncodingError.Wrap(err, "failed to encode account to JSON")
	}

	return nil
}

func (handler *Handler) HandleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	role, err := handler.Repository.FindRoleByAccountEmail(email)
	if err != nil {
		return types.AuthorizationError.Wrap(err, "couldn't find role for account")
	}

	if (role != types.Tier0) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission denied"))
		return types.AuthorizationError.New("Permission denied")
	}

	params := mux.Vars(r)

	var acc types.Account

	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		return types.JSONDecodingError.Wrap(err, "failed to decode account from request body")
	}

	if err := handler.Repository.UpdateAccount(ctx, &acc, params["acc_id"]); err != nil {
		return types.SQLExecutionError.Wrap(err, "failed to update account")
	}

	if err := json.NewEncoder(w).Encode("Account updated successfully"); err != nil {
		return types.JSONEncodingError.Wrap(err, "failed to encode success message to JSON")
	}

	return nil
}

func (handler *Handler) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	email := ctx.Value("email").(string)

	role, err := handler.Repository.FindRoleByAccountEmail(email)
	if err != nil {
		return types.AuthorizationError.Wrap(err, "couldn't find role for account")
	}

	if (role != types.Tier0) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Permission denied"))
		return types.AuthorizationError.New("Permission denied")
	}

	params := mux.Vars(r)

	if err := handler.Repository.DeleteAccount(ctx, params["acc_id"]); err != nil {
		return types.SQLExecutionError.Wrap(err, "failed to delete account")
	}

	if err := json.NewEncoder(w).Encode("Account deleted successfully"); err != nil {
		return types.JSONEncodingError.Wrap(err, "failed to encode success message to JSON")
	}

	return nil
}
