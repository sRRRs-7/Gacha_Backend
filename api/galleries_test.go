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

func TestGetGalleryAPI(t *testing.T) {
	user, _ := randomUser(t)
	gallery := randomGallery()

	testCases := []struct {
		name          string
		galleryID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			galleryID: gallery.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGallery(gomock.Any(), gomock.Eq(gallery.ID)).
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
			galleryID: gallery.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGallery(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			galleryID: gallery.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetGallery(gomock.Any(), gomock.Eq(gallery.ID)).
					Times(1).
					Return(db.Gallery{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/gallery/get/%d", tc.galleryID)
			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestListGalleriesByIdApi(t *testing.T) {
	user, _ := randomUser(t)
	n := 10
	galleries := make([]db.Gallery, n)
	for i := 0; i < n; i++ {
		galleries[i] = randomGallery()
	}

	testCases := []struct {
		name          	string
		body     		gin.H
		setupAuth     	func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    	func(store *mockdb.MockStore)
		checkResponse 	func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "OK",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListGalleriesByIdParams{
						OwnerID: galleries[0].OwnerID,
						Limit: int32(10),
						Offset: 0,
					}

					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(galleries, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchGalleries(t, recorder.Body, galleries)
				},
			},
			{
				name:      "No Authorization",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Unsupported authentication",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "unsupported", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Invalid Authorization Format",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Expired Token",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, -time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Not Found",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListGalleriesByIdParams{
						OwnerID: galleries[0].OwnerID,
						Limit: int32(10),
						Offset: 0,
					}

					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(galleries, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			{
				name:      "Invalid ID",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(-1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
			{
				name:      "Internal Server Error",
				body: gin.H{
					"owner_id": galleries[0].OwnerID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListGalleriesByIdParams{
						OwnerID: galleries[0].OwnerID,
						Limit: int32(10),
						Offset: 0,
					}

					store.EXPECT().
						ListGalleriesById(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(galleries, sql.ErrConnDone)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				// Marshal body data to JSON
				data, err := json.Marshal(tc.body)
				require.NoError(t, err)

				url := "/gallery/listById"
				request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)
				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			})
		}
}

func TestListGalleriesByItemIdAPI(t *testing.T) {
	user, _ := randomUser(t)
	n := 10
	galleries := make([]db.Gallery, n)
	for i := 0; i < n; i++ {
		galleries[i] = randomGallery()
	}

	testCases := []struct {
		name          	string
		body     		gin.H
		setupAuth     	func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    	func(store *mockdb.MockStore)
		checkResponse 	func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name:      "OK",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListGalleriesByItemIdParams{
						ItemID: galleries[0].ItemID,
						Limit: int32(10),
						Offset: 0,
					}

					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(galleries, nil)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusOK, recorder.Code)
					requireBodyMatchGalleries(t, recorder.Body, galleries)
				},
			},
			{
				name:      "No Authorization",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Unsupported authentication",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "unsupported", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Invalid Authorization Format",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, "", user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Expired Token",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, -time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusUnauthorized, recorder.Code)
				},
			},
			{
				name:      "Not Found",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListGalleriesByItemIdParams{
						ItemID: galleries[0].ItemID,
						Limit: int32(10),
						Offset: 0,
					}

					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(galleries, sql.ErrNoRows)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusNotFound, recorder.Code)
				},
			},
			{
				name:      "Invalid ID",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(-1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Any()).
						Times(0)
				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, http.StatusBadRequest, recorder.Code)
				},
			},
			{
				name:      "Internal Server Error",
				body: gin.H{
					"item_id": galleries[0].ItemID,
					"page_id": int32(1),
					"page_size": int32(10),
				},
				setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
					addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.UserName, time.Minute)
				},
				buildStubs: func(store *mockdb.MockStore) {
					arg := db.ListGalleriesByItemIdParams{
						ItemID: galleries[0].ItemID,
						Limit: int32(10),
						Offset: 0,
					}

					store.EXPECT().
						ListGalleriesByItemId(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(galleries, sql.ErrConnDone)
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

				server := newTestServer(t, store)
				recorder := httptest.NewRecorder()

				// Marshal body data to JSON
				data, err := json.Marshal(tc.body)
				require.NoError(t, err)

				url := "/gallery/listByItemId"
				request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(data))
				require.NoError(t, err)

				tc.setupAuth(t, request, server.tokenMaker)
				server.router.ServeHTTP(recorder, request)
				tc.checkResponse(t, recorder)
			})
		}
}

func randomGallery() db.Gallery {
	return db.Gallery{
		ID: utils.RandomInt(1, 10),
		OwnerID: utils.RandomInt(1, 10),
		ItemID: utils.RandomInt(1, 10),
	}
}

func requireBodyMatchGalleries(t *testing.T, body *bytes.Buffer, galleries []db.Gallery) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotGalleries []db.Gallery
	err = json.Unmarshal(data, &gotGalleries)
	require.NoError(t, err)
	require.Equal(t, galleries, gotGalleries)
}