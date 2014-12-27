package wc

import "yasty.org/peter/wc/ecs"

type User struct {
	Id       int64  `json:"id,string"`
	Username string `json:"username"`
}

var users = []User{User{}}

func exists(name string) bool {
	_, ok := getUser(name)
	return ok
}

func getUser(name string) (User, bool) {
	for _, u := range users {
		if u.Username == name {
			return u, true
		}
	}

	return User{}, false
}

func createUser(u User) (User, error) {
	if u.Id != 0 {
		return User{}, ecs.ErrHasId
	}

	if exists(u.Username) {
		return User{}, ecs.ErrExists
	}

	u.Id = int64(len(users))
	users = append(users, u)
	return u, nil
}

func saveUser(u User) User {
	users[u.Id] = u
	return u
}
