package data

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId    int64  `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"user_name"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
	Sex       string `json:"sex"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	CreateAt  string `json:"create_at"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(user *User) error {
	log.Println(user.Username, user.Password, user.Email)
	query := `insert into users (user_name, password, email) values (?, ?, ?)`
	_, err := u.DB.Exec(query,
		user.Username,
		user.Password,
		user.Email)
	if err != nil {
		log.Println("exec insert fail: ", err)
		return err
	}
	return nil
}

func (u *UserModel) Update(user *User) error {
	query := `update users
		set first_name = ?, last_name = ?, birthday = ?, sex = ?, address = ?, phone = ?
		where user_id = ?`
	_, err := u.DB.Exec(query,
		user.FirstName,
		user.LastName,
		user.Birthday,
		user.Sex,
		user.Address,
		user.Phone,
		user.UserId)
	if err != nil {
		log.Println("update user fail ", err)
		return err
	}
	return nil
}

func (u *UserModel) Delete(id string) error {
	query := `delete from users where user_id = ?`
	_, err := u.DB.Exec(query, id)
	if err != nil {
		log.Println("delete user fail ", err)
		return err
	}
	return nil
}

func (u *UserModel) GetUserById(id float64, user *User) error {
	query := `select user_id, user_name, password, IFNULL(first_name, ""), IFNULL(last_name,""), IFNULL(birthday,""), IFNULL(sex,""), IFNULL(address,""), IFNULL(phone,""), create_at
		from users where user_id = ?`
	r := u.DB.QueryRow(query, id)
	err := r.Scan(
		&user.UserId,
		&user.Username,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Birthday,
		&user.Sex,
		&user.Address,
		&user.Phone,
		&user.CreateAt)
	if err != nil {
		log.Println("get user fail ", err)
		return err
	}
	return err
}

func (u *UserModel) CheckUserName(userName string) (bool, error) {
	rs, err := u.DB.Query(`select user_name from users where user_name = ?`, userName)
	defer rs.Close()
	if err != nil {
		log.Println("exec check username fail")
		return false, err
	}
	if rs.Next() {
		return true, nil //available user_name
	}
	return false, nil //not available
}

func (u *UserModel) ChangePass(password string, email string) (bool, error) {
	rs, err := u.DB.Exec(`update users set password = ? where email = ?`, password, email)
	if err != nil {
		log.Println("exec check username fail")
		return false, err
	}
	r, err := rs.RowsAffected()
	log.Println(r)
	if r == 0 {
		return false, err
	}
	return true, nil
}

func (u *UserModel) CheckEmail(email string) (bool, error) {
	rs, err := u.DB.Query(`select email from users where email = ?`, email)
	defer rs.Close()
	if err != nil {
		log.Println("exec check email fail")
		return false, err
	}
	if rs.Next() {
		return true, nil //available email
	}
	return false, nil //not available
}

func (u *UserModel) Login(username string, email string, password string) (User, error) {
	var user User
	query := `select user_id, user_name, email, IFNULL(first_name, ""), IFNULL(last_name,""), IFNULL(birthday,""), IFNULL(sex,""), IFNULL(address,""), IFNULL(phone,""), create_at
		from users where user_name = ? or email = ? and password = ?`
	rs, err := u.DB.Query(query, username, email, password)
	if err != nil {
		fmt.Println("exec check login fail")
		return user, err
	}
	if rs.Next() {
		err = rs.Scan(
			&user.UserId,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Birthday,
			&user.Sex,
			&user.Address,
			&user.Phone,
			&user.CreateAt)
		if err != nil {
			log.Println("scan id user login fail")
			return user, err
		}
	}
	return user, nil
}

func (u *UserModel) GetAllUser(payload Payload) ([]User, error) {
	offset := (payload.Stt - 1) * payload.Limit // bat dau la row 0
	var users []User
	query := `select user_id, user_name, email, first_name, last_name, birthday, sex, address, phone, create_at from users LIMIT ?, ?`
	rs, err := u.DB.Query(query, offset, payload.Limit)
	if err != nil {
		log.Println("exec get all users fail ", err)
		return users, err
	}
	for rs.Next() {
		var user User
		err = rs.Scan(&user.UserId, &user.Username, &user.Email, &user.FirstName, &user.LastName,
			&user.Birthday, &user.Sex, &user.Address, &user.Phone, &user.CreateAt)
		if err != nil {
			log.Println("exec scan 1 user fail ", err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}
