package service

import (
	"github.com/stretchr/testify/mock"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/repository/mocks"
	"testing"
)

//func TestHashPassword(t *testing.T) {
//	tests := []struct {
//		password string
//		expected string
//	}{
//		{
//			password: "password123",
//			expected: "482c811da5d5b4bc6d497ffa98491e38",
//		},
//		{
//			password: "abc123",
//			expected: "e99a18c428cb38d5f260853678922e03",
//		},
//	}
//
//	service := &serviceImpl{}
//
//	for _, test := range tests {
//		result := service.HashPassword(test.password)
//		if result != test.expected {
//			t.Errorf("HashPassword(%s) = %s, expected %s", test.password, result, test.expected)
//		}
//	}
//}

func Test_serviceImpl_Register(t *testing.T) {

	type mckR func(r *mocks.Repository)

	type fields struct {
		userRequest model.RegisterRequest
	}
	type args struct {
		m mckR
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Valid registration1",
			fields: fields{
				userRequest: model.RegisterRequest{
					Login: "john",
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("SetUser", mock.Anything).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "Valid registration2",
			fields: fields{
				userRequest: model.RegisterRequest{
					Login: "johnawdawd",
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("SetUser", mock.Anything).Return(nil)
				},
			},
			wantErr: false,
		},
		{
			name: "invalid registration",
			fields: fields{
				userRequest: model.RegisterRequest{
					Login: "jo",
				},
			},
			args: args{
				m: func(r *mocks.Repository) {},
			},
			wantErr: true,
		},
		{
			name: "LoginAlreadyExists",
			fields: fields{
				userRequest: model.RegisterRequest{
					Login: "jo2ew",
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("SetUser", mock.Anything).Return(exception.ErrLoginAlreadyExists)
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewRepository(t)
			tt.args.m(storage)
			service := &serviceImpl{
				repository: storage,
			}
			err := service.Create(tt.fields.userRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_serviceImpl_ValidateDataRegister(t *testing.T) {
	tests := []struct {
		name    string
		user    model.RegisterRequest
		wantErr bool
	}{
		{
			name: "Valid registration data",
			user: model.RegisterRequest{
				Login: "john",
			},

			wantErr: false,
		},
		{
			name: "Invalid username",
			user: model.RegisterRequest{
				Login: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mocks.NewRepository(t)
			service := &serviceImpl{
				repository: storage,
			}
			if err := service.ValidateDataRegister(tt.user); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDataRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
