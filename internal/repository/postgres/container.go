package postgres

import (
	"context"

	"gorm.io/gorm"
)


type ContainerRepository struct {
	DB *gorm.DB
}


func NewContainerRepo(db *gorm.DB) *ContainerRepository {
	return &ContainerRepository{db}
}


func (r *ContainerRepository) GetById(ctx context.Context) string {
	return "tes"
}