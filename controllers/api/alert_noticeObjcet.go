package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
)

type AlertNoticeObjectController struct{}

func (ano *AlertNoticeObjectController) List(ctx *gin.Context) {

	object := alertNoticeService.SearchNoticeObject()
	response.Success(ctx, object, "success")

}

func (ano *AlertNoticeObjectController) Create(ctx *gin.Context) {

	var alertNotice models.AlertNotice
	_ = ctx.ShouldBindJSON(&alertNotice)

	object, err := alertNoticeService.CreateNoticeObject(alertNotice)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, object, "success")

}

func (ano *AlertNoticeObjectController) Update(ctx *gin.Context) {

	var alertNotice models.AlertNotice
	_ = ctx.ShouldBindJSON(&alertNotice)

	object, err := alertNoticeService.UpdateNoticeObject(alertNotice)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, object, "success")

}

func (ano *AlertNoticeObjectController) Delete(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	err := alertNoticeService.DeleteNoticeObject(uuid)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}
	response.Success(ctx, "", "success")

}

func (ano *AlertNoticeObjectController) Get(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	object := alertNoticeService.GetNoticeObject(uuid)
	response.Success(ctx, object, "success")

}

func (ano *AlertNoticeObjectController) CheckNoticeStatus(ctx *gin.Context) {

	uuid := ctx.Query("uuid")
	status := alertNoticeService.CheckNoticeObjectStatus(uuid)
	response.Success(ctx, status, "success")

}
