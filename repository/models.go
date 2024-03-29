// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package repository

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleSupplier UserRole = "supplier"
	UserRoleAdmin    UserRole = "admin"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole
	Valid    bool // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type Session struct {
	ID           uuid.UUID
	UserID       int64
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type SessionIp struct {
	ID        int64
	SessionID uuid.UUID
	UserAgent sql.NullString
	Ip        string
}

type User struct {
	ID                int64
	Email             string
	UserName          string
	Role              UserRole
	ActiveStatus      bool
	HashedPassword    string
	PasswordUpdatedAt time.Time
	CreatedAt         time.Time
	Gender            sql.NullString
	Phone             sql.NullString
	Name              sql.NullString
	Address           sql.NullString
	Note              sql.NullString
}

type UserAddress struct {
	ID        int64
	UserID    int64
	Address   string
	Note      sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserProfile struct {
	ID        int64
	UserID    int64
	Phone     sql.NullString
	Avatar    sql.NullString
	CreatedAt time.Time
	UpdatedAt time.Time
}
