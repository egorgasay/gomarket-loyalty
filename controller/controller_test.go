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

func TestController_Register(t *testing.T) {

	type mckS func(r *mocks.Service)

	t1, _ := json.Marshal(&model.RegisterRequest{Login: "login", Password: "waxXd2dAwd3"})
	t2, _ := json.Marshal(&model.RegisterRequest{Login: "login1", Password: "awdSwd4/7"})
	t3, _ := json.Marshal(&model.RegisterRequest{Login: "", Password: "pa6Wtk7ssword"})
	t4, _ := json.Marshal(`Logn: "dg"`)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
		m mckS
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
				r: httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user/register", bytes.NewBuffer(t1)),
				m: func(r *mocks.Service) {
					r.On("Register", model.RegisterRequest{Login: "login", Password: "waxXd2dAwd3"}).Return("ZDRGDRHG-DRGDRGDRG-2TJFYTUKAW", nil)
				},
			},
			wantCode: 200,
		},
		{
			name: "positiveTest2",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user/register", bytes.NewBuffer(t2)),
				m: func(r *mocks.Service) {
					r.On("Register", model.RegisterRequest{Login: "login1", Password: "awdSwd4/7"}).Return("ZDRGDRHG-DRioRG-2TJFYTUKAW", nil)
				},
			},
			wantCode: 200,
		},
		{
			name: "badJSON",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user/register", bytes.NewBuffer(t4)),
				m: func(r *mocks.Service) {},
			},

			wantCode: 400,
		},
		{
			name: "emtyLogin",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user/register", bytes.NewBuffer(t3)),
				m: func(r *mocks.Service) {
					r.On("Register", model.RegisterRequest{Login: "", Password: "pa6Wtk7ssword"}).Return("", exception.ErrEnabledData)
				},
			},

			wantCode: 400,
		},
		{
			name: "unexpectedError",
			args: args{
				w: &httptest.ResponseRecorder{},
				r: httptest.NewRequest(http.MethodPost, "http://127.0.0.1:8080/v1/user/register", bytes.NewBuffer(t3)),
				m: func(r *mocks.Service) {
					r.On("Register", model.RegisterRequest{Login: "", Password: "pa6Wtk7ssword"}).Return("", errors.New("unexpectedError"))
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
			tt.args.r.Header.Set("Content-Type", "application/json")

			controller := &Controller{
				service: logic,
			}

			app.Post("/v1/user/register", controller.Register)
			resp, err := app.Test(tt.args.r)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantCode, resp.StatusCode)

		})
	}
}
