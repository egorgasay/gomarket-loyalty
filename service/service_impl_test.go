package service

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/repository/mocks"
	"testing"
)

func Test_serviceImpl_Create(t *testing.T) {

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
					r.On("SetUser", mock.Anything).Return(exception.ErrAlreadyExists)
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

func Test_serviceImpl_AddMechanic(t *testing.T) {

	type mckR func(r *mocks.Repository)

	type fields struct {
		mechanic model.Mechanic
	}
	type args struct {
		m mckR
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "Valid data1",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Огненная вода",
					RewardType: "pt",
					Reward:     10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("AddMechanic", mock.Anything).Return(nil)
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "Valid data2",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Дорогой воздух",
					RewardType: "%",
					Reward:     10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("AddMechanic", mock.Anything).Return(nil)
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "invalid data1",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Дорогой воздух",
					RewardType: "",
					Reward:     10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {},
			},
			wantErr: true,
			err:     exception.ErrEnabledData,
		},
		{
			name: "invalid data2",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Дорогой воздух",
					RewardType: "awd",
					Reward:     10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {},
			},
			wantErr: true,
			err:     exception.ErrEnabledData,
		},
		{
			name: "already exists data",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Дорогой воздух",
					RewardType: "pt",
					Reward:     10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("AddMechanic", mock.Anything).Return(exception.ErrAlreadyExists)
				},
			},
			wantErr: true,
			err:     exception.ErrAlreadyExists,
		},
		{
			name: "negative data",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Дорогой воздух",
					RewardType: "%",
					Reward:     -10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {},
			},
			wantErr: true,
			err:     exception.ErrEnabledData,
		},
		{
			name: "unxpected error",
			fields: fields{
				mechanic: model.Mechanic{
					Match:      "Дорогой воздух",
					RewardType: "%",
					Reward:     10,
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("AddMechanic", mock.Anything).Return(errors.ErrUnsupported)
				},
			},
			wantErr: true,
			err:     errors.ErrUnsupported,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storage := mocks.NewRepository(t)
			tt.args.m(storage)
			service := &serviceImpl{
				repository: storage,
			}
			err := service.AddMechanic(tt.fields.mechanic)
			if (err != nil) != tt.wantErr || !errors.Is(err, tt.err) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
