package userHandler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"todoBackend/app/config"
	"todoBackend/app/models"
	"todoBackend/app/service/userService"
	"todoBackend/utils"
	"todoBackend/utils/token"

	"github.com/gin-gonic/gin"
)

// Register 用于处理用户注册请求。
func Register(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	err := userService.CreateUser(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(inputUser, "success"))
}

// Login 用于处理用户登录的请求
func Login(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	loginCheck, err := userService.LoginCheck(&inputUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "登录失败"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(loginCheck, "get loginCheck success!"))
}

// CurrentUser 用来获取当前用户的信息并返回给前端
func CurrentUser(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := userService.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(u, "success"))
}

// UpdateUser 用来更新用户信息
func UpdateUser(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userFromDB, err := userService.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		return
	}
	if err := userService.UpdateUser(&inputUser, &userFromDB); err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(err.Error(), "update failed"))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(userFromDB, "success update"))

}

// UploadAvatar 用于处理上传头像的请求
func UploadAvatar(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	// 获取用户id
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "error"))
		return
	}

	// 获取文件名
	fileName := file.Filename
	// 获取文件后缀
	extName := filepath.Ext(fileName)
	// 通过id获取用户信息
	userFromDB, err := userService.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	// 将文件保存到指定路径
	err = c.SaveUploadedFile(file, "app/static/avatars/"+strconv.Itoa(int(userId))+extName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error(), "保存文件失败"))
		return
	}
	// 更新用户头像
	server := config.Cfg.Server.URL
	uploadUrl := fmt.Sprintf("%s/users/avatars/", server) + strconv.Itoa(int(userId)) + extName
	err = userService.UpdateAvatar(&userFromDB, uploadUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(err.Error(), "更新用户头像失败"))
		return
	}
	// 返回成功信息
	c.JSON(http.StatusOK, utils.SuccessResponse(userFromDB, "avatar uploaded successfully"))
}
func UpdateBio(c *gin.Context) {
	var updateUser models.User
	if err := c.BindJSON(&updateUser); err != nil {
		c.JSON(400, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	//获取这个用户的ID
	userid, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//获取信息
	userFromDB, err := userService.GetUserByID(userid)
	if err != nil {
		c.JSON(400, utils.ErrorResponse(err.Error(), "error"))
		return
	}
	if len(updateUser.Bio) > 250 {
		c.JSON(400, gin.H{"error": "Bio is too long"})
		return
	}
	if err := userService.Updatebio(&updateUser, userFromDB.Bio); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, gin.H{"message": "Bio updated successfully"})
	}
}
func DleBio(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userFromDB, err := userService.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = userService.Dlebio(&userFromDB); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userFromDB})
}
