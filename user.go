package wc

import "yasty.org/peter/wc/ecs"

type User struct {
	Id       int64  `json:"id,string"`
	Username string `json:"username"`
}

type UserStore interface {
	Get(name string) (User, error)
	Create(User) (User, error)
	Save(User) error
}

type MemUsers struct {
	Users []User
}

func (m *MemUsers) Get(name string) (User, error) {
	for _, u := range m.Users {
		if u.Username == name {
			return u, nil
		}
	}

	return User{}, ecs.ErrNotFound
}

func (m *MemUsers) Create(u User) (User, error) {
	if u.Id != 0 {
		return User{}, ecs.ErrHasId
	}

	if _, err := m.Get(u.Username); err == nil {
		return User{}, ecs.ErrExists
	}

	u.Id = int64(len(m.Users))
	m.Users = append(m.Users, u)
	return u, nil
}

func (m *MemUsers) Save(u User) error {
	m.Users[u.Id] = u
	return nil
}
