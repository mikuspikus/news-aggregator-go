package accounts

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"math"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 5
)

var (
	errNotFound             = errors.New("user not found")
	errNotCreated           = errors.New("user not created")
	errUsernameAlreadyTaken = errors.New("this username already taken")
)

// User defines inner representation of user
type User struct {
	Uid      uuid.UUID
	Username string
	Created  time.Time
	Edited   time.Time
	IsAdmin  bool
}

type DataStoreHandler interface {
	List(pageNumber, pageSize int32) ([]*User, int32, error)
	Get(uid uuid.UUID) (*User, error)
	Create(username, password string) (*User, error)
	Update(uid uuid.UUID, username, password string) (*User, error)
	Delete(uid uuid.UUID) error

	// Admin panel
	AdminUpdate(uid uuid.UUID, username string, isAdmin bool) (*User, error)

	CheckPassword(uid uuid.UUID, password string) (bool, error)
	GetUserByUsername(username string) (*User, error)
	Close()
}

type db struct {
	*pgxpool.Pool
}

func NewDB(connString string) (*db, error) {
	dbpool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	return &db{dbpool}, nil
}

func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcryptCost)
	return string(hash), err
}

func checkPassword(pwd, hashedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd))
	return err == nil
}

// _pgErrCode tries to convert err into pgconn.PgErr and extract its code. Returns code and conversion success status
func _pgErrCode(err error) (string, bool) {
	if pgerr, ok := err.(*pgconn.PgError); ok {
		return pgerr.Code, true
	}
	return "", false

}

func (db *db) pageCount(pageSize int32) (int32, error) {
	query := "select count(*) from users "

	row := db.QueryRow(context.Background(), query)
	var rowsCount int32
	err := row.Scan(&rowsCount)
	if err != nil {
		return 0, err
	}
	pageCount := math.Ceil(float64(rowsCount) / float64(pageSize))

	return int32(pageCount), nil
}

func (db *db) List(pageNumber, pageSize int32) ([]*User, int32, error) {
	query := "select uid, username, created_at, edited_at, is_admin from users limit $1 offset $2"
	lastRecord := pageNumber * pageSize

	rows, err := db.Query(context.Background(), query, pageSize, lastRecord)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		var userUUID string

		err := rows.Scan(&userUUID, &user.Username, &user.Created, &user.Edited, &user.IsAdmin)
		if err != nil {
			return nil, 0, err
		}

		user.Uid, err = uuid.Parse(userUUID)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	pageCount, err := db.pageCount(pageSize)
	if err != nil {
		return nil, 0, err
	}

	return users, pageCount, nil
}

func (db *db) Get(uid uuid.UUID) (*User, error) {
	query := "select username, created_at, edited_at, is_admin from users where uid=$1"
	user := new(User)
	user.Uid = uid
	err := db.QueryRow(context.Background(), query, uid.String()).Scan(&user.Username, &user.Created, &user.Edited, &user.IsAdmin)
	switch err {
	case nil:
		return user, nil
	case pgx.ErrNoRows:
		return nil, errNotFound
	default:
		return nil, err

	}
}

func (db *db) Create(username, password string) (*User, error) {
	query := "insert into users (uid, username, password_hash, created_at, edited_at) " +
		"values ($1, $2, $3, $4, $5) " +
		"returning uid"
	user := new(User)
	now := time.Now()
	uid := uuid.New()

	user.Uid = uid
	user.Username = username
	user.Created = now
	user.Edited = now

	hashedPwd, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	cmd, err := db.Exec(context.Background(), query, uid.String(), username, hashedPwd, now, now)
	if err != nil {
		if code, ok := _pgErrCode(err); ok && code == "23505" {
			//Code 203505 for 'unique key violation' error
			//username is considered to be the only one unique field in users table
			return nil, errUsernameAlreadyTaken
		}
		return nil, err
	}

	rowsAffected := cmd.RowsAffected()
	if rowsAffected != 1 {
		return nil, errNotCreated
	}

	return user, nil
}

func (db *db) Update(uid uuid.UUID, username, password string) (*User, error) {
	query := "update users set username=$1 password_hash=$2 edited_at=$3 where uid=$4 returning created_at, is_admin"

	user := new(User)
	now := time.Now()
	hashedPwd, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user.Uid = uid
	user.Username = username
	user.Edited = now

	err = db.QueryRow(context.Background(), query, username, hashedPwd, now).Scan(&user.Created, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *db) Delete(uid uuid.UUID) error {
	query := "delete from users where uid=$1"

	cmd, err := db.Exec(context.Background(), query, uid)
	if err != nil {
		return err
	}
	if rowsAffected := cmd.RowsAffected(); rowsAffected == 0 {
		return errNotFound
	}
	return nil
}

func (db *db) AdminUpdate(uid uuid.UUID, username string, isAdmin bool) (*User, error) {
	query := "update users set username=$1, is_admin=$2 where uid=$3 returning created_at, edited_at"

	user := new(User)
	user.Uid = uid
	user.Username = username
	user.IsAdmin = isAdmin

	err := db.QueryRow(context.Background(), query, username, isAdmin, uid).Scan(&user.Created, &user.Edited)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *db) CheckPassword(uid uuid.UUID, password string) (bool, error) {
	query := "select password_hash from users where uid=$1"
	var hashedPwd string
	err := db.QueryRow(context.Background(), query, uid.String()).Scan(&hashedPwd)
	switch err {
	case nil:
		return checkPassword(password, hashedPwd), nil
	case pgx.ErrNoRows:
		return false, errNotFound
	default:
		return false, err
	}
}

func (db *db) GetUserByUsername(username string) (*User, error) {
	query := "select uid, created_at, edited_at, is_admin from users where username=$1"
	user := new(User)
	var struid string
	err := db.QueryRow(context.Background(), query, username).Scan(&struid, &user.Created, &user.Edited, &user.IsAdmin)
	switch err {
	case nil:
		uid, err := uuid.Parse(struid)
		if err != nil {
			return nil, err
		}
		user.Uid = uid
		user.Username = username
		return user, nil

	case pgx.ErrNoRows:
		return nil, errNotFound
	default:
		return nil, err
	}
}
