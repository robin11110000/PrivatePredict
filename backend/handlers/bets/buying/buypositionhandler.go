package buybetshandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"PrivatePredict/aleo"
	betutils "PrivatePredict/handlers/bets/betutils"
	"PrivatePredict/middleware"
	"PrivatePredict/models"
	"PrivatePredict/setup"
	"PrivatePredict/util"

	"gorm.io/gorm"
)

func PlaceBetHandler(loadEconConfig setup.EconConfigLoader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db := util.GetDB()
		user, httperr := middleware.ValidateUserAndEnforcePasswordChangeGetUser(r, db)
		if httperr != nil {
			http.Error(w, httperr.Error(), httperr.StatusCode)
			return
		}

		var betRequest models.Bet
		err := json.NewDecoder(r.Body).Decode(&betRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		bet, err := PlaceBetCore(user, betRequest, db, loadEconConfig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(bet)
	}
}

func PlaceBetCore(user *models.User, betRequest models.Bet, db *gorm.DB, loadEconConfig setup.EconConfigLoader) (*models.Bet, error) {
	if err := betutils.CheckMarketStatus(db, betRequest.MarketID); err != nil {
		return nil, err
	}

	sumOfBetFees := betutils.GetBetFees(db, user, betRequest)

	if err := checkUserBalance(user, betRequest, sumOfBetFees, loadEconConfig); err != nil {
		return nil, err
	}

	bet := models.CreateBet(user.Username, betRequest.MarketID, betRequest.Amount, betRequest.Outcome)

	if err := betutils.ValidateBuy(db, &bet); err != nil {
		return nil, err
	}

	totalCost := bet.Amount + sumOfBetFees
	user.AccountBalance -= totalCost

	if err := db.Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user balance: %w", err)
	}

	// 🔒 Fire private Aleo bet — outcome and amount stay private on-chain
	aleoTxID, err := aleo.PlacePrivateBet(bet.MarketID, bet.Amount, bet.Outcome)
	if err != nil {
		return nil, fmt.Errorf("aleo transaction failed: %w", err)
	}
	bet.AleoTxID = aleoTxID

	if err := db.Create(&bet).Error; err != nil {
		return nil, fmt.Errorf("failed to create bet: %w", err)
	}

	return &bet, nil
}

func checkUserBalance(user *models.User, betRequest models.Bet, sumOfBetFees int64, loadEconConfig setup.EconConfigLoader) error {
	appConfig := loadEconConfig()
	maximumDebtAllowed := appConfig.Economics.User.MaximumDebtAllowed

	if user.AccountBalance-betRequest.Amount-sumOfBetFees < -maximumDebtAllowed {
		return fmt.Errorf("Insufficient balance")
	}
	return nil
}