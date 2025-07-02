package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/entity"
	"github.com/mrbelka12000/wallet_calc/pkg/gorm/postgres"
)

func TestUserUsecase_Register(t *testing.T) {
	var (
		ctrl         = gomock.NewController(t)
		userRepoMock = NewMockuserRepository(ctrl)
		ctx          = context.Background()
	)
	defer ctrl.Finish()

	db, mockSQL, err := CreateMockDB()
	if err != nil {
		t.Errorf("CreateMockDB() error = %v", err)
		return
	}

	pgGorm := &postgres.Gorm{
		DB: db,
	}

	cases := []struct {
		name string

		req domainmodel.User

		mocks func()

		wantErr         bool
		wantClientError bool
	}{
		{
			name: "ok",
			req: domainmodel.User{
				Email:    "beka.teka11@gmail.com",
				Password: "test",
				Name:     "test",
			},

			mocks: func() {
				mockSQL.ExpectBegin()

				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{Email: "beka.teka11@gmail.com"}).Return(entity.User{}, gorm.ErrRecordNotFound)
				userRepoMock.EXPECT().Create(gomock.Any(), entity.User{
					ID:        uuid.New(),
					Email:     "beka.teka11@gmail.com",
					Name:      "test",
					CreatedAt: time.Now().UTC(),
				}).Return(nil)

				mockSQL.ExpectCommit()
			},
		},
		{
			name: "error get user",
			req: domainmodel.User{
				Email:    "beka.teka11@gmail.com",
				Password: "test",
				Name:     "test",
			},

			mocks: func() {
				defer mockSQL.ExpectRollback()
				mockSQL.ExpectBegin()

				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{Email: "beka.teka11@gmail.com"}).Return(entity.User{}, assert.AnError)
			},

			wantErr: true,
		},
		{
			name: "error create user",
			req: domainmodel.User{
				Email:    "beka.teka11@gmail.com",
				Password: "test",
				Name:     "test",
			},

			mocks: func() {
				defer mockSQL.ExpectRollback()
				mockSQL.ExpectBegin()

				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{Email: "beka.teka11@gmail.com"}).Return(entity.User{}, gorm.ErrRecordNotFound)
				userRepoMock.EXPECT().Create(gomock.Any(), entity.User{
					ID:        uuid.New(),
					Email:     "beka.teka11@gmail.com",
					Name:      "test",
					CreatedAt: time.Now().UTC(),
				}).Return(assert.AnError)
			},

			wantErr: true,
		},
		{
			name: "error email validation failed",
			req: domainmodel.User{
				Email:    "beka.teka11",
				Password: "test",
				Name:     "test",
			},

			wantErr:         true,
			wantClientError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mocks != nil {
				tc.mocks()
			}

			uc := NewUserUsecase(userRepoMock, pgGorm)

			err := uc.Register(ctx, tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				var clientErr ClientError
				if errors.As(err, &clientErr) && !tc.wantClientError {
					t.Errorf("Received client error, want server error: %v", err)
				} else if !errors.As(err, &clientErr) && tc.wantClientError {
					t.Errorf("Received server error, want client error: %v", err)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestUserUsecase_Login(t *testing.T) {
	var (
		ctrl         = gomock.NewController(t)
		userRepoMock = NewMockuserRepository(ctrl)
		ctx          = context.Background()
	)
	defer ctrl.Finish()

	db, _, err := CreateMockDB()
	if err != nil {
		t.Errorf("CreateMockDB() error = %v", err)
		return
	}

	pgGorm := &postgres.Gorm{
		DB: db,
	}

	cases := []struct {
		name string

		req domainmodel.User

		mocks func()

		wantClientError bool
		wantErr         bool
	}{
		{
			name: "ok",
			req: domainmodel.User{
				Email:    "beka.teka11",
				Password: "test",
			},

			mocks: func() {
				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{
					Email: "beka.teka11",
				}).Return(entity.User{
					Password: "$2a$10$MBhVGL.VgsL4gqyV2lRQoeEYW16wmzSeX.3YGs.ZyN.0IWdPNtWRC",
				}, nil)
			},
		},
		{
			name: "error invalid password",
			req: domainmodel.User{
				Email:    "beka.teka11",
				Password: "testtest",
			},

			mocks: func() {
				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{
					Email: "beka.teka11",
				}).Return(entity.User{
					Password: "$2a$10$MBhVGL.VgsL4gqyV2lRQoeEYW16wmzSeX.3YGs.ZyN.0IWdPNtWRC",
				}, nil)
			},

			wantErr:         true,
			wantClientError: true,
		},
		{
			name: "error get user",
			req: domainmodel.User{
				Email:    "beka.teka11",
				Password: "test",
			},

			mocks: func() {
				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{
					Email: "beka.teka11",
				}).Return(entity.User{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "error user not found",
			req: domainmodel.User{
				Email:    "beka.teka11",
				Password: "test",
			},

			mocks: func() {
				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{
					Email: "beka.teka11",
				}).Return(entity.User{}, gorm.ErrRecordNotFound)
			},
			wantErr:         true,
			wantClientError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mocks != nil {
				tc.mocks()
			}

			uc := NewUserUsecase(userRepoMock, pgGorm)
			_, err := uc.Login(ctx, tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				var clientErr ClientError
				if errors.As(err, &clientErr) && !tc.wantClientError {
					t.Errorf("Received client error, want server error: %v", err)
				} else if !errors.As(err, &clientErr) && tc.wantClientError {
					t.Errorf("Received server error, want client error: %v", err)
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestUserUsecase_Get(t *testing.T) {
	var (
		ctrl         = gomock.NewController(t)
		userRepoMock = NewMockuserRepository(ctrl)
		ctx          = context.Background()
	)
	defer ctrl.Finish()

	db, _, err := CreateMockDB()
	if err != nil {
		t.Errorf("CreateMockDB() error = %v", err)
		return
	}

	pgGorm := &postgres.Gorm{
		DB: db,
	}

	cases := []struct {
		name string

		req domainmodel.User

		mocks func()

		wantErr bool
	}{
		{
			name: "ok",
			req: domainmodel.User{
				Email: "beka.teka11",
			},
			mocks: func() {
				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{
					Email: "beka.teka11",
				}).Return(entity.User{}, nil)
			},
		},
		{
			name: "error get user",
			req: domainmodel.User{
				Email: "beka.teka11",
			},
			mocks: func() {
				userRepoMock.EXPECT().Get(gomock.Any(), entity.User{
					Email: "beka.teka11",
				}).Return(entity.User{}, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.mocks != nil {
				tc.mocks()
			}

			uc := NewUserUsecase(userRepoMock, pgGorm)
			_, err := uc.Get(ctx, tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
