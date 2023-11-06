package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController_Create(t *testing.T) {

	type mckS func(r *mocks.Service)
	req := func(t []byte) *http.Request {
		return httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user", bytes.NewBuffer(t))
	}

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		m mckS
		t any
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "positiveTest1",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("Create", model.RegisterRequest{Login: "login"}).Return(nil)
				},
				t: model.RegisterRequest{Login: "login"},
			},
			wantCode: 200,
		},
		{
			name: "positiveTest2",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("Create", model.RegisterRequest{Login: "login1"}).Return(nil)
				},
				t: model.RegisterRequest{Login: "login1"},
			},
			wantCode: 200,
		},
		{
			name: "badJSON",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {},
				t: `Logn: "dg"`,
			},

			wantCode: 400,
		},
		{
			name: "emtyLogin",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("Create", model.RegisterRequest{Login: ""}).Return(exception.ErrEnabledData)
				},
				t: model.RegisterRequest{Login: ""},
			},

			wantCode: 400,
		},
		{
			name: "unexpectedError",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("Create", model.RegisterRequest{Login: ""}).Return(errors.New("unexpectedError"))
				},
				t: model.RegisterRequest{Login: ""},
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			app := fiber.New()

			logic := mocks.NewService(t)
			tt.args.m(logic)

			body, _ := json.Marshal(&tt.args.t)
			tt.args.r = req(body)

			tt.args.r.Header.Set("Content-Type", "application/json")
			controller := &Controller{
				service: logic,
			}

			app.Post("/v1/user", controller.Create)
			resp, err := app.Test(tt.args.r)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantCode, resp.StatusCode)

		})
	}
}

func TestController_RegisterMechanic(t *testing.T) {
	type mckS func(r *mocks.Service)
	req := func(t []byte) *http.Request {
		return httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/mechanics", bytes.NewBuffer(t))
	}

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		m mckS
		t any
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "positiveTest1",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match", RewardType: "pt", Reward: 10}).Return(nil)
				},
				t: model.Mechanic{Match: "match", RewardType: "pt", Reward: 10},
			},
			wantCode: 200,
		},
		{
			name: "positiveTest2",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match2", RewardType: "%", Reward: 10}).Return(nil)
				},
				t: model.Mechanic{Match: "match2", RewardType: "%", Reward: 10},
			},
			wantCode: 200,
		},
		{
			name: "badJSON",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {},
				t: `Logn: "dg"`,
			},

			wantCode: 400,
		},
		{
			name: "emtyReward",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match3", RewardType: "%"}).Return(exception.ErrEnabledData)
				},
				t: model.Mechanic{Match: "match3", RewardType: "%"},
			},

			wantCode: 400,
		},
		{
			name: "emtyRewardType",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match3", RewardType: ""}).Return(exception.ErrEnabledData)
				},
				t: model.Mechanic{Match: "match3", RewardType: ""},
			},

			wantCode: 400,
		},
		{
			name: "enegativeRewardType",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match3", RewardType: "$"}).Return(exception.ErrEnabledData)
				},
				t: model.Mechanic{Match: "match3", RewardType: "$"},
			},

			wantCode: 400,
		},
		{
			name: "negativeReward",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match4", Reward: -1, RewardType: "%"}).Return(exception.ErrEnabledData)
				},
				t: model.Mechanic{Match: "match4", Reward: -1, RewardType: "%"},
			},

			wantCode: 400,
		},
		{
			name: "errAlreadyExists",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match34", Reward: 1, RewardType: "%"}).Return(exception.ErrAlreadyExists)
				},
				t: model.Mechanic{Match: "match34", Reward: 1, RewardType: "%"},
			},

			wantCode: 409,
		},
		{
			name: "unexpectedError",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("AddMechanic", model.Mechanic{Match: "match34", Reward: 1, RewardType: "%"}).Return(errors.New("unexpectedError"))
				},
				t: model.Mechanic{Match: "match34", Reward: 1, RewardType: "%"},
			},
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			app := fiber.New()

			logic := mocks.NewService(t)
			tt.args.m(logic)

			body, _ := json.Marshal(&tt.args.t)
			tt.args.r = req(body)

			tt.args.r.Header.Set("Content-Type", "application/json")
			controller := &Controller{
				service: logic,
			}

			app.Post("/v1/mechanics", controller.RegisterMechanic)
			resp, err := app.Test(tt.args.r)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantCode, resp.StatusCode)

		})
	}
}
