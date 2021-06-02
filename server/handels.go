package server

import (
	"encoding/json"
	"fmt"
	"github.com/DAlconoiD/balance-service/models"
	"github.com/DAlconoiD/balance-service/storage"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}
}

func handleGetBalance(storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strId := params["id"]
		id, err := strconv.Atoi(strId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(err.Error())
			return
		}

		account, cErr := storage.GetBalance(id)
		if cErr != nil {
			http.Error(w, fmt.Sprintf("[%v]", cErr.Err.Error()), http.StatusInternalServerError)
			log.Error(cErr.Err.Error())
			return
		}
		data, err := json.Marshal(account)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Marshalling failed. [%v]", err), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}
		w.Write(data)
	}
}

func handleChangeBalance(storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(err.Error())
			return
		}

		chBR := &models.ChangeBalanceRequest{}
		if err = json.Unmarshal(data, chBR); err != nil {
			http.Error(w, fmt.Sprintf("JSON Unmarshalling failed. [%v]", err), http.StatusBadRequest)
			log.Error(err.Error())
			return
		}

		v := validator.New()
		errs := v.Struct(chBR)
		if errs != nil {
			logMsg := "Validation error(s):\n"
			fmt.Fprint(w, "Validation error(s): \n")
			for _, e := range errs.(validator.ValidationErrors) {
				w.Write([]byte(fmt.Sprintf("%v\n", e)))
				logMsg += fmt.Sprintf("[%v]\n", e)
			}
			w.WriteHeader(http.StatusBadRequest)
			log.Error(logMsg)
			return
		}

		transaction, cErr := storage.UpdateBalance(chBR)
		if cErr != nil {
			if cErr.ErrorCode != models.ErrorInsufficientFundsCode {
				http.Error(w, cErr.Err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, cErr.Err.Error(), http.StatusInternalServerError)
			log.Error(cErr.Err.Error())
			return
		}

		data, err = json.Marshal(transaction)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Marshalling failed. [%v]", err), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		w.Write(data)
	}
}

func handleTransfer(storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(err.Error())
			return
		}

		tR := &models.TransferRequest{}
		if err = json.Unmarshal(data, tR); err != nil {
			http.Error(w, fmt.Sprintf("JSON Unmarshalling failed. [%v]", err), http.StatusBadRequest)
			log.Error(err.Error())
			return
		}

		v := validator.New()
		errs := v.Struct(tR)
		if errs != nil {
			logMsg := "Validation error(s):\n"
			fmt.Fprint(w, "Validation error(s): ")
			for _, e := range errs.(validator.ValidationErrors) {
				logMsg += fmt.Sprintf("[%v]\n", e)
				w.Write([]byte(fmt.Sprintf("%v", e)))
			}
			w.WriteHeader(http.StatusBadRequest)
			log.Error(logMsg)
			return
		}

		transaction, cErr := storage.MakeTransfer(tR)
		if cErr != nil {
			if cErr.ErrorCode != models.ErrorInsufficientFundsCode {
				http.Error(w, cErr.Err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, cErr.Err.Error(), http.StatusInternalServerError)
			log.Error(cErr.Err.Error())
			return
		}

		data, err = json.Marshal(transaction)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Marshalling failed. [%v]", err), http.StatusInternalServerError)
			log.Error(err.Error())
			return
		}

		w.Write(data)
	}
}

func handleGetTransactions(storage storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		strId := params["id"]
		id, err := strconv.Atoi(strId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(err)
			return
		}

		sorting := strings.ToLower(r.URL.Query().Get("sort"))
		if sorting != "" && sorting != models.SortBySumString && sorting != models.SortByTimeString {
			msg := fmt.Sprintf("Query param [sort] not valid: valid options are [%s], [%s]", models.SortBySumString, models.SortByTimeString)
			http.Error(w, msg, http.StatusBadRequest)
			log.Error(msg)
			return
		}
		if sorting == "" {
			sorting = models.SortByTimeString
		}

		order := strings.ToLower(r.URL.Query().Get("order"))
		if order != "" && order != models.OrderAscendingString && order != models.OrderDescendingString {
			msg := fmt.Sprintf("Query param [order] not valid: valid options are [%s], [%s]", models.OrderAscendingString, models.OrderDescendingString)
			http.Error(w, msg, http.StatusBadRequest)
			log.Error(msg)
			return
		}
		if order == "" {
			order = models.OrderAscendingString
		}

		var page int
		strPage := r.URL.Query().Get("page")
		if strPage != "" {
			page, err = strconv.Atoi(strPage)
			if err != nil {
				msg := "Query param [page] not valid: param must be integer number"
				http.Error(w, msg, http.StatusBadRequest)
				log.Error(msg)
				return
			}
			if page == 0 {
				page = 1
			}
		} else {
			page = -1
		}

		history, cErr := storage.GetTransactionHistory(id, sorting, order, page)
		if cErr != nil {
			http.Error(w, cErr.Err.Error(), http.StatusBadRequest)
			log.Error(cErr.Err.Error())
			return
		}

		if len(history) == 0 {
			w.Write([]byte("Transaction history is empty"))
			return
		}

		data, err := json.Marshal(history)
		if err != nil {
			http.Error(w, fmt.Sprintf("JSON Marshalling failed. [%v]", err), http.StatusInternalServerError)
			log.Error(err)
			return
		}

		w.Write(data)
	}
}
