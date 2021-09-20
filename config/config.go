/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : config.go
*
*                   @Programmer   : 刘泽奇
*
*                   @Start Date   : 2021/8/6 11:07
*
*                   @Last Update  : 2021/8/6 11:07
*
*-------------------------------------------------------------------*
*
* FUNCTIONS:
*
* WARNING:
*
* HISTORY:
*
* DESCRIPTION:
*
*********************************************************************/

package config

import (
	"database/sql"
	"fmt"
)

type MysqlConfStruct struct {
	Host string
	User string
	Password string
	Port string
	DbName string
}

var (
	Plugins string
	PluginsArray map[string]string
	MysqlConfig *MysqlConfStruct
    MysqlDb *sql.DB
	UpdatePlugins chan int
	UpdatePluginsArray chan int
	AddPluginContent chan string
)

func init(){

	MysqlConfig = new(MysqlConfStruct)
	MysqlConfig.Host = "localhost"
	MysqlConfig.User = "admin"
	MysqlConfig.Password = ""
	MysqlConfig.Port = "3306"
	MysqlConfig.DbName = "test"
	UpdatePlugins = make( chan int )
	UpdatePluginsArray = make( chan int )
	PluginsArray = make( map[string]string )
	AddPluginContent = make (chan string)
}

func FirstUpper(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
