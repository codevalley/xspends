package impl

import (
	"xspends/models/interfaces"
)

type ModelsServiceContainer struct {
	DBService           *DBService
	CategoryModel       interfaces.CategoryService
	SourceModel         interfaces.SourceService
	UserModel           interfaces.UserService
	TagModel            interfaces.TagService
	TransactionTagModel interfaces.TransactionTagService
	// other models...
}

func NewModelsServiceContainer() *ModelsServiceContainer {
	return &ModelsServiceContainer{
		DBService:           GetDBService(),   // Using the existing DBService instance
		CategoryModel:       &CategoryModel{}, // Initialize other models as needed
		SourceModel:         &SourceModel{},
		UserModel:           &UserModel{},
		TagModel:            &TagModel{},
		TransactionTagModel: &TransactionTagModel{},
		// initialize other models...
	}
}

var ModelsService *ModelsServiceContainer

func InitModelsService() error {
	err := InitDB()
	if err == nil {
		ModelsService = NewModelsServiceContainer()
		return nil
	}
	return err
}

func GetModelsService() *ModelsServiceContainer {
	return ModelsService
}
