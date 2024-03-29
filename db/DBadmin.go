package db

import (
	"context"
)

type temp_user struct {
	Role_id  uint8
	Username string
	Password string
	Id       uint16
}

type User struct {
	Id   uint16
	Role string
	Name string
}

type Role struct {
	Id   uint8
	Name string
}

func GetUsers() []User {
	tmp_users := make([]temp_user, 0)
	all_users := make([]User, 0)

	rows, err := conn.Query(context.Background(), "select * from goodbuy.users")
	if err != nil {
		println(err.Error())
		return all_users
	}

	for rows.Next() {
		var usr temp_user
		err = rows.Scan(
			&usr.Role_id,
			&usr.Username,
			&usr.Password,
			&usr.Id,
		)
		if err != nil {
			println("Something on 36", err.Error())
			return all_users
		}

		tmp_users = append(tmp_users, usr)
	}
	rows.Close()

	for _, usr := range tmp_users {
		user := User{usr.Id, GetRoleById(usr.Role_id), usr.Username}
		all_users = append(all_users, user)
	}

	return all_users
}

func NewUser() User {
	return User{0, "Unknown", "NotUser"}
}

func RegisterUser(username, password string) {
	_, err := conn.Exec(context.Background(), "insert into goodbuy.users values (2, $1, $2);", username, password)
	if err != nil {
		println(err.Error())
	}
}

func GetUsernameById(id uint8) string {
	var name string

	err := conn.QueryRow(context.Background(), "select username from goodbuy.users where id=$1", id).Scan(&name)
	if err != nil {
		println("Something on 67", err.Error())
		return "Unknown"
	}

	return name
}

func GetRoles() []Role {
	roles := make([]Role, 0)

	rows, err := conn.Query(context.Background(), "select * from goodbuy.roles")
	if err != nil {
		println("Something on 87.", err.Error())
	}

	for rows.Next() {
		var role Role
		rows.Scan(&role.Name, &role.Id)
		roles = append(roles, role)
	}
	rows.Close()

	return roles
}

func GetRoleOfUser(username string) string {
	role := "error"
	err := conn.QueryRow(context.Background(), "select name from goodbuy.roles where id = (select role_id from goodbuy.users where username=$1)", username).Scan(&role)
	if err != nil {
		println(err.Error())
	}
	return role
}

func GetRoleById(id uint8) string {
	var role string

	err := conn.QueryRow(context.Background(), "select name from goodbuy.roles where id=$1", id).Scan(&role)
	if err != nil {
		println("Something on 76", err.Error())
		return "Unknown"
	}

	return role
}

func GetIdOfRole(role string) uint8 {
	var id uint8

	err := conn.QueryRow(context.Background(), "select id from goodbuy.roles where name=$1", role).Scan(&id)
	if err != nil {
		println("Something on 126", err.Error())
		return 0
	}

	return id
}

func GrantRoleToUser(user string, role_id int) {
	_, err := conn.Exec(context.Background(), "update goodbuy.users set role_id=$1 where username=$2", role_id, user)
	if err != nil {
		println("Failed to grant the role.", err.Error())
	}
}

func AddMarket(name string, date_start, date_end string, fee float64) {
	_, err := conn.Exec(context.Background(), "call goodbuy.add_market($1, $2, $3, $4)", name, date_start, date_end, fee)
	if err != nil {
		println("Failed to add the market.", err.Error())
	}
}
