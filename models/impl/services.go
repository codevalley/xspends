package impl

import (
	"sync"
	"xspends/models/interfaces"
)

const (
	RoleOwner = "owner"
	RoleWrite = "write"
	RoleView  = "view"
)

var roleHierarchy = map[string]int{
	RoleOwner: 3,
	RoleWrite: 2,
	RoleView:  1,
}

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
	ScopeModel          interfaces.ScopeService
	GroupModel          interfaces.GroupService
	UserScopeModel      interfaces.UserScopeService
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
	ScopeModel          interfaces.ScopeService
	GroupModel          interfaces.GroupService
	UserScopeModel      interfaces.UserScopeService
}

var isTesting bool

// InitModelsService is now a function that takes a ModelsConfig struct
func InitModelsService(config *ModelsConfig) {
	if isTesting {
		initializeModelsService(config)
		return
	}
	once.Do(func() {
		initializeModelsService(config)
	})
}

func initializeModelsService(config *ModelsConfig) {
	ModelsService = &ModelsServiceContainer{
		DBService:           config.DBService,
		CategoryModel:       config.CategoryModel,
		SourceModel:         config.SourceModel,
		UserModel:           config.UserModel,
		TagModel:            config.TagModel,
		TransactionTagModel: config.TransactionTagModel,
		TransactionModel:    config.TransactionModel,
		ScopeModel:          config.ScopeModel,
		GroupModel:          config.GroupModel,
		UserScopeModel:      config.UserScopeModel,
	}
}

func GetModelsService() *ModelsServiceContainer {
	return ModelsService
}
