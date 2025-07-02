package converter

import (
	domainmodel "github.com/mrbelka12000/wallet_calc/internal/domain_model"
	"github.com/mrbelka12000/wallet_calc/internal/dto"
	"github.com/mrbelka12000/wallet_calc/internal/entity"
)

func ConvertFromEntityUserToDTO(u entity.User, withPassword bool) dto.User {
	user := dto.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}

	if withPassword {
		user.Password = u.Password
	}

	return user
}

func ConvertFromDTOUserToDomainModel(u dto.User) domainmodel.User {
	return domainmodel.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func ConvertFromDTOUserCUToDomainModel(u dto.UserCU) domainmodel.User {
	return domainmodel.User{
		Email:    u.Email,
		Name:     u.Name,
		Password: u.Password,
	}
}
