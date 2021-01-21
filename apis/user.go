/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: user.go
* @Created:   2021-01-18 10:26:46
* @Modified:  2021-01-21 15:02:14
 */

package apis

import (
	"blog_api/models"
	"blog_api/utils"
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Verify 验证已注册并登录用户的某些操作
func Verify(c *fiber.Ctx) error {
	return nil
}

// SendValidateCodeViaMail 通过邮件发送验证码
func SendValidateCodeViaMail(c *fiber.Ctx) error {
	var emailJSON models.SendValidateCodeJSON
	if err := c.BodyParser(&emailJSON); err != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, err)
	}

	if err := emailJSON.Validate(); err != nil {
		return utils.ErrorJSON(c, fiber.StatusBadRequest, err)
	}

	// utils.ParseVisitInfo(string(c.Request().Header.UserAgent()), "120.6.51.129")

	go utils.SendValidateCodeViaMail(emailJSON.Email, string(c.Request().Header.UserAgent()), "120.6.51.129")

	return c.Status(fiber.StatusOK).JSON(models.BaseResponse{
		Code:    fiber.StatusOK,
		Message: "Validate code is sent to your email",
	})
}

// Register 用户注册，
// 验证码无论正确与否，都将新注册的账号保存到数据库里。
// 验证码正确的，在下次登录时可正常进行登录后的操作。
// 验证码错误的，标记为未验证用户，登录后需重新进行验证。
func Register(c *fiber.Ctx) error {
	var registerJSON models.RegisterJSON
	if err := c.BodyParser(&registerJSON); err != nil {
		return utils.ErrorJSON(c, fiber.StatusForbidden, err)
	}

	if err := registerJSON.Validate(); err != nil {
		return utils.ErrorJSON(c, fiber.StatusForbidden, err)
	}

	user := registerJSON.User
	user.IsAdmin = utils.IsAdmin(user.Username)
	user.RegisterIP = c.IP()

	user.IsVerified = utils.ValidateCodeIsValid(registerJSON.Email, registerJSON.ValidateCode)

	db := utils.GetDB()

	password, err := utils.GeneratePassword(user.Password)
	if err != nil {
		return utils.ErrorJSON(c, fiber.StatusForbidden, err)
	}
	user.Password = password

	result := db.Create(&user)
	if result.Error != nil {
		db.Rollback()
		errStr := result.Error.Error()
		if strings.Contains(errStr, "Error 1062: Duplicate entry") {
			err := utils.DatabaseExistError(errStr)
			return utils.ErrorJSON(c, fiber.StatusForbidden, err)
		}

	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      fiber.StatusCreated,
		"msg":         "Account registration is successful",
		"username":    user.Username,
		"user_id":     user.ID,
		"is_verified": user.IsVerified,
	})
}

// Login 登录，
func Login(c *fiber.Ctx) error {
	var loginJSON models.LoginJSON
	if err := c.BodyParser(&loginJSON); err != nil {
		return utils.ErrorJSON(c, fiber.StatusForbidden, err)
	}

	if err := loginJSON.Validate(); err != nil {
		return utils.ErrorJSON(c, fiber.StatusForbidden, err)
	}

	var user models.User
	db := utils.GetDB()
	if strings.Contains(loginJSON.UserID, "@") {
		db.Where("email = ?", loginJSON.UserID).First(&user)
	} else {
		// 先用用户名查询
		db.Where("username = ?", loginJSON.UserID).First(&user)
		if user.ID == 0 {
			// 用户名查询不到时，用手机号查询
			// TODO: 用手机号查询前，先用正则判断是否为正确的手机号格式，如果是才开始查询
			db.Where("phone = ?", loginJSON.UserID).First(&user)
		}
	}

	if user.ID == 0 {
		return utils.ErrorJSON(c, fiber.StatusUnauthorized, errors.New("User not exists"))
	}

	// 如果用户已注销，返回错误信息
	if user.IsWriteOff {
		return utils.ErrorJSON(c, fiber.StatusForbidden, errors.New("The currently logged-in account has been wrote off"))
	}

	// 未验证用户重新发送验证邮件，否则登录失败
	if !user.IsVerified {
		go utils.SendValidateCodeViaMail(user.Email, string(c.Request().Header.UserAgent()), "120.6.51.129")
		return utils.ErrorJSON(c,
			fiber.StatusForbidden,
			errors.New("Your account has not been verified. We just sent an email containing a verification code to your registered email address"),
		)
	}

	if !utils.CheckPassword(user.Password, loginJSON.Password) {
		return utils.ErrorJSON(c, fiber.StatusForbidden, errors.New("Wrong password"))
	}

	loginIP := c.IP()

	// 使用协程将登录时间和登录ip保存到数据库中
	go func() {
		db.Model(&user).Updates(map[string]interface{}{
			"last_login_time": time.Now(),
			"last_login_ip":   loginIP,
		})
	}()

	token := utils.GenerateToken()
	if err := utils.SetToken(token, user.ID, 0); err != nil {
		return utils.ErrorJSON(
			c,
			fiber.StatusInternalServerError,
			errors.New("The system is busy, please try again later"),
		)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   60 * 60 * 24,
		Secure:   false,
		HTTPOnly: true,
		Path:     "/",
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"username": user.Username,
		"user_id":  user.ID,
		"msg":      "Successful login",
		"token":    token,
	})
}
