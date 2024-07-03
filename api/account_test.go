package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/uditsaurabh/simple_bank/db/mock"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
	"github.com/uditsaurabh/simple_bank/util"
)

type MockTestStruct struct {
	name          string
	accountID     int64
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func TestGetAccountApi(t *testing.T) {
	account := randomAccount()
	testCases := []MockTestStruct{{
		name:      "ok",
		accountID: account.ID,
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().
				GetAccount(gomock.Any(),
					gomock.Eq(account.ID)).Times(1).Return(account, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusOK, recorder.Code)
			requireBodyMatchAccount(t, recorder.Body, account)
		},
	},
		{
			name:      "No Account found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(),
						gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "Internal server error",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(),
						gomock.Eq(account.ID)).Times(1).Return(db.Account{}, errors.New("Server broke"))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid Id",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(),
						gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for i := 0; i < len(testCases); i++ {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			var url string
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			//start the test http server
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			if tc.name == "Invalid Id" {
				url = fmt.Sprintf("/account/%s", "$$")
			} else {
				url = fmt.Sprintf("/account/%d", tc.accountID)
			}
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount() db.Account {

	return db.Account{
		ID:        int64(rand.Intn(100)),
		Owner:     util.RandomOwner(),
		Balance:   util.RandomMoney(),
		Currency:  "EUR",
		CreatedAt: time.Now(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	fmt.Println(gotAccount)
	fmt.Println(reflect.DeepEqual(account, gotAccount))
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.Equal(t, account.Owner, gotAccount.Owner)
}
