package impl

import (
	"sync"
	"xspends/models/interfaces"
)

var (
	ModelsService *ModelsServiceContainer
	once          sync.Once
)

type ModelsServiceContainer struct {
	DBService           *DBService
	CategoryModel       interfaces.CategoryService
	SourceModel         interfaces.SourceService
	UserModel           interfaces.UserService
	TagModel            interfaces.TagService
	TransactionTagModel interfaces.TransactionTagService
	TransactionModel    interfaces.TransactionService
}

// ModelsConfig struct to group all the dependencies
type ModelsConfig struct {
	DBService           *DBService
	CategoryModel       interfaces.CategoryService
	SourceModel         interfaces.SourceService
	UserModel           interfaces.UserService
	TagModel            interfaces.TagService
	TransactionTagModel interfaces.TransactionTagService
	TransactionModel    interfaces.TransactionService
}

// InitModelsService is now a function that takes a ModelsConfig struct
func InitModelsService(config *ModelsConfig) {
	once.Do(func() {
		ModelsService = &ModelsServiceContainer{
			DBService:           config.DBService,
			CategoryModel:       config.CategoryModel,
			SourceModel:         config.SourceModel,
			UserModel:           config.UserModel,
			TagModel:            config.TagModel,
			TransactionTagModel: config.TransactionTagModel,
			TransactionModel:    config.TransactionModel,
			// initialize other models...
		}
	})
}

func GetModelsService() *ModelsServiceContainer {
	return ModelsService
}
