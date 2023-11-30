package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/khafizullokh02/simplebank/db/mock"
	db "github.com/khafizullokh02/simplebank/db/sqlc"
	"github.com/stretchr/testify/require"
)

func requireBodyMatchTransferResult(t *testing.T, body io.Reader, expected db.Transfer) {
	transfer := testQueries.CreateTransfer(context.Background(), arg)
	var arg db.Transfer
	err := json.NewDecoder(body).Decode(&arg)
	require.NoError(t, err)

	require.Equal(t, expected.ID, arg.ID)
	require.Equal(t, expected.FromAccountID, arg.FromAccountID)
	require.Equal(t, expected.ToAccountID, arg.ToAccountID)
	require.Equal(t, expected.Amount, arg.Amount)
}

func TestCreateTransferAPI(t *testing.T) {
	testCases := []struct {
		name          string
		request       transferRequest
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			request: transferRequest{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "USD",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetTransfer{
						FromAccount: db.Account{
							ID:       1,
							Owner:    "Alice",
							Currency: "USD",
							Balance:  900,
						},
						ToAccount: db.Account{
							ID:       2,
							Owner:    "Bob",
							Currency: "USD",
							Balance:  1100,
						},
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTransferResult(t, recorder.Body, db.TransferResult{
					FromAccount: db.Account{
						ID:       1,
						Owner:    "Alice",
						Currency: "USD",
						Balance:  900,
					},
					ToAccount: db.Account{
						ID:       2,
						Owner:    "Bob",
						Currency: "USD",
						Balance:  1100,
					},
				})
			},
		},
		{
			name: "BadRequest",
			request: transferRequest{
				FromAccountID: 0,
				ToAccountID:   0,
				Amount:        0,
				Currency:      "",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Forbidden",
			request: transferRequest{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				Currency:      "invalid currency",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					TransferTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			server := NewServer(store)
			tc.buildStubs(store)
			server.store = store

			router := gin.Default()
			router.POST("/transfer", server.createTransfer)

			body, _ := json.Marshal(tc.request)
			req, _ := http.NewRequest("POST", "/transfer", bytes.NewReader(body))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestValidAccount(t *testing.T) {
	testCases := []struct {
		name          string
		accountID     int64
		currency      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: 1,
			currency:  "USD",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(1))).
					Times(1).
					Return(db.Account{
						ID:       1,
						Owner:    "Alice",
						Currency: "USD",
						Balance:  1000,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "NotFound",
			accountID: 2,
			currency:  "USD",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(2))).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalServerError",
			accountID: 3,
			currency:  "USD",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(3))).
					Times(1).
					Return(db.Account{}, errors.New("internal error"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 4,
			currency:  "EUR",
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(int64(4))).
					Times(1).
					Return(db.Account{
						ID:       4,
						Owner:    "Bob",
						Currency: "USD",
						Balance:  2000,
					}, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)
			server := NewServer(store)
			tc.buildStubs(store)
			server.store = store

			router := gin.Default()
			router.GET("/account/:id", server.validAccount)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/account/%d", tc.accountID), nil)
			router.ServeHTTP(w, req)

			tc.checkResponse(t, w)
		})
	}
}
