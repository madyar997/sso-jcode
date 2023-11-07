package usecase

import (
	"context"
	"errors"
	"github.com/madyar997/sso-jcode/config"
	"github.com/madyar997/sso-jcode/internal/entity"
	"github.com/madyar997/sso-jcode/internal/mocks/repomocks"
	"github.com/madyar997/sso-jcode/internal/usecase/repo"
	"github.com/madyar997/sso-jcode/pkg/logger"
	"reflect"
	"testing"
)

func TestUser_GetUserByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "success: user exists",
			args: args{id: 1},
			want: &entity.User{
				Id:    1,
				Email: "madiar",
			},
			wantErr: false,
		},
		{
			name:    "fail: user not exists",
			args:    args{id: 1},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//ctx := context.Background()
			repomock := &repomocks.IUserRepo{}
			u := &User{
				repo: repomock,
			}
			if !tt.wantErr {
				repomock.On("GetUserByID", context.TODO(), tt.args.id).Return(tt.want, nil)
			} else {
				repomock.On("GetUserByID", context.TODO(), tt.args.id).Return(nil, errors.New("error"))
			}

			got, err := u.GetUserByID(context.TODO(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_GetUserByEmail(t *testing.T) {
	type fields struct {
		cfg    *config.Config
		repo   repo.IUserRepo
		logger *logger.Logger
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				cfg:    tt.fields.cfg,
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
			}
			got, err := u.GetUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserByEmail() got = %v, want %v", got, tt.want)
			}
		})
	}
}
