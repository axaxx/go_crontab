/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : main.go
*
*                   @Programmer   : 刘泽奇
*
*                   @Start Date   : 2021/8/2 18:18
*
*                   @Last Update  : 2021/8/2 18:18
*
*-------------------------------------------------------------------*
*
* FUNCTIONS:
*
* WARNING:
* 扩展编译方法,
*
* HISTORY:
*
* DESCRIPTION:
*
*********************************************************************/

package main

import (
	"fmt"
	"net/http"
	"go_crontab/cron"
	"go_crontab/web"
	"time"
)
var (
	C *cron.Cron
	W *web.Web
)


func main() {
	C = new(cron.Cron)
	W = new(web.Web)
	go adminWeb()
	//对齐时间
	fmt.Println("开始对齐时间")
	timeAlignment := make(chan int)
	go func() {
		for {
			NT := <-time.After(time.Second * 1)
			if NT.Second() == 0 {
				timeAlignment <- 1
				return
			}
		}
	}()
	<-timeAlignment
	fmt.Println("时间已对齐")
	CronTicker := time.NewTicker( time.Minute )

	for {
		fmt.Println("开始执行任务集")
		T := <-CronTicker.C
		go C.Exec(T)
	}

}

func adminWeb() {

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { _ , _ = fmt.Fprintf(w, "I am ok.") })
	//初始化plugins
	W.InitPlugins()

	go C.Split()
	go C.InitCTickers()
	go C.AddCTicker()

	http.HandleFunc("/addPlugins", W.AddPlugins)
	//http.HandleFunc("/deletePlugins", W.DeletePlugins)
	http.HandleFunc("/listPlugins", W.ListPlugins)
	//http.HandleFunc("/updatePlugins", W.UpdatePlugins)
	http.HandleFunc("/build", W.Build)

	_ = http.ListenAndServe("0.0.0.0:8081", nil)

}
