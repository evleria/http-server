package main

import "errors"

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNotFound       = errors.New("not found")
)

type UserRepository interface {
	GetAllUsers() ([]User, error)
	CreateUser(email string, firstName string, lastName string, phone string) (*User, error)
	DeleteUser(id int) error
	GetUserById(id int) (*User, error)
}

type InMemoryUserRepository struct {
	users          []*User
	userIdIndex    map[int]*User
	userEmailIndex map[string]*User
	lastId         int
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:          []*User{},
		userIdIndex:    map[int]*User{},
		userEmailIndex: map[string]*User{},
		lastId:         0,
	}
}

func (r *InMemoryUserRepository) GetAllUsers() ([]User, error) {
	users := make([]User, len(r.users))
	for i, user := range r.users {
		users[i] = *user
	}
	return users, nil
}

func (r *InMemoryUserRepository) GetUserById(id int) (*User, error) {
	user, ok := r.userIdIndex[id]
	if !ok {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) CreateUser(email string, firstName string, lastName string, phone string) (*User, error) {
	if _, ok := r.userEmailIndex[email]; ok {
		return nil, ErrDuplicateEmail
	}

	user := &User{
		Id:        r.getNewId(),
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
	}

	r.users = append(r.users, user)
	r.userIdIndex[user.Id] = user
	r.userEmailIndex[user.Email] = user

	return user, nil
}

func (r *InMemoryUserRepository) DeleteUser(id int) error {
	user, ok := r.userIdIndex[id]
	if !ok {
		return ErrNotFound
	}

	i := findUserIndex(r.users, user)
	r.users = append(r.users[:i], r.users[i+1:]...)
	delete(r.userIdIndex, user.Id)
	delete(r.userEmailIndex, user.Email)

	return nil
}

func (r *InMemoryUserRepository) getNewId() int {
	r.lastId++
	return r.lastId
}

func findUserIndex(users []*User, user *User) int {
	for i, u := range users {
		if u == user {
			return i
		}
	}
	return -1
}
