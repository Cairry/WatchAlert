package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/internal/models"
	"watchAlert/internal/services"
)

type DetectionCollector struct{}

func (dc DetectionCollector) API(gin *gin.RouterGroup) {
	detection := gin.Group("detection")
	{
		detection.POST("http", dc.HTTP)
		detection.GET("listSite", dc.ListSite)
		detection.POST("createSite", dc.CreateSite)
	}
}

func (dc DetectionCollector) ListSite(ctx *gin.Context) {
	r := new(models.DetectionSiteQuery)
	BindQuery(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DetectionService.ListSite()
	})
}

func (dc DetectionCollector) CreateSite(ctx *gin.Context) {
	r := new(models.DetectionSite)
	BindJson(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DetectionService.CreateSite(r)
	})
}

func (dc DetectionCollector) HTTP(ctx *gin.Context) {
	r := new(models.DetectionReq)
	BindJson(ctx, r)
	Service(ctx, func() (interface{}, interface{}) {
		return services.DetectionService.HTTP(r)
	})
}
