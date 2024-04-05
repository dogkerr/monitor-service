package postgres

import (
	"context"
	"dogker/lintang/monitor-service/pkg/gorm"
)


type ContainerRepository struct {
	DB *gorm.Gorm
}


func NewContainerRepo(db *gorm.Gorm) *ContainerRepository {
	return &ContainerRepository{db}
}


func (r *ContainerRepository) GetById(ctx context.Context) string {
	return "tes"
}