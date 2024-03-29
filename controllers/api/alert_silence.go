package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/models"
	jwtUtils "watchAlert/utils/jwt"
)

type AlertSilenceController struct {
}

func (asc *AlertSilenceController) Create(ctx *gin.Context) {

	var silence models.AlertSilences
	_ = ctx.ShouldBindJSON(&silence)

	user := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	silence.CreateBy = user

	err := alertSilenceService.CreateAlertSilence(silence)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (asc *AlertSilenceController) Update(ctx *gin.Context) {

	var silence models.AlertSilences
	_ = ctx.ShouldBindJSON(&silence)

	user := jwtUtils.GetUser(ctx.Request.Header.Get("Authorization"))
	silence.UpdateBy = user

	data, err := alertSilenceService.UpdateAlertSilence(silence)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}

func (asc *AlertSilenceController) Delete(ctx *gin.Context) {

	id := ctx.Query("id")
	err := alertSilenceService.DeleteAlertSilence(id)
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, "", "success")

}

func (asc *AlertSilenceController) List(ctx *gin.Context) {

	data, err := alertSilenceService.ListAlertSilence()
	if err != nil {
		response.Fail(ctx, err.Error(), "failed")
		return
	}

	response.Success(ctx, data, "success")

}
