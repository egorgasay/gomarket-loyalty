package repository

import (
	"context"
	"fmt"
	"github.com/egorgasay/dockerdb/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gomarket-loyalty/model"
	"testing"
	"time"
)

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
			name: "Valid user",

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
			name: "Duplicate user",
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
			if err := repository.SetUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("SetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
