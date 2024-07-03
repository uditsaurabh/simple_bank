package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/uditsaurabh/simple_bank/db/mock"
	db "github.com/uditsaurabh/simple_bank/db/sqlc"
	"github.com/uditsaurabh/simple_bank/util"
)

type MockUserTestStruct struct {
	name          string
	username      string
	buildStubs    func(store *mockdb.MockStore)
	checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
}

func TestCreateUserApi(t *testing.T) {
	user, userReqBody := randomUser()
	testCases := []MockUserTestStruct{{
		name:     "Create user-valid input",
		username: user.Username,
		buildStubs: func(store *mockdb.MockStore) {
			store.EXPECT().CreateUser(gomock.Any(),
				gomock.Any()).Times(1).Return(user, nil)
		},
		checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
			require.Equal(t, http.StatusCreated, recorder.Code)
			requireBodyMatchUser(t, recorder.Body)
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
			url = "/user"
			data, err := json.Marshal(userReqBody)
			if err != nil {
				log.Fatal(err)
			}
			reader := bytes.NewReader(data)
			request, err := http.NewRequest(http.MethodPost, url, reader)
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser() (user db.User, userRequest CreateUserRequest) {
	Username := util.RandomOwner()
	password := util.RandomHashPassword()
	HashedPassword, _ := util.HashPassword(password)
	Email := util.RandomEmail()
	FullName := util.RandomOwner()
	PasswordChangedAt := time.Now()
	CreatedAt := time.Now()
	return db.User{
			Username:          Username,
			HashedPassword:    password,
			Email:             Email,
			FullName:          FullName,
			PasswordChangedAt: PasswordChangedAt,
			CreatedAt:         CreatedAt,
		}, CreateUserRequest{
			Username: Username,
			Password: HashedPassword,
			Email:    Email,
			FullName: FullName,
		}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotResp string
	err = json.Unmarshal(data, &gotResp)
	require.NoError(t, err)

	require.Equal(t, gotResp, "user created")
}
