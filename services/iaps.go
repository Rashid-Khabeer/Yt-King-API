package services

import (
	"backend/helpers"
	"backend/models"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"
)

type Iap interface {
	ReadAllTransactions() ([]*models.Transactions, error)
	HandleRevenueCatWebhooks(notificationData []byte) (interface{}, error)
}

func NewTransactions() Iap {
	notificationService := NewNotificationService()
	return &transactionsService{db: helpers.GetDB(), notificationService: *notificationService}
}

func (n *transactionsService) ReadAllTransactions() ([]*models.Transactions, error) {
	// n.notificationService.SendNotificationToTopic("YT", "Hello Youtubers", "260")
	result, err := n.db.Query("select * from transactions")
	if err != nil {
		return make([]*models.Transactions, 0), err
	}
	defer result.Close()
	var transactions []*models.Transactions
	for result.Next() {
		row := models.Transactions{}
		var webhook string
		result.Scan(&row.Id, &row.Email, &row.Date, &row.Price, &row.TransactionId, &webhook, &row.CreatedAt, &row.UpdatedAt, &row.Type, &row.ProductId)
		var jsonWebhook map[string]interface{}
		json.Unmarshal([]byte(webhook), &jsonWebhook)
		row.WebhookResponse = &jsonWebhook
		transactions = append(transactions, &row)
	}
	return transactions, nil
}

func (s *transactionsService) HandleRevenueCatWebhooks(notificationData []byte) (interface{}, error) {
	var dat map[string]interface{}
	if err := json.Unmarshal(notificationData, &dat); err != nil {
		panic(err)
	}
	event := dat["event"].(map[string]interface{})
	sql := "INSERT INTO transactions(email, date, price, transaction_id, webhook_response, created_at, updated_at, type, product_id) VALUES(?,?,?,?,?,?,?,?,?)"
	insert, err := s.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	email := event["app_user_id"]
	date := ""
	if event["purchased_at_ms"] != nil {
		purchaseAtMs := event["purchased_at_ms"].(float64)
		date = time.Unix(0, int64(purchaseAtMs)*int64(time.Millisecond)).String()
	}
	price := event["price"]
	transactionId := event["transaction_id"]
	wr, err := json.Marshal(dat)
	webhookResponse := string(wr)
	pType := event["type"]
	productId := event["product_id"]
	createdAt := time.Now()
	updatedAt := time.Now()
	response, err := insert.Exec(email, date, price, transactionId, webhookResponse, createdAt, updatedAt, pType, productId)
	if err != nil {
		return nil, err
	}
	defer insert.Close()
	affected, err := response.RowsAffected()
	if affected < 1 {
		return nil, err
	}
	if pType == "NON_RENEWING_PURCHASE" {
		result, err := s.db.Query("select * from users where email = ?", email)
		if err != nil {
			return nil, err
		}
		defer result.Close()
		user := models.User{}
		if result.Next() {
			result.Scan(&user.Id, &user.Name, &user.Email, &user.Image, &user.TotalCoins, &user.PremiumType, &user.HasPremium, &user.LastDate, &user.Password, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt, &user.AppVersion, &user.IsBlocked, &user.BlockedDays)
		}
		coins := 0
		premiumType := ""
		hasPremium := *user.HasPremium
		switch productId {
		case "buy_coin1":
			coins = *user.TotalCoins + 10000
			premiumType = "buy_coin1"
			break
		case "buy_coin2":
			coins = *user.TotalCoins + 50000
			premiumType = "buy_coin2"
			break
		case "buy_coin3":
			coins = *user.TotalCoins + 100000
			premiumType = "buy_coin3"
			break
		case "buy_coin4":
			coins = *user.TotalCoins + 2500000
			premiumType = "buy_coin4"
			break
		case "buy_coin5":
			coins = *user.TotalCoins + 1000000
			premiumType = "buy_coin5"
			break
		case "buy_coin6":
			coins = *user.TotalCoins + 2000000
			premiumType = "buy_coin6"
			break
		case "vip1":
			hasPremium = true
			premiumType = "vip1"
			now := time.Now()
			now = now.Add(time.Hour * 24 * 7)
			lastDate := now.Format("2006-01-02 15:04:05")
			user.LastDate = &lastDate
			break
		case "vip1.5":
			hasPremium = true
			premiumType = "vip1.5"
			now := time.Now()
			now = now.Add(time.Hour * 24 * 30)
			lastDate := now.Format("2006-01-02 15:04:05")
			user.LastDate = &lastDate
			break
		case "vip3":
			hasPremium = true
			premiumType = "vip3"
			now := time.Now()
			now = now.Add(time.Hour * 24 * 30 * 12 * 100)
			lastDate := now.Format("2006-01-02 15:04:05")
			user.LastDate = &lastDate
			break
		}
		user.TotalCoins = &coins
		user.HasPremium = &hasPremium
		user.PremiumType = &premiumType
		updateQuery := "UPDATE users SET total_coins = ?, premium_type = ?, has_premium = ?, last_date=?, updated_at = ? WHERE id = ?"
		update, err := s.db.Prepare(updateQuery)
		if err != nil {
			return nil, err
		}
		user.UpdatedAt = time.Now()
		updateResponse, err := update.Exec(user.TotalCoins, user.PremiumType, user.HasPremium, user.LastDate, user.UpdatedAt, user.Id)
		if err != nil {
			return nil, err
		}
		no, err := updateResponse.RowsAffected()
		if err != nil {
			return nil, err
		}
		if no < 1 {
			return nil, err
		}
		token := strconv.Itoa(*user.Id)
		s.notificationService.SendNotificationToTopic("Purchase Successfull!", "Your purchase has been completed, kindly reopen the app.", token)
		defer update.Close()
	}
	return affected, nil
}

type transactionsService struct {
	db                  *sql.DB
	notificationService notificationService
}
