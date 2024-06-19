package repo

import (
	"gorm.io/gorm"
	"watchAlert/internal/models"
	"watchAlert/pkg/utils/cmd"
)

type (
	detectionRepo struct {
		entryRepo
	}

	InterDetectionRepo interface {
		Create(r models.DetectionSite) error
		List() ([]models.DetectionSite, error)
		Get(r models.DetectionSiteQuery) (models.DetectionSite, error)
	}
)

func newInterDetectionRepo(db *gorm.DB, g InterGormDBCli) InterDetectionRepo {
	return detectionRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

func (dr detectionRepo) List() ([]models.DetectionSite, error) {
	var list []models.DetectionSite

	db := dr.DB().Model(models.DetectionSite{})
	err := db.Find(&list).Error
	if err != nil {
		return list, err
	}

	return list, nil
}

func (dr detectionRepo) Get(r models.DetectionSiteQuery) (models.DetectionSite, error) {
	db := dr.DB().Model(models.DetectionSite{})
	if r.ID != "" {
		db.Where("id = ?", r.ID)
	}
	if r.Name != "" {
		db.Where("name = ?", r.Name)
	}

	var get models.DetectionSite
	err := db.First(&get).Error
	if err != nil {
		return get, err
	}

	return get, nil
}

func (dr detectionRepo) Create(r models.DetectionSite) error {
	r.ID = cmd.RandId()
	err := dr.g.Create(models.DetectionSite{}, r)
	if err != nil {
		return err
	}

	return nil
}
