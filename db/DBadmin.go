package db

import (
	"context"
	"fmt"
	"strings"

	pgquery "github.com/pganalyze/pg_query_go/v5"
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

	rows, err := conn.Query(context.Background(), "select role_id, username, id from goodbuy.users where isDeleted=false order by id")
	if err != nil {
		println(err.Error())
		return all_users
	}

	for rows.Next() {
		var usr temp_user
		err = rows.Scan(
			&usr.Role_id,
			&usr.Username,
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
		user := User{usr.Id, GetRolenameByID(usr.Role_id), usr.Username}
		all_users = append(all_users, user)
	}

	return all_users
}

func IsUserExistByID(userID uint8) bool {
	username := GetUsernameById(userID)

	if username != "unknown" {
		return IsUserExist(username)
	}

	return false
}

func IsUserExist(username string) bool {
	exist := false

	rows, err := conn.Query(context.Background(), "select id from goodbuy.users where isDeleted=false and username=$1", username)
	if err != nil {
		println("Error on getting user", err.Error())
		return exist
	}

	for rows.Next() {
		exist = true
		break
	}
	rows.Close()

	return exist
}

func NewUser() User {
	return User{0, "Unknown", "NotUser"}
}

func RegisterUser(roleID uint8, username, password string) {
	_, err := conn.Exec(context.Background(), "insert into goodbuy.users values ($1, $2, $3);", roleID, username, password)
	if err != nil {
		println(err.Error())
	}
}

func DisableUserByID(userID uint8) {
	_, err := conn.Exec(context.Background(), "update goodbuy.users set isDeleted=true where id=$1", userID)
	if err != nil {
		println(err.Error())
	}
}

func GetUsernameById(id uint8) string {
	var name string

	err := conn.QueryRow(context.Background(), "select username from goodbuy.users where isDeleted=false and id=$1", id).Scan(&name)
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

func GetRolenameOfUserByName(username string) string {
	role := "error"
	err := conn.QueryRow(context.Background(), "select name from goodbuy.roles where id = (select role_id from goodbuy.users where isDeleted=false and username=$1)", username).Scan(&role)
	if err != nil {
		println(err.Error())
	}
	return role
}

func GetRolenameOfUserById(id uint8) string {
	var role string

	err := conn.QueryRow(context.Background(), "select name from goodbuy.roles where id=$1", id).Scan(&role)
	if err != nil {
		println("Something on 76", err.Error())
		return "Unknown"
	}

	return role
}

func GetRolenameByID(roleID uint8) string {
	role := "error"
	err := conn.QueryRow(context.Background(), "select name from goodbuy.roles where id=$1", roleID).Scan(&role)
	if err != nil {
		println(err.Error())
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

func GrantRoleToUser(userID uint8, role_id int) {
	_, err := conn.Exec(context.Background(), "update goodbuy.users set role_id=$1 where isDeleted=false and id=$2", role_id, userID)
	if err != nil {
		println("Failed to grant the role.", err.Error())
	}
}

func ChangeUserPassword(userID uint8, password string) {
	_, err := conn.Exec(context.Background(), "update goodbuy.users set password=$1 where isDeleted=false and id=$2", password, userID)
	if err != nil {
		println("Failed to change the password.", err.Error())
	}
}

func AddMarket(name string, date_start, date_end string, fee float64) {
	_, err := conn.Exec(context.Background(), "call goodbuy.add_market($1, $2, $3, $4)", name, date_start, date_end, fee)
	if err != nil {
		println("Failed to add the market.", err.Error())
	}
}

func RunSqlSelect(query string) string {
	builder := new(strings.Builder)
	rows, err := conn.Query(context.Background(), query)

	if err != nil {
		builder.WriteString("<p style=\"font-weight: bold; text-align: center; width: 100%;\">ОШИБКА: " + err.Error() + "</p>")
	} else if notSafe := isSafe(query); notSafe != "" {
		builder.WriteString("<p style=\"font-weight: bold; text-align: center; width: 100%;\">ОПАСНО: " + notSafe + "</p>")
	} else {
		builder.WriteString("<p style=\"font-weight: bold; text-align: center; width: 100%;\">РЕЗУЛЬТАТ:</p>")

		builder.WriteString("<table style=\"text-align: center; width: 100%;\">")

		builder.WriteString("<thead> <tr>")
		descriptions := rows.FieldDescriptions()
		for _, desc := range descriptions {
			builder.WriteString("<th>")
			builder.Write(desc.Name)
			builder.WriteString("</th>")
		}
		builder.WriteString("</tr> </thead>")

		builder.WriteString("<tbody>")
		for rows.Next() {
			columnValues, _ := rows.Values()

			builder.WriteString("<tr>")
			for _, value := range columnValues {
				builder.WriteString("<td>")
				builder.WriteString(fmt.Sprint(value))
				builder.WriteString("</td>")
			}
			builder.WriteString("</tr>")
		}
		rows.Close()
		builder.WriteString("</tbody>")

		builder.WriteString("</table>")
	}

	return builder.String()
}

func isSafe(query string) string {
	tree, err := pgquery.Parse(query)
	if err != nil {
		return fmt.Sprintf("Синтаксическая ошибка в запросе: %v", err)
	}

	for _, stmt := range tree.GetStmts() {
		node := stmt.GetStmt()

		if node.GetDropStmt() != nil {
			return "Операция DROP запрещена. Она безвозвратно удаляет объекты базы данных."
		}

		if node.GetTruncateStmt() != nil {
			return "Операция TRUNCATE запрещена. Она немедленно удаляет все данные из таблицы."
		}

		if deleteStmt := node.GetDeleteStmt(); deleteStmt != nil {
			if deleteStmt.GetWhereClause() == nil {
				return "Операция DELETE без условия WHERE запрещена. Это приведет к удалению ВСЕХ строк."
			}
		}

		if updateStmt := node.GetUpdateStmt(); updateStmt != nil {
			if updateStmt.GetWhereClause() == nil {
				return "Операция UPDATE без условия WHERE запрещена. Это приведет к обновлению ВСЕХ строк."
			}
		}
	}

	return ""
}
