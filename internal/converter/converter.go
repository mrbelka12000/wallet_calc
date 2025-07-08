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

func ConvertFromEntityCategoriesToDM(categories []entity.Category) []domainmodel.Category {
	dmCategories := make([]domainmodel.Category, len(categories))
	for i, category := range categories {
		dmCategories[i] = domainmodel.Category{
			ID:   category.ID,
			Name: category.Name,
			Type: category.Type,
		}
	}
	return dmCategories
}

func ConvertFromDMCategoriesToDTO(categories []domainmodel.Category) []dto.CategoryResp {
	dmCategories := make([]dto.CategoryResp, len(categories))
	for i, category := range categories {
		dmCategories[i] = dto.CategoryResp{
			ID:   category.ID,
			Name: category.Name,
			Type: category.Type,
		}
	}
	return dmCategories
}

func ConvertFromEntityAccountsToDTO(accounts []entity.Account) []dto.AccountResp {
	dtoAccounts := make([]dto.AccountResp, len(accounts))
	for i, account := range accounts {
		dtoAccounts[i] = dto.AccountResp{
			ID:        account.ID,
			UserID:    account.UserID,
			Name:      account.Name,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
		}
	}

	return dtoAccounts
}

func ConvertFromDomainModelBtToDTO(trxs []domainmodel.BaseTransaction) []dto.BaseTransaction {
	dtoTransactions := make([]dto.BaseTransaction, len(trxs))
	for i, trx := range trxs {
		dtoTransactions[i] = dto.BaseTransaction{
			Date:        trx.Date,
			Description: trx.Description,
			Amount:      trx.Amount,
			Transaction: trx.Transaction,
			Category:    trx.Category,
			Confidence:  trx.Confidence,
		}
	}

	return dtoTransactions
}
