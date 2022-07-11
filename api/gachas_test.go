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

func TestCreateGachaAPI(t *testing.T) {
	user, _ := randomUser(t)
	gacha := randomGacha()
	item := randomItem()
	items := randomItems()
	gallery := randomGallery()

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
				"id": gacha.ID,
				"account_id": gacha.AccountID,
				"item_id": gacha.ItemID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg0 := db.ListItemsByIdParams{
					Limit: int32(1000),
					Offset: 0,
				}

				arg1 := db.CreateGachaParams{
					AccountID: gacha.AccountID,
					ItemID: item.ID,
				}

				arg2 := db.CreateGalleryParams{
					OwnerID: gacha.AccountID,
					ItemID: gacha.ItemID,
				}

				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Eq(arg0)).
					Times(1).
					Return(items, nil)

				store.EXPECT().
					GetItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(item, nil)

				store.EXPECT().
					CreateGacha(gomock.Any(), gomock.Eq(arg1)).
					Times(1).
					Return(gacha, nil)

				store.EXPECT().
					CreateGallery(gomock.Any(), gomock.Eq(arg2)).
					Times(1).
					Return(gallery, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGallery(t, recorder.Body, gallery)
			},
		},
		{
			name: "NoAuthorization",
			body: gin.H{
				"id": gacha.ID,
				"account_id": gacha.AccountID,
				"item_id": gacha.ItemID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateGacha(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id": gacha.ID,
				"account_id": gacha.AccountID,
				"item_id": gacha.ItemID,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg0 := db.ListItemsByIdParams{
					Limit: int32(1000),
					Offset: 0,
				}

				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Eq(arg0)).
					Times(1).
					Return(items, nil)

				store.EXPECT().
					GetItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(item, nil)

				store.EXPECT().
					CreateGacha(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Gacha{}, sql.ErrConnDone)
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/gacha/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetGachaApi(t *testing.T) {
	user, _ := randomUser(t)
	gacha := randomGacha()

	testCases := []struct {
		name          	string
		gachaID     	int64
		setupAuth     	func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    	func(store *mockdb.MockStore)
		checkResponse 	func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "OK",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Eq(gacha.ID)).
						Times(1).
						Return(gacha, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchGacha(t, recorder.Body, gacha)
				},
			},
			{
				name:      "No Authorization",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Unsupported authentication",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "unsupported", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Eq(gacha.ID)).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Invalid Authorization Format",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Eq(gacha.ID)).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Expired Token",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, -time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Eq(gacha.ID)).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Not Found",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Eq(gacha.ID)).
						Times(1).
						Return(db.Gacha{}, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			{
				name:      "Invalid ID",
				gachaID: 0,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
			{
				name:      "Internal Server Error",
				gachaID: gacha.ID,
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetGacha(gomock.Any(), gomock.Eq(gacha.ID)).
						Times(1).
						Return(db.Gacha{}, sql.ErrConnDone)
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

				url := fmt.Sprintf("/gacha/get/%d", tc.gachaID)
				req, err := http.NewRequest("GET", url, nil)
				require.NoError(t, err)

				tc.setupAuth(t, req, server.tokenMaker)
				server.router.ServeHTTP(recorder, req)
				tc.checkResponse(t, recorder)
			})
		}
}

func TestListGachasAPI(t *testing.T) {
	user, _ := randomUser(t)

	n := 10
	gachas := make([]db.Gacha, n)
	for i := 0; i < n; i++ {
		gachas[i] = randomGacha()
	}

	type Query struct {
		pageID   int
		pageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListGachasParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListGachas(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(gachas, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchGachas(t, recorder.Body, gachas)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListGachas(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				pageID:   1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListGachas(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Gacha{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			query: Query{
				pageID:   -1,
				pageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListGachas(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			query: Query{
				pageID:   1,
				pageSize: -1,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListGachas(gomock.Any(), gomock.Any()).
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

			url := "/gacha/list"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// Add query parameters to request URL
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomGacha() db.Gacha {
	return db.Gacha{
		ID: utils.RandomInt(1, 100),
		AccountID: utils.RandomInt(1, 10),
		ItemID: utils.RandomInt(1, 10),
	}
}

func requireBodyMatchGallery(t *testing.T, body *bytes.Buffer, gallery db.Gallery) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGallery db.Gallery
	err = json.Unmarshal(data, &gotGallery)
	require.NoError(t, err)
	require.Equal(t, gallery, gotGallery)
}

func requireBodyMatchGacha(t *testing.T, body *bytes.Buffer, gacha db.Gacha) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGacha db.Gacha
	err = json.Unmarshal(data, &gotGacha)
	require.NoError(t, err)
	require.Equal(t, gacha, gotGacha)
}

func requireBodyMatchGachas(t *testing.T, body *bytes.Buffer, gachas []db.Gacha) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGachas []db.Gacha
	err = json.Unmarshal(data, &gotGachas)
	require.NoError(t, err)
	require.Equal(t, gachas, gotGachas)
}