package api

import (
	"encoding/json"
	"inoscipta/service"
	response "inoscipta/utils"
	"net/http"
)

// GetBalanceHandler retrieves the balance for a given user.
func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		response.RespondWithHTML(w, http.StatusBadRequest, "user_id is required")
		return
	}
	balance, err := service.GetUserBalance(userID)
	if err != nil {
		response.RespondWithHTML(w, http.StatusInternalServerError, "Error retrieving balance")
		return
	}
	response.RespondWithJSON(w, http.StatusOK, GetBalanceResponse{
		UserID:  userID,
		Balance: balance,
	})
}

// AddAmountHandler adds funds to a user's account.
func AddAmountHandler(w http.ResponseWriter, r *http.Request) {
	var body AmountOpRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := service.AddAmount(body.UserID, body.Amount)
	if err != nil {
		response.RespondWithHTML(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.RespondWithJSON(w, http.StatusOK, "amount added successfully")
}

// DeductAmountHandler deducts funds from a user's account.
func DeductAmountHandler(w http.ResponseWriter, r *http.Request) {
	var req AmountOpRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := service.DeductAmount(req.UserID, req.Amount)
	if err != nil {
		response.RespondWithHTML(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.RespondWithJSON(w, http.StatusOK, "amount deducted successfully")
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req AmountOpRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := service.CreateAccount(req.UserID, req.Amount)
	if err != nil {
		response.RespondWithHTML(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	response.RespondWithJSON(w, http.StatusOK, "account created successfully")

}
