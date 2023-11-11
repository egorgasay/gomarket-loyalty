package repository

import (
	"context"
	"fmt"
	"github.com/egorgasay/dockerdb/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomarket-loyalty/config"
	"gomarket-loyalty/model"
	"reflect"
	"testing"
	"time"
)

type Result struct {
	client *mongo.Client
	vdb    *dockerdb.VDB
	errase func()
}

func upMongo(ctx context.Context, t *testing.T) *Result {
	var cl *mongo.Client
	var err error
	cfg := dockerdb.EmptyConfig().Vendor("mongo").DBName("golang_test").
		NoSQL(func(c dockerdb.Config) (stop bool) {
			opt := options.Client()
			dsn := fmt.Sprintf("mongodb://127.0.0.1:%s", c.GetActualPort())
			opt.ApplyURI(dsn).SetTimeout(1 * time.Second)

			cl, err = mongo.Connect(ctx, opt)
			if err != nil {
				t.Log("can't connect to mongodb")
				return false
			}

			if err := cl.Ping(ctx, nil); err != nil {
				t.Log("can't ping mongodb")
				return false
			}
			return true
		}, 30, 2*time.Second).PullImage().StandardDBPort("27017").Build()

	vdb, err := dockerdb.New(ctx, cfg)
	if err != nil {
		t.Fatalf("can't up mongo for tests %v", err)
	}

	res := new(Result)
	res.vdb = vdb
	res.client = cl
	res.errase = func() {
		err = res.vdb.Clear(ctx)
		if err != nil {
			t.Errorf("can't clear container, possible container leak and wrong results in the future tests: %v", err)
		}
	}
	return res
}

func Test_repositoryImpl_SetUser(t *testing.T) {

	type args struct {
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ValidUser",

			args: args{
				user: model.User{
					Login:      "John Doe",
					Bonus:      0,
					SpentBonus: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "DuplicateUser",
			args: args{
				user: model.User{
					Login: "John Doe",
				},
			},
			wantErr: true,
		},
	}
	res := upMongo(context.Background(), t)
	defer res.errase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repository := &repositoryImpl{
				db: res.client.Database("golang_test"),
			}
			if err := repository.SetUser(context.Background(), tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repositoryImpl_AddMechanic(t *testing.T) {
	type fields struct {
		db *mongo.Database
	}
	type args struct {
		bonus model.Mechanic
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ValidBonus1",
			args: args{
				bonus: model.Mechanic{
					Match:      "match",
					RewardType: "pt",
					Reward:     10,
				},
			},
			wantErr: false,
		},
		{
			name: "ValidBonus2",
			args: args{
				bonus: model.Mechanic{
					Match:      "match2",
					RewardType: "%",
					Reward:     10,
				},
			},
			wantErr: false,
		},
		{
			name: "alreadyExists",
			args: args{
				bonus: model.Mechanic{
					Match:      "match",
					RewardType: "%",
					Reward:     10,
				},
			},
			wantErr: true,
		},
	}
	res := upMongo(context.Background(), t)
	defer res.errase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repository := &repositoryImpl{
				db: res.client.Database("golang_test"),
			}
			if err := repository.AddMechanic(context.Background(), tt.args.bonus); (err != nil) != tt.wantErr {
				t.Errorf("Mechanic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repositoryImpl_UpdateBonusUser(t *testing.T) {
	type fields struct {
		db *mongo.Database
	}
	type args struct {
		id    string
		bonus int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantBonus int
	}{
		{
			name: "ValidBonus",
			args: args{
				id:    "John Doe",
				bonus: 12,
			},
			wantErr:   false,
			wantBonus: 12,
		},
		{
			name: "ValidBonus",
			args: args{
				id:    "John Doe",
				bonus: 12,
			},
			wantErr:   false,
			wantBonus: 24,
		},
	}
	res := upMongo(context.Background(), t)
	defer res.errase()
	func() {
		ctx, cancel := config.NewMongoContext()
		defer cancel()
		_, _ = res.client.Database("golang_test").Collection("users").InsertOne(ctx, model.User{
			Login:      "John Doe",
			Bonus:      0,
			SpentBonus: 0,
		})
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repository := &repositoryImpl{
				db: res.client.Database("golang_test"),
			}
			if err := repository.UpdateBonusUser(context.Background(), tt.args.id, tt.args.bonus); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBonusUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			func() {
				ctx, cancel := config.NewMongoContext()
				defer cancel()

				var user model.User
				err := res.client.Database("golang_test").Collection("users").FindOne(ctx, bson.M{"_id": tt.args.id}).Decode(&user)
				if err != nil {
					t.Errorf("1 UpdateBonusUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				if user.Bonus != tt.wantBonus {
					t.Errorf("3 UpdateBonusUser() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
		})
	}
}

func Test_repositoryImpl_CreateOrder(t *testing.T) {
	type fields struct {
		db *mongo.Database
	}
	type args struct {
		order model.Order
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "ValidOrder",
			args: args{
				order: model.Order{
					Order: "order",
					Bonus: 12,
				},
			},
			wantErr: false,
		},
		{name: "AlreadyExists",
			args: args{
				order: model.Order{
					Order: "order",
					Bonus: 12,
				},
			},
			wantErr: true,
		},
	}
	res := upMongo(context.Background(), t)
	defer res.errase()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repository := &repositoryImpl{
				db: res.client.Database("golang_test"),
			}

			if err := repository.CreateOrder(context.Background(), tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repositoryImpl_GetAllMechanics(t *testing.T) {
	type fields struct {
		db *mongo.Database
	}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "GetAllMechanics",
			wantErr: false,
		},
	}
	res := upMongo(context.Background(), t)
	defer res.errase()
	mechanic := []model.Mechanic{
		{Match: "match", RewardType: "pt", Reward: 10},
		{Match: "asc", RewardType: "%", Reward: 10},
		{Match: "awe", RewardType: "pt", Reward: 768},
		{Match: "awvyj", RewardType: "%", Reward: 34},
		{Match: "awdi", RewardType: "pt", Reward: 234},
	}
	func() {
		ctx, cancel := config.NewMongoContext()
		defer cancel()
		for _, mec := range mechanic {
			_, _ = res.client.Database("golang_test").Collection("mechanics").InsertOne(ctx, mec)
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repository := &repositoryImpl{
				db: res.client.Database("golang_test"),
			}

			got, err := repository.GetAllMechanics(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllMechanics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, mechanic) {
				t.Errorf("GetAllMechanics() got = %v, want %v", got, mechanic)
			}
		})
	}
}
