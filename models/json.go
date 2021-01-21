/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: json.go
* @Created:   2021-01-18 10:28:03
* @Modified:  2021-01-21 14:00:13
 */

package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// SendValidateCodeJSON 验证码json
type SendValidateCodeJSON struct {
	Email string `json:"email"`
}

// Validate 发送验证码时的邮箱验证
func (scm SendValidateCodeJSON) Validate() error {
	return validation.ValidateStruct(&scm,
		validation.Field(&scm.Email, validation.Required, validation.Length(5, 40)),
	)
}

// RegisterJSON 注册时序列化用的json
type RegisterJSON struct {
	User
	ValidateCode string `json:"validate_code"`
}

// Validate 注册信息验证
func (rj RegisterJSON) Validate() error {
	return validation.ValidateStruct(&rj,
		validation.Field(&rj.Username, validation.Required, validation.Length(3, 20)),
		validation.Field(&rj.Password, validation.Required, validation.Length(8, 40)),
		validation.Field(&rj.Email, validation.Required, is.Email, validation.Length(5, 40)),
		// TODO: 手机号的正则规则待改进
		validation.Field(&rj.Phone,
			validation.Required,
			validation.Match(regexp.MustCompile("^1[0-9]{10}$")).Error("must be a string with 11 digits"),
		),
		validation.Field(&rj.ValidateCode, validation.Required, validation.Length(6, 6)),
	)
}

// LoginJSON 登录时序列化用的json
type LoginJSON struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

// Validate 注册信息验证
func (lj LoginJSON) Validate() error {
	return validation.ValidateStruct(&lj,
		validation.Field(&lj.UserID, validation.Required, validation.Length(3, 40)),
		validation.Field(&lj.Password, validation.Required, validation.Length(8, 40)),
	)
}
