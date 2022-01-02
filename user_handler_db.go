package userapp

import (
	"errors"

	"github.com/guregu/dynamo"
)

func scanUsers() ([]User, error) {
	var resp []User
	table := gdb.Table(usersTable)
	if err := table.Scan().All(&resp); err != nil {
		if errors.Is(err, dynamo.ErrNotFound) {
			return nil, nil
		}
		return resp, err
	}
	return resp, nil
}

func putUser(u User) error {
	if err := gdb.Table(usersTable).Put(u).If("attribute_not_exists(user_id)").Run(); err != nil {
		return err
	}
	return nil
}

func deleteUser(user_id string) error {
	if err := gdb.Table(usersTable).Delete("user_id", user_id).Run(); err != nil {
		return err
	}
	return nil
}

func updateUser(u User) error {
	if err := gdb.Table(usersTable).Update("user_id", u.UserID).Set("user_name", u.UserName).Run(); err != nil {
		return err
	}
	return nil
}

func getUser(user_id string) (User, error) {
	var resp User
	if err := gdb.Table(usersTable).Get("user_id", user_id).One(&resp); err != nil {
		return resp, err
	}
	return resp, nil
}
