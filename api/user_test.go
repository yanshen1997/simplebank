package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/yanshen1997/simplebank/db/mock"
	db "github.com/yanshen1997/simplebank/db/sqlc"
	"github.com/yanshen1997/simplebank/util"
)

type eqCreateUserParamMatcher struct {
	password string
	args     db.CreateUserParams
}

func (e eqCreateUserParamMatcher) Matches(x interface{}) bool {

	args, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, args.HashedPassword)
	if err != nil {
		return false
	}
	e.args.HashedPassword = args.HashedPassword
	return reflect.DeepEqual(e.args, args)
}

func (e eqCreateUserParamMatcher) String() string {
	return fmt.Sprintf("is equal to %v (%T)", e.args, e.args)
}

func EqCreateUserParam(args db.CreateUserParams, passsword string) gomock.Matcher {
	return eqCreateUserParamMatcher{
		args:     args,
		password: passsword,
	}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStub: func(store *mockdb.MockStore) {
				args := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParam(args, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for _, v := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		v.buildStub(store)
		server := NewServer(store)
		recorder := httptest.NewRecorder()

		data, err := json.Marshal(v.body)
		require.NoError(t, err)

		url := "/users"
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
		require.NoError(t, err)
		server.router.ServeHTTP(recorder, req)
		v.checkResponse(t, recorder)
	}

}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	return db.User{
		Username:       util.GetRandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.GetRandomOwner(),
		Email:          util.GetRandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)

	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Empty(t, gotUser.HashedPassword)
}
