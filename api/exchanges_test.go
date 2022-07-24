package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/sRRRs-7/GachaPon/db/mock"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/token"
	"github.com/sRRRs-7/GachaPon/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateExchangeAPI(t *testing.T) {
	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	account1 := randomAccount(user1.UserName)
	account2 := randomAccount(user2.UserName)
	item1 := randomItem()
	item2 := randomItem()
	exchange1 := randomExchange(account1, account2, item1)
	exchange2 := randomExchange(account2, account1, item2)
	gallery1 := randomGallery()
	gallery2 := randomGallery()

	exchangeResult := db.ExchangeTxResult{
		Exchange1: exchange1,
		Exchange2: exchange2,
		Gallery1: gallery1,
		Gallery2: gallery2,
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": exchange1.FromAccountID,
				"to_account_id": exchange1.ToAccountID,
				"item_id_1": item1.ID,
				"item_id_2": item2.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
					Times(1).
					Return(account1, nil)

				arg := db.ExchangeTxParams{
					FromAccountID: exchange1.FromAccountID,
					ToAccountID: exchange1.ToAccountID,
					ItemID1: item1.ID,
					ItemID2: item2.ID,
				}

				store.EXPECT().
					ExchangeTx(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(exchangeResult, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchExchangeResult(t, recorder.Body, exchangeResult)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"from_account_id": exchange1.FromAccountID,
				"to_account_id": exchange1.ToAccountID,
				"item_id_1": item1.ID,
				"item_id_2": item2.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ExchangeTx(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"from_account_id": exchange1.FromAccountID,
				"to_account_id": exchange1.ToAccountID,
				"item_id_1": item1.ID,
				"item_id_2": item2.ID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user1.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(account1, nil)

				store.EXPECT().
					ExchangeTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.ExchangeTxResult{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/exchange/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetExchangeApi(t *testing.T) {
	user, _ := randomUser(t)
	account1 := randomAccount(user.UserName)
	account2 := randomAccount(user.UserName)
	item := randomItem()
	exchange := randomExchange(account1, account2, item)

	testCases := []struct {
		name          string
		exchangeID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "OK",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Eq(exchange.ID)).
						Times(1).
						Return(exchange, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchExchange(t, recorder.Body, exchange)
				},
			},
			{
				name:      "No Authorization",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Unsupported authentication",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "unsupported", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Eq(exchange.ID)).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Invalid Authorization Format",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Eq(exchange.ID)).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Expired Token",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, -time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Eq(exchange.ID)).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Not Found",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Eq(exchange.ID)).
						Times(1).
						Return(db.Exchange{}, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			{
				name:      "Invalid ID",
				exchangeID: 0,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
			{
				name:      "Internal Server Error",
				exchangeID: exchange.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetExchange(gomock.Any(), gomock.Eq(exchange.ID)).
						Times(1).
						Return(db.Exchange{}, sql.ErrConnDone)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusInternalServerError, recorder.Code)
				},
			},
		}

		for i := range testCases {
			tc := testCases[i]

			t.Run(tc.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()

				store := mockdb.NewMockStore(ctrl)
				tc.buildStubs(store)

				// start test server and send request
				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				url := fmt.Sprintf("/exchange/get/%d", tc.exchangeID)
				req, err := http.NewRequest("GET", url, nil)
				require.NoError(t, err)

				tc.setupAuth(t, req, server.tokenMaker)
				server.router.ServeHTTP(recorder, req)
				tc.checkResponse(t, recorder)
			})
		}
}

func TestListExchangeFromAccountAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	accounts1 := make([]db.Account, n)
	accounts2 := make([]db.Account, n)
	items := make([]db.Item, n)
	exchanges := make([]db.Exchange, n)
	for i := 0; i < n; i++ {
		accounts1[i] = randomAccount(user.UserName)
		accounts2[i] = randomAccount(user.UserName)
		exchanges[i] = randomExchange(accounts1[i], accounts2[i], items[i])
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H {
				"from_account_id": 1,
				"page_id":   1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListExchangeFromAccountParams{
					FromAccountID: 1,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListExchangeFromAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(exchanges, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchExchanges(t, recorder.Body, exchanges)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H {
				"from_account_id": 1,
				"page_id":   1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeFromAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H {
				"from_account_id": 1,
				"page_id":   1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeFromAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Exchange{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H {
				"from_account_id": 1,
				"page_id":   -1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeFromAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H {
				"from_account_id": 1,
				"page_id":   1,
				"page_size": 11,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeFromAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/exchange/listFromExchange"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListExchangeToAccountAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 5
	accounts1 := make([]db.Account, n)
	accounts2 := make([]db.Account, n)
	items := make([]db.Item, n)
	exchanges := make([]db.Exchange, n)
	for i := 0; i < n; i++ {
		accounts1[i] = randomAccount(user.UserName)
		accounts2[i] = randomAccount(user.UserName)
		exchanges[i] = randomExchange(accounts1[i], accounts2[i], items[i])
	}

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H {
				"to_account_id": 1,
				"page_id":   1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListExchangeToAccountParams{
					ToAccountID: 1,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListExchangeToAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(exchanges, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchExchanges(t, recorder.Body, exchanges)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H {
				"to_account_id": 1,
				"page_id":   1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeToAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H {
				"to_account_id": 1,
				"page_id":   1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeToAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Exchange{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H {
				"to_account_id": 1,
				"page_id":   -1,
				"page_size": n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeToAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H {
				"to_account_id": 1,
				"page_id":   1,
				"page_size": 11,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListExchangeToAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			// Marshal body data to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/exchange/listToExchange"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}



func randomExchange(account1, account2 db.Account, item db.Item) db.Exchange {
	return db.Exchange{
	    ID: utils.RandomInt(1, 100),
    	FromAccountID: account1.ID,
    	ToAccountID: account2.ID,
    	ItemID: item.ID,
	}
}

func requireBodyMatchExchangeResult(t *testing.T, body *bytes.Buffer, result db.ExchangeTxResult) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotResult db.ExchangeTxResult
	err = json.Unmarshal(data, &gotResult)
	require.NoError(t, err)
	require.Equal(t, result, gotResult)
}

func requireBodyMatchExchange(t *testing.T, body *bytes.Buffer, exchange db.Exchange) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotExchange db.Exchange
	err = json.Unmarshal(data, &gotExchange)
	require.NoError(t, err)
	require.Equal(t, exchange, gotExchange)
}

func requireBodyMatchExchanges(t *testing.T, body *bytes.Buffer, exchanges []db.Exchange) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotExchanges []db.Exchange
	err = json.Unmarshal(data, &gotExchanges)
	require.NoError(t, err)
	require.Equal(t, exchanges, gotExchanges)
}