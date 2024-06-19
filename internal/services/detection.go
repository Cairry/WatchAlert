package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"watchAlert/internal/models"
	"watchAlert/pkg/ctx"
	utilsHttp "watchAlert/pkg/utils/http"
)

type (
	detectionService struct {
		ctx *ctx.Context
	}

	InterDetectionService interface {
		HTTP(req interface{}) (interface{}, interface{})
		ListSite() (interface{}, interface{})
		CreateSite(req interface{}) (interface{}, interface{})
	}
)

func NewInterDetectionService(ctx *ctx.Context) InterDetectionService {
	return detectionService{
		ctx: ctx,
	}
}

func (d detectionService) ListSite() (interface{}, interface{}) {
	list, err := d.ctx.DB.Detection().List()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (d detectionService) CreateSite(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DetectionSite)
	err := d.ctx.DB.Detection().Create(*r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (d detectionService) HTTP(req interface{}) (interface{}, interface{}) {
	r := req.(*models.DetectionReq)

	for _, site := range r.SiteList {
		get, err := d.ctx.DB.Detection().Get(models.DetectionSiteQuery{ID: site.ID})
		if err != nil {
			return nil, err
		}

		b := strings.NewReader(fmt.Sprintf(`{"protocol": "%s","target": "%s"}`, r.Protocol, r.Target))
		res, err := utilsHttp.Post("http://"+get.Endpoint+"/http", b)

		if err != nil {
			logrus.Errorf(err.Error())
			return nil, err
		}

		body, _ := io.ReadAll(res.Body)
		fmt.Println("===>", string(body))
	}

	return "", nil
}
