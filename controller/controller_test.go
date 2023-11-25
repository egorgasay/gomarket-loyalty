package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/service/mocks"
	"io"
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
					r.On("Create", mock.Anything, model.RegisterRequest{Login: "login"}).Return(nil)
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
					r.On("Create", mock.Anything, model.RegisterRequest{Login: "login1"}).Return(nil)
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
					r.On("Create", mock.Anything, model.RegisterRequest{Login: ""}).Return(exception.ErrEnabledData)
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
					r.On("Create", mock.Anything, model.RegisterRequest{Login: ""}).Return(errors.New("unexpectedError"))
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match", RewardType: "pt", Reward: 10}).Return(nil)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match2", RewardType: "%", Reward: 10}).Return(nil)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match3", RewardType: "%"}).Return(exception.ErrEnabledData)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match3", RewardType: ""}).Return(exception.ErrEnabledData)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match3", RewardType: "$"}).Return(exception.ErrEnabledData)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match4", Reward: -1, RewardType: "%"}).Return(exception.ErrEnabledData)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match34", Reward: 1, RewardType: "%"}).Return(exception.ErrAlreadyExists)
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
					r.On("AddMechanic", mock.Anything, model.Mechanic{Match: "match34", Reward: 1, RewardType: "%"}).Return(errors.New("unexpectedError"))
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

func TestController_CreateOrder(t *testing.T) {
	type mckS func(r *mocks.Service)
	req := func(t []byte) *http.Request {
		return httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/orders?order_id=234&client_id=234", bytes.NewBuffer(t))
	}

	type fields struct {
		items model.Items
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
		fields   fields
		wantCode int
	}{
		{
			name: "ValidData1",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 234234324324, Price: 10, Count: 1},
							{Id: 1235, Price: 13457623, Count: 345},
							{Id: 235, Price: 323456, Count: 13464664},
						},
					}).Return(nil)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 234234324324, Price: 10, Count: 1},
						{Id: 1235, Price: 13457623, Count: 345},
						{Id: 235, Price: 323456, Count: 13464664},
					},
				},
			},

			wantCode: 200,
		},
		{
			name: "ValidData2",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 4634324324, Price: 2134, Count: 235},
						},
					}).Return(nil)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 4634324324, Price: 2134, Count: 235},
					},
				},
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
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 23421, Price: 0, Count: 1},
						},
					}).Return(exception.ErrEnabledData)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 23421, Price: 0, Count: 1},
					},
				},
			},

			wantCode: 400,
		},
		{
			name: "emtyRewardCount",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 2324324, Price: 7, Count: 0},
						},
					}).Return(exception.ErrEnabledData)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 2324324, Price: 7, Count: 0},
					},
				},
			},

			wantCode: 400,
		},
		{
			name: "enegativeRewardPrice",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 324, Price: -7, Count: 0},
						},
					}).Return(exception.ErrEnabledData)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 324, Price: -7, Count: 0},
					},
				},
			},

			wantCode: 400,
		},
		{
			name: "negativeRewardCount",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 224, Price: 7, Count: -2},
						},
					}).Return(exception.ErrEnabledData)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 224, Price: 7, Count: -2},
					},
				},
			},

			wantCode: 400,
		},
		{
			name: "errAlreadyExists",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 24324, Price: 2, Count: 1},
						},
					}).Return(exception.ErrAlreadyExists)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 24324, Price: 2, Count: 1},
					},
				},
			},

			wantCode: 409,
		},
		{
			name: "unexpectedError",
			args: args{
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("CreateOrder", mock.Anything, "234", "234", model.Items{
						Items: []model.Item{
							{Id: 234213, Price: 2, Count: 1},
						},
					}).Return(errors.ErrUnsupported)
				},
				t: model.Items{
					Items: []model.Item{
						{Id: 234213, Price: 2, Count: 1},
					},
				},
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

			app.Post("/v1/orders", controller.CreateOrder)
			resp, err := app.Test(tt.args.r)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantCode, resp.StatusCode)

		})
	}
}

func TestController_GetInfoOrders(t *testing.T) {
	type mckS func(r *mocks.Service)

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
		wantRes  string
	}{
		{
			name: "positiveTest1",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/v1/orders?client_id=234", nil),
				m: func(r *mocks.Service) {
					r.On("GetInfoOrders", mock.Anything, "234").Return([]model.Order{
						{Order: "233", Bonus: 2, Time: "2020-12-10T15:15:45+03:10"},
						{Order: "213", Bonus: 2, Time: "2020-12-10T15:15:45+03:03"},
						{Order: "23", Bonus: 2, Time: "2020-12-10T15:15:45+03:00"},
					}, nil)
				},
			},
			wantRes:  "[{\"number\":\"233\",\"accrual\":2,\"upload_time\":\"2020-12-10T15:15:45+03:10\"},{\"number\":\"213\",\"accrual\":2,\"upload_time\":\"2020-12-10T15:15:45+03:03\"},{\"number\":\"23\",\"accrual\":2,\"upload_time\":\"2020-12-10T15:15:45+03:00\"}]",
			wantCode: 200,
		},
		{
			name: "negativeTest1",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/v1/orders?client_d=234", nil),
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {},
			},
			wantCode: 400,
		},
		{
			name: "negativeTest2",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/v1/orders?client_id=234", nil),
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("GetInfoOrders", mock.Anything, "234").Return([]model.Order{
						{Order: "2313", Bonus: 2, Time: "2020-12-10T15:15:45+03:00"},
					}, exception.ErrNotFound)
				},
			},
			wantCode: 204,
		},
		{
			name: "negativeTest3",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/v1/orders?client_id=234", nil),
				w: &httptest.ResponseRecorder{},
				m: func(r *mocks.Service) {
					r.On("GetInfoOrders", mock.Anything, "234").Return([]model.Order{
						{Order: "2313", Bonus: 2, Time: "2020-12-10T15:15:45+03:00"},
					}, errors.ErrUnsupported)
				},
			},
			wantCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			app := fiber.New()

			logic := mocks.NewService(t)
			tt.args.m(logic)
			tt.args.r.Header.Set("Content-Type", "text/plain")
			tt.args.r.Header.Set("Accept", "application/json")

			controller := &Controller{
				service: logic,
			}

			app.Get("/v1/orders", controller.GetInfoOrders)
			resp, err := app.Test(tt.args.r)
			c, _ := io.ReadAll(resp.Body)
			assert.Equal(t, tt.wantRes, string(c))
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantCode, resp.StatusCode)

		})
	}
}
