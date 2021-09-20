/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : test.go
*
*                   @Programmer   : 刘泽奇
*
*                   @Start Date   : 2021/9/18 16:05
*
*                   @Last Update  : 2021/9/18 16:05
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

package main

import "fmt"
type test int
var Test test
func( *test) Exec()( ErrorCode int , ErrorMsg string , ErrorInfo error ){

	fmt.Println("test")

	return 0,"success",nil

}
