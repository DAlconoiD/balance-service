package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dalconoid/balance-service/models"
	mockdb "github.com/dalconoid/balance-service/storage/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestAliveHandle(t *testing.T) {
	req, _ := http.NewRequest("GET", "/hello", nil)

	rr := httptest.NewRecorder()
	handler := handleAlive()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestGetBalanceHandleStandardBehaviour(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 20
	id := rand.Intn(max-min+1) + 1
	vars := map[string]string{
		"id": strconv.Itoa(id),
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%v", id), nil)
	req = mux.SetURLVars(req, vars)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDb := mockdb.NewMockStore(mockCtrl)
	dummyAccount := models.Account{ID: id, Balance: float64(id * 100)}
	mockDb.EXPECT().GetBalance(id).Return(&dummyAccount, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := handleGetBalance(mockDb)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestChangeBalanceHandleStandardBehaviour(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	minId := 1
	maxId := 20
	id := rand.Intn(maxId-minId+1) + 1
	minDelta := 0.01
	maxDelta := 9999.99
	delta := minDelta + rand.Float64() * (maxDelta - minDelta)
	delta = math.Round(delta*100)/100
	chBR := models.ChangeBalanceRequest{ID: id, Delta: 100}
	entryData, _ := json.Marshal(chBR)
	req, _ := http.NewRequest("GET", "change-balance", bytes.NewBuffer(entryData))

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDb := mockdb.NewMockStore(mockCtrl)
	dummyTransaction := models.Transaction{
		ID:        1,
		AccountID: id,
		CreatedAt: time.Now(),
		Delta:     delta,
		Remaining: 100 + delta,
	}
	dummyTransaction.Message = fmt.Sprintf("Account [%v]: balance changed by [%.2f], [%.2f] remaining", dummyTransaction.AccountID, dummyTransaction.Delta, dummyTransaction.Remaining)
	mockDb.EXPECT().UpdateBalance(&chBR).Return(&dummyTransaction, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := handleChangeBalance(mockDb)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestTransferHandleStandardBehaviour(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	minId := 1
	maxId := 20
	id := rand.Intn(maxId-minId+1) + 1
	minDelta := 0.01
	maxDelta := 9999.99
	delta := minDelta + rand.Float64() * (maxDelta - minDelta)
	delta = math.Round(delta*100)/100
	tR := models.TransferRequest{ID1: id, ID2: id + 1, Delta: 100}
	entryData, _ := json.Marshal(tR)
	req, _ := http.NewRequest("GET", "change-balance", bytes.NewBuffer(entryData))

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDb := mockdb.NewMockStore(mockCtrl)
	dummyTransaction := models.Transaction{
		ID:        1,
		AccountID: id,
		CreatedAt: time.Now(),
		Delta:     -delta,
		Remaining: 10000 - delta,
	}
	dummyTransaction.Message = fmt.Sprintf("Transfer from account [%v] to account [%v]: balance changed by [%.2f], [%.2f] remaining",
		dummyTransaction.ID, tR.ID2, dummyTransaction.Delta, dummyTransaction.Remaining)
	mockDb.EXPECT().MakeTransfer(&tR).Return(&dummyTransaction, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := handleTransfer(mockDb)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}

func TestGetTransactionsHandleStandardBehaviour(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 20
	id := rand.Intn(max-min+1) + 1
	vars := map[string]string{
		"id": strconv.Itoa(id),
	}
	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactions/%v?", id), nil)
	req = mux.SetURLVars(req, vars)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDb := mockdb.NewMockStore(mockCtrl)
	dummyTransactions := make([]models.Transaction, 0, 15)
	for i := 1; i < 15; i++ {
		t := models.Transaction{
			ID: i,
			AccountID: id,
			CreatedAt: time.Now().Add(1*time.Hour),
			Delta: 10.00,
			Remaining: float64(i * 10),
			Message: fmt.Sprintf("Account [%v]: balance changed by [%.2f], [%.2f] remaining", id, 10.00, float64(i * 10)),
		}
		dummyTransactions = append(dummyTransactions, t)
	}
	mockDb.EXPECT().GetTransactionHistory(id, "by-time", "asc", -1).Return(dummyTransactions, nil).Times(1)

	rr := httptest.NewRecorder()
	handler := handleGetTransactions(mockDb)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
}