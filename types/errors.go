package types

import "github.com/joomcode/errorx"

var (
	DBNamespace        = errorx.NewNamespace("database")
	ConnectionError    = DBNamespace.NewType("connection_error")
	TableCreationError = DBNamespace.NewType("table_creation_error")
	AdminAccountError  = DBNamespace.NewType("admin_account_error")

	SQLNamespace       = errorx.NewNamespace("sql")
	SQLExecutionError  = SQLNamespace.NewType("sql_execution_error")

	ValidationNamespace = errorx.NewNamespace("validation")
	ValidationError     = ValidationNamespace.NewType("validation_error")

	AuthNamespace      = errorx.NewNamespace("authorization")
	AuthorizationError = AuthNamespace.NewType("authorization_error")

	TokenNamespace     = errorx.NewNamespace("token")
	TokenGenerationError = TokenNamespace.NewType("token_generation_error")

	NotFoundNamespace = errorx.NewNamespace("not_found")
	NotFoundError     = NotFoundNamespace.NewType("resource_not_found_error")

	JSONNamespace     = errorx.NewNamespace("json")
	JSONEncodingError = JSONNamespace.NewType("json_encoding_error")
	JSONDecodingError = JSONNamespace.NewType("json_decoding_error")
)
