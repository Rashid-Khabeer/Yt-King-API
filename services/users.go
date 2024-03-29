package services

import (
	"backend/helpers"
	"backend/models"
	"database/sql"
	"time"
)

type Users interface {
	ReadAllUsers() ([]*models.User, error)
	ReadOneUser(id int) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
}

func NewUsers() Users {
	return &userService{db: helpers.GetDB()}
}

func (n *userService) ReadAllUsers() ([]*models.User, error) {
	result, err := n.db.Query("select * from users")
	if err != nil {
		return make([]*models.User, 0), err
	}
	defer result.Close()
	var users []*models.User
	for result.Next() {
		row := models.User{}
		result.Scan(&row.Id, &row.Name, &row.Email, &row.Image, &row.TotalCoins, &row.PremiumType, &row.HasPremium, &row.LastDate, &row.Password, &row.RememberToken, &row.CreatedAt, &row.UpdatedAt, &row.AppVersion, &row.IsBlocked, &row.BlockedDays)
		users = append(users, &row)
	}
	return users, nil
}

func (n *userService) ReadOneUser(id int) (*models.User, error) {
	result, err := n.db.Query("select * from users where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	row := models.User{}
	for result.Next() {
		result.Scan(&row.Id, &row.Name, &row.Email, &row.Image, &row.TotalCoins, &row.PremiumType, &row.HasPremium, &row.LastDate, &row.Password, &row.RememberToken, &row.CreatedAt, &row.UpdatedAt, &row.AppVersion, &row.IsBlocked, &row.BlockedDays)
	}
	return &row, nil
}

func (n *userService) CreateUser(user *models.User) (*models.User, error) {
	query := "select * from users where email = ?"
	stmt := n.db.QueryRow(query, user.Email)
	var hasFound bool
	switch err := stmt.Scan(&user.Id, &user.Name, &user.Email, &user.Image, &user.TotalCoins, &user.PremiumType, &user.HasPremium, &user.LastDate, &user.Password, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt, &user.AppVersion, &user.IsBlocked, &user.BlockedDays); err {
	case sql.ErrNoRows:
		hasFound = false
	case nil:
		hasFound = true
	default:
		return nil, err
	}
	if !hasFound {
		sql := "INSERT INTO users(name, email, image, total_coins, has_premium, created_at, updated_at, app_version) VALUES(?,?,?,?,?,?,?,?)"
		insert, err := n.db.Prepare(sql)
		if err != nil {
			return nil, err
		}
		user.TotalCoins = helpers.GetIntPointer(200)
		user.HasPremium = helpers.GetBoolPointer(false)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		response, err := insert.Exec(user.Name, user.Email, user.Image, user.TotalCoins, user.HasPremium, user.CreatedAt, user.UpdatedAt, user.AppVersion)
		if err != nil {
			return nil, err
		}
		insert.Close()
		id, err := response.LastInsertId()
		if err != nil {
			return nil, err
		}
		user.Id = helpers.GetIntPointer(int(id))
	}
	return user, nil
}

type userService struct {
	db *sql.DB
}

func (n *userService) UpdateUser(user *models.User) (*models.User, error) {
	sql := "UPDATE users SET total_coins = ?, premium_type = ?, has_premium=?, last_date=?, updated_at=?, app_version = ?, is_blocked = ?, blocked_days = ? WHERE id = ?"
	insert, err := n.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	user.UpdatedAt = time.Now()
	response, err := insert.Exec(user.TotalCoins, user.PremiumType, user.HasPremium, user.LastDate, user.UpdatedAt, user.AppVersion, user.IsBlocked, user.BlockedDays, user.Id)
	if err != nil {
		return nil, err
	}
	no, err := response.RowsAffected()
	if err != nil {
		return nil, err
	}
	if no < 1 {
		return nil, err
	}
	defer insert.Close()
	return user, nil
}
