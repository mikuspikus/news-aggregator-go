package accounts

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mikuspikus/news-aggregator-go/pkg/simple-token-storage"
	"github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/refresh-token-storage"
	"github.com/mikuspikus/news-aggregator-go/services/accounts/pkg/user-token-storage"
	pb "github.com/mikuspikus/news-aggregator-go/services/accounts/proto"
	"math"
	"testing"
	"time"
)

type mockDB struct {
	Users  []*User
	LastID int32
}

func _safeSlice(users []*User, pageNumber, pageSize int32) []*User {
	if pageNumber < 0 || pageSize < 0 {
		return make([]*User, 0)
	}

	startIndex := int(pageNumber * pageSize)
	if startIndex < 0 {
		return make([]*User, 0)
	}

	endIndex := startIndex + int(pageSize)
	if endIndex > len(users) {
		endIndex = len(users)
	}

	return users[startIndex:endIndex]

}

func (db *mockDB) find(uid uuid.UUID) (int, *User, error) {
	for idx, user := range db.Users {
		if user.Uid == uid {
			return idx, user, nil
		}
	}

	return -1, nil, errNotFound
}

func (db *mockDB) findByUsername(username string) (int, *User, error) {
	for idx, user := range db.Users {
		if user.Username == username {
			return idx, user, nil
		}
	}

	return -1, nil, errNotFound
}

func (db *mockDB) push(username string) *User {
	user := new(User)

	now := time.Now()

	user.Uid = uuid.New()
	user.Username = username
	user.Created = now
	user.Edited = now
	user.IsAdmin = false

	update := append(db.Users, user)
	db.Users = update
	db.LastID++
	return user
}

func (db *mockDB) List(pageNumber, pageSize int32) ([]*User, int32, error) {
	users := _safeSlice(db.Users, pageNumber, pageSize)
	pageCount := math.Ceil(float64(len(db.Users)) / float64(pageSize))

	return users, int32(pageCount), nil
}

func (db *mockDB) Get(uid uuid.UUID) (*User, error) {
	_, user, err := db.find(uid)
	return user, err
}

func (db *mockDB) Create(username, _ string) (*User, error) {
	user := db.push(username)
	return user, nil
}

func (db *mockDB) Update(uid uuid.UUID, username, _ string) (*User, error) {
	_, user, err := db.find(uid)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Edited = time.Now()
	return user, nil
}

func (db *mockDB) Delete(uid uuid.UUID) error {
	idx, _, err := db.find(uid)
	if err != nil {
		return err
	}

	users := db.Users
	users[len(users)-1], users[idx] = users[idx], users[len(users)-1]
	db.Users = users[:len(users)-1]

	return nil
}

func (db *mockDB) AdminUpdate(uid uuid.UUID, username string, isAdmin bool) (*User, error) {
	_, user, err := db.find(uid)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.IsAdmin = isAdmin

	return user, nil
}

func (db *mockDB) CheckPassword(_ uuid.UUID, _ string) (bool, error) {
	return true, nil
}

func (db *mockDB) GetUserByUsername(username string) (*User, error) {
	_, user, err := db.findByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *mockDB) Close() {}

func setupDB() *mockDB {
	db := new(mockDB)

	for i := 0; i < 25; i++ {
		db.push(fmt.Sprintf("testuser-%d", i))
	}

	return db
}

func setupServer(db *mockDB) *Service {
	server := &Service{
		Store:         db,
		ServiceTokens: &simple_token_storage.MockStorage{},
		RefreshTokens: &refresh_token_storage.MockStorage{},
		UserTokens:    &user_token_storage.MockStorage{},
	}
	return server
}

func TestService_ListUsers(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.ListUsersRequest{
		PageSize:   3,
		PageNumber: 1,
	}

	_, err := server.ListUsers(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_GetUser(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	request := &pb.GetUserRequest{Uid: db.Users[0].Uid.String()}
	_, err := server.GetUser(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AddUser(t *testing.T) {
	db := setupDB()
	server := setupServer(db)

	request := &pb.AddUserRequest{
		ApiToken: "",
		Username: "exampleuser",
		Password: "exampleuser",
	}
	_, err := server.AddUser(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_EditUser(t *testing.T) {
	db := setupDB()
	server := setupServer(db)

	target := db.Users[0]

	request := &pb.EditUserRequest{
		ApiToken: "",
		Uid:      target.Uid.String(),
		Username: "new-username",
		Password: "new-password",
	}

	_, err := server.EditUser(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_DeleteUser(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	target := db.Users[0]

	request := &pb.DeleteUserRequest{
		ApiToken: "",
		Uid:      target.Uid.String(),
	}
	_, err := server.DeleteUser(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}

func TestService_AdminEditUser(t *testing.T) {
	db := setupDB()
	server := setupServer(db)
	target := db.Users[0]

	request := &pb.AdminEditUserRequest{
		ApiToken: "",
		Uid:      target.Uid.String(),
		IsAdmin:  true,
		Username: "admin-edit-user",
	}
	_, err := server.AdminEditUser(context.Background(), request)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}
}
