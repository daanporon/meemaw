// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: query.sql

package server

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/tabbed/pqtype"
)

const addDevice = `-- name: AddDevice :one
INSERT INTO devices (user_id, wallet_id, user_agent)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING
RETURNING id, user_id, wallet_id, user_agent
`

type AddDeviceParams struct {
	UserId    int64
	WalletId  int64
	UserAgent string
}

func (q *Queries) AddDevice(ctx context.Context, arg AddDeviceParams) (Device, error) {
	row := q.db.QueryRowContext(ctx, addDevice, arg.UserId, arg.WalletId, arg.UserAgent)
	var i Device
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.WalletID,
		&i.UserAgent,
	)
	return i, err
}

const addUser = `-- name: AddUser :one

INSERT INTO users (foreign_key)
VALUES ($1)
ON CONFLICT DO NOTHING
RETURNING id, foreign_key
`

// ----- INSERTS -------
func (q *Queries) AddUser(ctx context.Context, foreignKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, addUser, foreignKey)
	var i User
	err := row.Scan(&i.ID, &i.ForeignKey)
	return i, err
}

const addWallet = `-- name: AddWallet :one
INSERT INTO wallets (user_id, public_address, share, params)
VALUES ($1, $2, $3, $4)
ON CONFLICT DO NOTHING
RETURNING id, user_id, public_address, share, params
`

type AddWalletParams struct {
	UserId        int64
	PublicAddress string
	Share         string
	Params        json.RawMessage
}

func (q *Queries) AddWallet(ctx context.Context, arg AddWalletParams) (Wallet, error) {
	row := q.db.QueryRowContext(ctx, addWallet,
		arg.UserId,
		arg.PublicAddress,
		arg.Share,
		arg.Params,
	)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PublicAddress,
		&i.Share,
		&i.Params,
	)
	return i, err
}

const getFirstUser = `-- name: GetFirstUser :one

SELECT id, foreign_key FROM users
LIMIT 1
`

// ----- SELECTS -------
func (q *Queries) GetFirstUser(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getFirstUser)
	var i User
	err := row.Scan(&i.ID, &i.ForeignKey)
	return i, err
}

const getUserByAddress = `-- name: GetUserByAddress :one
SELECT users.id, users.foreign_key 
FROM wallets
INNER JOIN users ON wallets.user_id = users.id
WHERE wallets.public_address = $1
`

func (q *Queries) GetUserByAddress(ctx context.Context, publicaddress string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByAddress, publicaddress)
	var i User
	err := row.Scan(&i.ID, &i.ForeignKey)
	return i, err
}

const getUserByForeignKey = `-- name: GetUserByForeignKey :one
SELECT id, foreign_key FROM users
WHERE foreign_key = $1
LIMIT 1
`

func (q *Queries) GetUserByForeignKey(ctx context.Context, foreignkey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByForeignKey, foreignkey)
	var i User
	err := row.Scan(&i.ID, &i.ForeignKey)
	return i, err
}

const getUserDevices = `-- name: GetUserDevices :many
SELECT id, user_id, wallet_id, user_agent FROM devices
WHERE user_id = $1
`

func (q *Queries) GetUserDevices(ctx context.Context, userid int64) ([]Device, error) {
	rows, err := q.db.QueryContext(ctx, getUserDevices, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Device
	for rows.Next() {
		var i Device
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.WalletID,
			&i.UserAgent,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSigningParameters = `-- name: GetUserSigningParameters :one
SELECT wallets.id, wallets.user_id, wallets.public_address, wallets.share, wallets.params
FROM users
LEFT JOIN wallets ON users.id = wallets.user_id
WHERE users.foreign_key = $1
`

type GetUserSigningParametersRow struct {
	ID            sql.NullInt64
	UserID        sql.NullInt64
	PublicAddress sql.NullString
	Share         sql.NullString
	Params        pqtype.NullRawMessage
}

func (q *Queries) GetUserSigningParameters(ctx context.Context, foreignkey string) (GetUserSigningParametersRow, error) {
	row := q.db.QueryRowContext(ctx, getUserSigningParameters, foreignkey)
	var i GetUserSigningParametersRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PublicAddress,
		&i.Share,
		&i.Params,
	)
	return i, err
}

const getUserWallets = `-- name: GetUserWallets :many
SELECT id, user_id, public_address, share, params FROM wallets
WHERE user_id = $1
`

func (q *Queries) GetUserWallets(ctx context.Context, userid int64) ([]Wallet, error) {
	rows, err := q.db.QueryContext(ctx, getUserWallets, userid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wallet
	for rows.Next() {
		var i Wallet
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.PublicAddress,
			&i.Share,
			&i.Params,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWalletByAddress = `-- name: GetWalletByAddress :one
SELECT id, user_id, public_address, share, params FROM wallets
WHERE public_address = $1
`

func (q *Queries) GetWalletByAddress(ctx context.Context, publicaddress string) (Wallet, error) {
	row := q.db.QueryRowContext(ctx, getWalletByAddress, publicaddress)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PublicAddress,
		&i.Share,
		&i.Params,
	)
	return i, err
}

const status = `-- name: Status :one
SELECT 1
`

func (q *Queries) Status(ctx context.Context) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, status)
	var column_1 interface{}
	err := row.Scan(&column_1)
	return column_1, err
}
