package service

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"gomarket-loyalty/constants"
	"gomarket-loyalty/exception"
	"gomarket-loyalty/model"
	"gomarket-loyalty/repository/mocks"
	mc "gomarket-loyalty/service/mocks"
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
			name: "ValidData1",
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
			name: "ValidData2",
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
			name: "invalidData1",
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
			name: "invalidData2",
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
			name: "alreadyExistsData",
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
			name: "negativeData",
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
			name: "unxpectedError",
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

func Test_serviceImpl_CreateOrder(t *testing.T) {
	type mckR func(r *mocks.Repository)
	type mckCL func(r *mc.Client)

	type fields struct {
		clientID string
		orderID  string
		order    model.Items
	}
	type args struct {
		m   mckR
		mCl mckCL
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "ValidData1",
			fields: fields{
				clientID: "1",
				orderID:  "1",
				order: model.Items{
					Items: []model.Item{
						{
							Id:    1,
							Count: 1,
							Price: 100,
						},
					},
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("GetAllMechanics").Return([]model.Mechanic{
						{Match: "Помидорка", Reward: 10, RewardType: "pt"},
						{Match: "ми", Reward: 10, RewardType: "pt"},
					}, nil)
					r.On("CreateOrder", model.Order{Order: "1", Bonus: 20}).Return(nil)
					r.On("UpdateBonusUser", "1", 20).Return(nil)
				},
				mCl: func(r *mc.Client) {
					r.On("JSONRequest", model.RequestNameItems{Offset: 0, Limit: 1000,
						Query: model.Query{Ids: []int{1}}}, &model.ResponseNameItems{}, constants.URLGETNameItems).
						Return(model.ResponseNameItems{Items: []model.ItemRes{
							{ID: 1, Name: "Помидорка", Price: 100}}}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "ValidData2",
			fields: fields{
				clientID: "2",
				orderID:  "2",
				order: model.Items{
					Items: []model.Item{
						{
							Id:    34,
							Count: 3,
							Price: 100,
						},
						{
							Id:    4,
							Count: 1,
							Price: 10,
						},
					},
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("GetAllMechanics").Return([]model.Mechanic{
						{Match: "По", Reward: 10, RewardType: "pt"},
						{Match: "ми", Reward: 10, RewardType: "%"},
					}, nil)
					r.On("CreateOrder", model.Order{Order: "2", Bonus: 70}).Return(nil)
					r.On("UpdateBonusUser", "2", 70).Return(nil)
				},
				mCl: func(r *mc.Client) {
					r.On("JSONRequest", model.RequestNameItems{Offset: 0, Limit: 1000,
						Query: model.Query{Ids: []int{34, 4}}}, &model.ResponseNameItems{}, constants.URLGETNameItems).
						Return(model.ResponseNameItems{Items: []model.ItemRes{
							{ID: 34, Name: "Помидорка", Price: 100},
							{ID: 4, Name: "Полка", Price: 10},
						}}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "ValidData3",
			fields: fields{
				clientID: "3",
				orderID:  "4",
				order: model.Items{
					Items: []model.Item{
						{
							Id:    34,
							Count: 3,
							Price: 10330,
						},
						{
							Id:    4,
							Count: 23,
							Price: 1550,
						},
						{
							Id:    234,
							Count: 234,
							Price: 13400,
						},
						{
							Id:    2333,
							Count: 12,
							Price: 10234,
						},
					},
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("GetAllMechanics").Return([]model.Mechanic{
						{Match: "По", Reward: 10, RewardType: "pt"},
						{Match: "ми", Reward: 10, RewardType: "%"},
						{Match: "си", Reward: 10, RewardType: "pt"},
						{Match: "окно", Reward: 10, RewardType: "%"},
					}, nil)
					r.On("CreateOrder", model.Order{Order: "4", Bonus: 17700}).Return(nil)
					r.On("UpdateBonusUser", "3", 17700).Return(nil)
				},
				mCl: func(r *mc.Client) {
					r.On("JSONRequest", model.RequestNameItems{Offset: 0, Limit: 1000,
						Query: model.Query{Ids: []int{34, 4, 234, 2333}}}, &model.ResponseNameItems{}, constants.URLGETNameItems).
						Return(model.ResponseNameItems{Items: []model.ItemRes{
							{ID: 34, Name: "Помидорка", Price: 10330},
							{ID: 4, Name: "поликака", Price: 1550},
							{ID: 234, Name: "сосиска", Price: 13400},
							{ID: 2333, Name: "окно", Price: 10234},
						}}, nil)
				},
			},
			wantErr: false,
		},
		{
			name: "ErrorJSONService",
			fields: fields{
				clientID: "3",
				orderID:  "4",
				order: model.Items{
					Items: []model.Item{
						{
							Id:    34,
							Count: 3,
							Price: 10330,
						},
						{
							Id:    4,
							Count: 23,
							Price: 1550,
						},
						{
							Id:    234,
							Count: 234,
							Price: 13400,
						},
						{
							Id:    2333,
							Count: 12,
							Price: 10234,
						},
					},
				},
			},
			args: args{
				m: func(r *mocks.Repository) {},
				mCl: func(r *mc.Client) {
					r.On("JSONRequest", model.RequestNameItems{Offset: 0, Limit: 1000,
						Query: model.Query{Ids: []int{34, 4, 234, 2333}}}, &model.ResponseNameItems{}, constants.URLGETNameItems).
						Return(nil, errors.ErrUnsupported)
				},
			},
			wantErr: true,
			err:     errors.ErrUnsupported,
		},
		{
			name: "ErrorAlreadyExists",
			fields: fields{
				clientID: "3",
				orderID:  "4",
				order: model.Items{
					Items: []model.Item{
						{
							Id:    34,
							Count: 3,
							Price: 10330,
						},
						{
							Id:    4,
							Count: 23,
							Price: 1550,
						},
						{
							Id:    234,
							Count: 234,
							Price: 13400,
						},
						{
							Id:    2333,
							Count: 12,
							Price: 10234,
						},
					},
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("GetAllMechanics").Return([]model.Mechanic{
						{Match: "По", Reward: 10, RewardType: "pt"},
						{Match: "ми", Reward: 10, RewardType: "%"},
						{Match: "си", Reward: 10, RewardType: "pt"},
						{Match: "окно", Reward: 10, RewardType: "%"},
					}, nil)
					r.On("CreateOrder", model.Order{Order: "4", Bonus: 17700}).Return(exception.ErrAlreadyExists)
				},
				mCl: func(r *mc.Client) {
					r.On("JSONRequest", model.RequestNameItems{Offset: 0, Limit: 1000,
						Query: model.Query{Ids: []int{34, 4, 234, 2333}}}, &model.ResponseNameItems{}, constants.URLGETNameItems).
						Return(model.ResponseNameItems{Items: []model.ItemRes{
							{ID: 34, Name: "Помидорка", Price: 10330},
							{ID: 4, Name: "поликака", Price: 1550},
							{ID: 234, Name: "сосиска", Price: 13400},
							{ID: 2333, Name: "окно", Price: 10234},
						}}, nil)
				},
			},
			wantErr: true,
			err:     exception.ErrAlreadyExists,
		},
		{
			name: "EmtyPrice",
			fields: fields{
				clientID: "3",
				orderID:  "4",
				order: model.Items{
					Items: []model.Item{
						{
							Id:    34,
							Count: 0,
							Price: 10330,
						},
						{
							Id:    4,
							Count: 23,
							Price: -1,
						},
						{
							Id:    234,
							Count: 2,
							Price: 100,
						},
						{
							Id:    2,
							Count: 12,
							Price: 100,
						},
					},
				},
			},
			args: args{
				m: func(r *mocks.Repository) {
					r.On("GetAllMechanics").Return([]model.Mechanic{
						{Match: "По", Reward: 10, RewardType: "pt"},
						{Match: "чт", Reward: 50, RewardType: "%"},
						{Match: "си", Reward: 10, RewardType: "pt"},
					}, nil)
					r.On("CreateOrder", model.Order{Order: "4", Bonus: 740}).Return(nil)
					r.On("UpdateBonusUser", "3", 740).Return(nil)
				},
				mCl: func(r *mc.Client) {
					r.On("JSONRequest", model.RequestNameItems{Offset: 0, Limit: 1000,
						Query: model.Query{Ids: []int{234, 2}}}, &model.ResponseNameItems{}, constants.URLGETNameItems).
						Return(model.ResponseNameItems{Items: []model.ItemRes{
							{ID: 234, Name: "Просика меда", Price: 100},
							{ID: 2, Name: "Почтовая книга", Price: 100},
						}}, nil)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientJSON := mc.NewClient(t)
			tt.args.mCl(clientJSON)
			storage := mocks.NewRepository(t)
			tt.args.m(storage)

			service := &serviceImpl{
				repository: storage,
				client:     clientJSON,
			}
			err := service.CreateOrder(tt.fields.clientID, tt.fields.orderID, tt.fields.order)
			if (err != nil) != tt.wantErr || !errors.Is(err, tt.err) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
