package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

func TestCheck(t *testing.T) {
	str := `
('63b99e15-1030-4c74-82c8-2d9bb69a2554', 'Groceries', 'expense'),
('19d354ad-ec0b-495f-bbd8-faeec8b3317c', 'Restaurants', 'expense'),
('daafc4c3-fd58-4d3d-a5c3-6014fe146eb6', 'Coffee shops', 'expense'),
('ec97c4bb-f0ef-4c67-a4df-059ab43a2601', 'Delivery services', 'expense'),

('ab3b1e7c-c244-42c0-a385-9c1e9fda52d5', 'Fuel / Gas', 'expense'),
('c6152b06-8573-4f17-822b-9a7e80c94fc7', 'Public transport', 'expense'),
('3dd99b25-7fd8-4e9e-a0d5-633bc86f4a01', 'Ride-hailing', 'expense'),
('5f43a058-fb5e-4c5c-b343-1a62c93c6632', 'Parking', 'expense'),
('888c6501-b08f-4725-b307-f3c02b6ffcc3', 'Vehicle maintenance', 'expense'),

('46f45d8d-b70b-49c5-a9a5-cf4ed7330593', 'Rent / Mortgage', 'expense'),
('f80a3c65-2116-4f10-b268-34d738059993', 'Electricity', 'expense'),
('ed7281ab-d258-4ef4-b1a3-0b79e22d0669', 'Water', 'expense'),
('bbf44e80-d89a-4b3a-b491-841a4004e9b7', 'Gas (Utilities)', 'expense'),
('b2dfd5e9-9b87-4d9d-9052-cc024c3be0de', 'Internet / Mobile', 'expense'),

('aab9366e-bd3f-46a1-8b7c-4d02c9d28296', 'Streaming services', 'expense'),
('2523e21d-2b60-4cbf-a950-71df4b43a973', 'Movies / Events', 'expense'),
('1ebd6dbb-16c2-4640-8674-831125a6bc9e', 'Hobbies', 'expense'),
('bc66d589-0a14-41c2-a96c-0c167c9ac20c', 'Gaming', 'expense'),
('d7a1a78a-b77e-46d9-83f1-b3f49fe17d20', 'Books', 'expense'),

('781ddf3e-83dc-45be-9bcd-80c3fa652d02', 'Clothing', 'expense'),
('30948b4d-13ea-4f64-bbe5-22d36f5482c7', 'Electronics', 'expense'),
('a967ac0f-401b-4b94-b40b-44de8e3136f4', 'Cosmetics', 'expense'),
('f327d28c-fc40-43e3-a5b5-5b84f62d5074', 'Haircuts', 'expense'),
('0f32e1d4-e452-4a6a-b765-987a62aa7ac9', 'Household items', 'expense'),

('cc5fa529-fec0-44ae-90f7-59b43e73e56f', 'Pharmacy', 'expense'),
('478b95ef-bb93-474c-901a-2048c212a136', 'Clinic visits', 'expense'),
('015daed1-402e-4c9c-8034-bf5c623e26e5', 'Insurance premiums', 'expense'),
('5d23c5f2-0177-4f9b-a74a-c2c0956a1a77', 'Fitness/Gym', 'expense'),

('96466f10-d3e4-4bbf-9286-73039028b16e', 'School / University fees', 'expense'),
('e50c6461-5a15-48d3-88db-0ddf2d4e5b87', 'Online courses', 'expense'),
('8fcabbb0-bad5-4f6d-a637-e20b97b1c96c', 'Study materials', 'expense'),

('9f05b45f-7480-4003-a720-e9b45ef2726e', 'Hotels', 'expense'),
('45f94939-d862-4a93-82df-dbe2d2a002d0', 'Flights', 'expense'),
('b9b2c32b-5bbf-4be3-a368-8e85e91f8a8c', 'Travel insurance', 'expense'),
('ba2c6db1-01cb-4705-80b0-8e38907f4ee9', 'Tours', 'expense'),

('a204e993-f2f0-4d76-a10e-f29adff8cba3', 'Bank fees', 'expense'),
('03b0b5f4-0f18-4e90-9b5e-4dbfc5ff14c2', 'Loan payments', 'expense'),
('7c9e2c5f-2b7e-4139-9142-b79e8bdf7e92', 'Investments', 'expense'),

('2aa3909c-d839-4e9d-8426-1548ecb53661', 'Charity', 'expense'),
('c91a66c0-8791-49ab-993c-48251f9baf5f', 'Presents', 'expense'),
('8aeff2bb-b3e9-49f1-8df5-80569b8d89df', 'Religious donations', 'expense'),

-- Income categories
('87b17b0d-8aa2-437e-9dc5-6b9b9cdb3c12', 'Salary', 'income'),
('ce84549a-7a44-4624-984c-e4b20f8a9f4b', 'Cashback', 'income'),
('2ed36c84-6d5e-4bd1-a4e6-9bc7e471a27c', 'Bonuses', 'income'),
('d7a6f618-99df-41f2-a2b2-b20c207d9813', 'Transfers received', 'income')`

	fmt.Println(strings.ToLower(str))
}
