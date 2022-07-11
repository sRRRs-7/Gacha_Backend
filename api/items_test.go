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

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	mockdb "github.com/sRRRs-7/GachaPon/db/mock"
	db "github.com/sRRRs-7/GachaPon/db/sqlc"
	"github.com/sRRRs-7/GachaPon/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateItemApi(t *testing.T) {
	item := randomItem()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"item_name": item.ItemName,
				"rating": item.Rating,
				"item_url": item.ItemUrl,
				"category_id": item.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateItemParams{
					ItemName: item.ItemName,
					Rating: item.Rating,
					ItemUrl: item.ItemUrl,
					CategoryID: item.CategoryID,
				}

				store.EXPECT().
					CreateItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(item, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItem(t, recorder.Body, item)
			},
		},
		{
			name: "unique_violation",
			body: gin.H{
				"item_name": item.ItemName,
				"rating": item.Rating,
				"item_url": item.ItemUrl,
				"category_id": item.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Item{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"item_name": item.ItemName,
				"rating": item.Rating,
				"item_url": item.ItemUrl,
				"category_id": item.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Item{}, sql.ErrConnDone)
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

			url := "/item/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetItemApi(t *testing.T) {
	item := randomItem()

	testCases := []struct {
		name          string
		itemID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "OK",
				itemID: item.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetItem(gomock.Any(), gomock.Eq(item.ID)).
						Times(1).
						Return(item, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchItem(t, recorder.Body, item)
				},
			},
			{
				name:      "Not Found",
				itemID: item.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetItem(gomock.Any(), gomock.Eq(item.ID)).
						Times(1).
						Return(db.Item{}, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			{
				name:      "Internal Server Error",
				itemID: item.ID,
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						GetItem(gomock.Any(), gomock.Eq(item.ID)).
						Times(1).
						Return(db.Item{}, sql.ErrConnDone)
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

				url := fmt.Sprintf("/item/get/%d", tc.itemID)
				req, err := http.NewRequest("GET", url, nil)
				require.NoError(t, err)

				server.router.ServeHTTP(recorder, req)
				tc.checkResponse(t, recorder)
			})
		}
}

func TestListItemByCategoryIdAPI(t *testing.T) {
	items := randomItems()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
					"category_id": items[0].CategoryID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListItemByCategoryIdParams{
					CategoryID: items[0].CategoryID,
					Limit:  int32(10),
					Offset: 0,
				}

				store.EXPECT().
					ListItemByCategoryId(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(items, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItems(t, recorder.Body, items)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"category_id": items[0].CategoryID,
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemByCategoryId(gomock.Any(), gomock.Any()).
					Times(1).
					Return(items, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"category_id": items[0].CategoryID,
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemByCategoryId(gomock.Any(), gomock.Any()).
					Times(1).
					Return(items, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H{
				"category_id": items[0].CategoryID,
				"page_id": int32(0),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemByCategoryId(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H{
				"category_id": items[0].CategoryID,
				"page_id": int32(1),
				"page_size": int32(-10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemByCategoryId(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/item/listByCategoryId"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListItemsByCategoryIdAPI(t *testing.T) {
	items := randomItems()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
					"page_id": int32(1),
					"page_size": int32(10),
				},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListItemsByCategoryIdParams{
					Limit:  int32(10),
					Offset: 0,
				}

				store.EXPECT().
					ListItemsByCategoryId(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(items, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItems(t, recorder.Body, items)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByCategoryId(gomock.Any(), gomock.Any()).
					Times(1).
					Return(items, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByCategoryId(gomock.Any(), gomock.Any()).
					Times(1).
					Return(items, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H{
				"page_id": int32(0),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByCategoryId(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(0),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByCategoryId(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/item/listByCategoriesId"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListItemsByIdAPI(t *testing.T) {
	items := randomItems()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
					"page_id": int32(1),
					"page_size": int32(10),
				},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListItemsByIdParams{
					Limit:  int32(10),
					Offset: 0,
				}

				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(items, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItems(t, recorder.Body, items)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Item{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Item{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H{
				"page_id": int32(0),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(0),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsById(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/item/listById"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListItemsByItemNameAPI(t *testing.T) {
	items := randomItems()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
					"page_id": int32(1),
					"page_size": int32(10),
				},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListItemsByItemNameParams{
					Limit:  int32(10),
					Offset: 0,
				}

				store.EXPECT().
					ListItemsByItemName(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(items, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItems(t, recorder.Body, items)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByItemName(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Item{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByItemName(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Item{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H{
				"page_id": int32(0),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByItemName(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(0),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByItemName(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/item/listByItemName"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListItemsByRatingAPI(t *testing.T) {
	items := randomItems()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
					"page_id": int32(1),
					"page_size": int32(10),
				},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListItemsByRatingParams{
					Limit:  int32(10),
					Offset: 0,
				}

				store.EXPECT().
					ListItemsByRating(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(items, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItems(t, recorder.Body, items)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByRating(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Item{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByRating(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Item{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidPageID",
			body: gin.H{
				"page_id": int32(0),
				"page_size": int32(10),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByRating(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPageSize",
			body: gin.H{
				"page_id": int32(1),
				"page_size": int32(0),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListItemsByRating(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/item/listByRating"
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateItemAPI(t *testing.T) {
	item := randomItem()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"id": item.ID,
				"item_name": item.ItemName,
				"rating": item.Rating,
				"item_url": item.ItemUrl,
				"category_id": item.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateItemParams{
					ID: item.ID,
					ItemName: item.ItemName,
					Rating: item.Rating,
					ItemUrl: item.ItemUrl,
					CategoryID: item.CategoryID,
				}

				store.EXPECT().
					UpdateItem(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(item, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchItem(t, recorder.Body, item)
			},
		},
		{
			name: "Not Found",
			body: gin.H{
				"id": item.ID,
				"item_name": item.ItemName,
				"rating": item.Rating,
				"item_url": item.ItemUrl,
				"category_id": item.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(item, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"id": item.ID,
				"item_name": item.ItemName,
				"rating": item.Rating,
				"item_url": item.ItemUrl,
				"category_id": item.CategoryID,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(item, sql.ErrConnDone)
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

			url := "/item/update"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteItemAPI(t *testing.T) {
	item := randomItem()

	testCases := []struct {
		name          string
		itemID        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			itemID: item.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteItem(gomock.Any(), gomock.Eq(item.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Not Found",
			itemID: item.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InternalError",
			itemID: item.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteItem(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
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

			// start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/item/delete/%d", tc.itemID)
			req, err := http.NewRequest("DELETE", url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func randomItem() db.Item {
	return db.Item{
		ID: utils.RandomInt(1, 10),
		ItemName: utils.RandomString(5),
		Rating: int32(utils.RandomInt(1, 7)),
		ItemUrl: fmt.Sprintf("http://%s", utils.RandomString(10)),
		CategoryID: int32(utils.RandomInt(1, 10)),
	}
}

func randomItems() []db.Item {
	n := 5
	items := make([]db.Item, n)
	for i := 0; i < n; i++ {
		items[i] = db.Item{
			ID: utils.RandomInt(1, 10),
			ItemName: utils.RandomString(5),
			Rating: int32(utils.RandomInt(1, 7)),
			ItemUrl: fmt.Sprintf("http://%s", utils.RandomString(10)),
			CategoryID: int32(utils.RandomInt(1, 10)),
		}
	}
	return items
}

func requireBodyMatchItem(t *testing.T, body *bytes.Buffer, item db.Item) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotItem db.Item
	err = json.Unmarshal(data, &gotItem)
	require.NoError(t, err)
	require.Equal(t, item, gotItem)
}

func requireBodyMatchItems(t *testing.T, body *bytes.Buffer, items []db.Item) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotItems []db.Item
	err = json.Unmarshal(data, &gotItems)
	require.NoError(t, err)
	require.Equal(t, items, gotItems)
}