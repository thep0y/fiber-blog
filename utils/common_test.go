/*
* @Author: thepoy
* @Email: thepoy@163.com
* @File Name: common_test.go (c) 2021
* @Created:  2021-01-18 08:43:33
* @Modified: 2021-01-18 16:51:25
 */

package utils

import "testing"

func TestLoadConfig(t *testing.T) {
	t.Logf("dsn: %s", dbConfig)
}
