package api

import (
	"github.com/gin-gonic/gin"
	"watchAlert/controllers/response"
	"watchAlert/globals"
	"watchAlert/models"
)

type UserPermissionsController struct {
}

func (urc *UserPermissionsController) List(ctx *gin.Context) {

	var data []models.UserPermissions
	globals.DBCli.Model(&models.UserPermissions{}).Find(&data)
	response.Success(ctx, data, "success")

}
