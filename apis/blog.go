/*
* @Author: thepoy
* @Email:  thepoy@aliyun.com
* @File Name: blog.go
* @Created:   2021-01-17 12:32:53
* @Modified:  2021-01-19 20:42:40
 */

package apis

import (
	"blog_api/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetBlog 根据id获取blog信息
func GetBlog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 0)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(&models.ErrorResponse{
			Error: fmt.Sprintf("`%s` is not an integer", c.Params("id")),
		})
	}
	blog := &models.Blog{}
	blog.ID = uint(id)
	blog.Title = "It is a blog."
	return c.JSON(blog)
}
