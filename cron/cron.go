/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : cron.go
*
*                   @Programmer   : 刘泽奇
*
*                   @Start Date   : 2021/8/12 16:56
*
*                   @Last Update  : 2021/8/12 16:56
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

package cron

import (
	"fmt"
	"go_crontab/config"
	"plugin"
	"strconv"
	"strings"
	"time"
)

type Cron struct {
}
//flag = 1 , 默认是就是某个月,某天,某时,某分执行的
//flag = 2 , 隔一段时间执行,只支持某一个,从前向后优先
type CTicker struct {

	Flag int
	PluginName string
	Year int
	Month int/**下次执行时间**/
	Day int
	Hour int
	Minute int/**下次执行时间**/
	YearIncrement int/**执行时间增量**/
	MonthIncrement int
	DayIncrement int
	HourIncrement int
	MinuteIncrement int/**执行时间增量**/
	MonthStart int/**执行时间区间开始时间**/
	DayStart int
	HourStart int
	MinuteStart int/**执行时间区间开始时间**/
	MonthEnd int/**执行时间区间结束时间**/
	DayEnd int
	HourEnd int
	MinuteEnd int/**执行时间区间结束时间**/
	UpdateFlag int

}

type DayStruct struct {
	Days int
	Months int
}

type ExecAndRet interface {
	Exec() ( ErrorCode int , ErrorMsg string , ErrorInfo error )
}

var (

	plugins []string
	cTickers []CTicker
	leapYearMonthToDay = map[int]int{ 0:0 , 1:31 , 2:60 , 3:91 , 4:121 , 5:152 , 6:182 , 7:213 , 8:244 , 9:274 , 10:305 , 11:335 , 12:366 }
	leapYearDayToYear [367]DayStruct
	nonLeapYearMonthToDay = map[int]int{ 0:0 , 1:31 , 2:59 , 3:90 , 4:120 , 5:151 , 6:181 , 7:212 , 8:243 , 9:273 , 10:304 , 11:334 , 12:365 }
	nonLeapYearDayToYear [366]DayStruct
)

func init (){

	for i:=1 ; i<= 366 ; i++ {
		ds := DayStruct{}
		if i <= 0 {
			ds.Months = 0
			ds.Days   = 0
		}else if i <= 31 {
			ds.Months = 1
			ds.Days   = i
		}else if i <= 60 {
			ds.Months = 2
			ds.Days   = i - 31
		}else if i <= 91 {
			ds.Months = 3
			ds.Days   = i - 60
		}else if i <= 121 {
			ds.Months = 4
			ds.Days   = i - 91
		}else if i <= 152 {
			ds.Months = 5
			ds.Days   = i - 121
		}else if i <= 182 {
			ds.Months = 6
			ds.Days   = i - 152
		}else if i <= 213 {
			ds.Months = 7
			ds.Days   = i - 182
		}else if i <= 244 {
			ds.Months = 8
			ds.Days   = i - 213
		}else if i <= 274 {
			ds.Months = 9
			ds.Days   = i - 244
		}else if i <= 305 {
			ds.Months = 10
			ds.Days   = i - 274
		}else if i <= 335 {
			ds.Months = 11
			ds.Days   = i - 305
		}else if i <= 366 {
			ds.Months = 12
			ds.Days   = i - 335
		}
		leapYearDayToYear[i] = ds
	}
	for i:=1 ; i<= 365 ; i++ {
		ds := DayStruct{}
		if i <= 0 {
			ds.Months = 0
			ds.Days   = 0
		}else if i <= 31 {
			ds.Months = 1
			ds.Days   = i
		}else if i <= 59 {
			ds.Months = 2
			ds.Days   = i - 31
		}else if i <= 90 {
			ds.Months = 3
			ds.Days   = i - 59
		}else if i <= 120 {
			ds.Months = 4
			ds.Days   = i - 90
		}else if i <= 151 {
			ds.Months = 5
			ds.Days   = i - 120
		}else if i <= 181 {
			ds.Months = 6
			ds.Days   = i - 151
		}else if i <= 212 {
			ds.Months = 7
			ds.Days   = i - 181
		}else if i <= 243 {
			ds.Months = 8
			ds.Days   = i - 212
		}else if i <= 273 {
			ds.Months = 9
			ds.Days   = i - 243
		}else if i <= 304 {
			ds.Months = 10
			ds.Days   = i - 273
		}else if i <= 334 {
			ds.Months = 11
			ds.Days   = i - 304
		}else if i <= 365 {
			ds.Months = 12
			ds.Days   = i - 334
		}
		nonLeapYearDayToYear[i] = ds
	}
}

func ( C *Cron )InitCTickers(){
	for {
		//cnt := 0
		cTickers = []CTicker{}
		for pName, pTime := range config.PluginsArray {
			tt           := ParseCronExpression(pTime)
			tt.PluginName = pName
			tt.UpdateFlag = 1
			cTickers      = append( cTickers , tt )
			//cnt++
		}
		<-config.UpdatePluginsArray
	}
}

func ( C *Cron )AddCTicker(){
	for{
		pluginContent := <-config.AddPluginContent
		pluginContents := strings.Split( pluginContent , "," )
		tt := ParseCronExpression(pluginContents[1])
		tt.PluginName = pluginContents[0]
		tt.UpdateFlag = 1
		cTickers      = append( cTickers , tt )
		plugins       = append( plugins , pluginContents[0] )
	}
}

func ( C *Cron )Split( ) {
	for {
		fmt.Println("初始化模块",config.Plugins)
		plugins = strings.Split( config.Plugins , "," )
		<-config.UpdatePlugins
	}
}

func ( C *Cron ) Exec( T time.Time ){

	fmt.Println( "开始执行任务集:" , plugins )
	month  := int( T.Month() )
	day    := T.Day()
	hour   := T.Hour()
	minute := T.Minute()
	fmt.Println( T , cTickers )
	for _ , p := range plugins {
		fmt.Println( "开始执行任务:" , p )
		execTask( p , month , day , hour, minute  , T)
	}

}

func Task( p string ){
	path := "./plugins/" + config.FirstUpper(p) + ".so"
	pG, err := plugin.Open(path)
	if err != nil {
		fmt.Println("PATH错误:" + path)
		return
	}

	pGE, err := pG.Lookup(config.FirstUpper(p))
	if err != nil {
		fmt.Println("寻找扩展错误:",err)
		return
	}

	pGEA, ok := pGE.(ExecAndRet)
	if ok {
		errorCode , errorMsg , err := pGEA.Exec()
		if err != nil {
			fmt.Println("执行错误")
			return
		}
		fmt.Printf("errorCode: %d, errorMsg: %s \n", errorCode , errorMsg )
	}
}

func execTask( p string , month int, day int, hour int, minute int , T time.Time ) {

	for CTKey , CTValue := range cTickers {
		if p == CTValue.PluginName &&
		   minute == CTValue.Minute && hour == CTValue.Hour &&
		   day == CTValue.Day && month == CTValue.Month {
			go updateCTicker( CTKey , T )
		    go Task(p)
		}
	}
}

func updateCTicker( CTKey int , T time.Time ){

	if  cTickers[CTKey].UpdateFlag != 1 {
		for {
			time.Sleep( time.Second * 1 )
			if  cTickers[CTKey].UpdateFlag == 1 {
				break
			}
		}
	}
	cTickers[CTKey].UpdateFlag = 2
	timeTemplate := "2006-1-2 15:04:05"

	//还原上次时间
	strTime := strconv.Itoa( cTickers[CTKey].Year ) + "-" + strconv.Itoa( cTickers[CTKey].Month ) + "-" + strconv.Itoa( cTickers[CTKey].Day ) +
		" " + strconv.Itoa( cTickers[CTKey].Hour ) + ":" + strconv.Itoa( cTickers[CTKey].Minute) + ":00"
	preTime , _ := time.ParseInLocation( timeTemplate, strTime , time.Local )
    fmt.Println( preTime , strTime )
	//计算下次时间
	nextTime := preTime.AddDate( cTickers[CTKey].YearIncrement , cTickers[CTKey].MonthIncrement , cTickers[CTKey].DayIncrement ).
		Add( time.Hour *  time.Duration(cTickers[CTKey].HourIncrement) ).Add( time.Minute * time.Duration(cTickers[CTKey].MinuteIncrement) )
	fmt.Println(nextTime)

	cTickers[CTKey].Year = nextTime.Year()
	cTickers[CTKey].Month = int(nextTime.Month())
	cTickers[CTKey].Day = nextTime.Day()
	cTickers[CTKey].Hour = nextTime.Hour()
	cTickers[CTKey].Minute = nextTime.Minute()

	cTickers[CTKey].UpdateFlag = 1

}

func ParseCronExpression( cronString string ) ( ct CTicker ) {
	//1 2 3 4 5
	// */ 1 2 3 4
	// 1-2
	nowTime := time.Now()
	cronSegments := strings.Fields( cronString )
	if len(cronSegments) != 4 {
		fmt.Println( "格式错误" )
		return
	}

	if cronSegments[0] != "*" {
		cronSegments0Slashes := strings.Split(cronSegments[0],"/")
		if len( cronSegments0Slashes ) > 1  { //*/x
			ct.Minute = ( nowTime.Minute() + 2 ) / 60

			if ct.Minute == 0{
				nowTime = nowTime.Add( time.Hour * 1 )
			}

			ct.MinuteIncrement , _ = strconv.Atoi( cronSegments0Slashes[1] )
		}else{
			cronSegments0Line := strings.Split( cronSegments[0] , "-" )
			if len( cronSegments0Line ) > 1 { //x-x
				ct.HourIncrement = 1
				ct.Minute , _ = strconv.Atoi( cronSegments0Line[0] )
				ct.MinuteStart = ct.Minute
				ct.MinuteEnd , _ = strconv.Atoi( cronSegments0Line[1] )
				if nowTime.Minute() >= ct.MinuteEnd{
					nowTime = nowTime.Add( time.Hour * 1 )
				}
			} else { //x
				ct.MinuteIncrement = 0
				ct.HourIncrement = 1
				ct.Minute, _ = strconv.Atoi( cronSegments[0] )
				if nowTime.Minute() >= ct.Minute {
					//查看是否需要重置时间
					nowTime = nowTime.Add( time.Hour * 1 )
				}
			}
		}
	}else{// *
		ct.MinuteIncrement = 1
		//留下足够的时间,初始化
		ct.Minute = ( nowTime.Minute() + 2 ) / 60
		if ct.Minute == 0 {
			nowTime = nowTime.Add( time.Hour * 1 )
		}
	}

	if cronSegments[1] != "*" {
		cronSegments1Slashes := strings.Split( cronSegments[1] , "/" )
		if len( cronSegments1Slashes ) > 1 { //*/x
			ct.Hour = nowTime.Hour()
			ct.HourIncrement , _ = strconv.Atoi( cronSegments1Slashes[1] )
		}else{
			cronSegments1Line := strings.Split( cronSegments[1] , "-" )
			if len( cronSegments1Line ) > 1 { //x-x
				ct.Hour , _ = strconv.Atoi( cronSegments1Line[0] )
				ct.HourStart = ct.Hour
				ct.HourEnd , _ = strconv.Atoi( cronSegments1Line[1] )
			} else { //x
				ct.HourIncrement = 0
				ct.DayIncrement = 1
				ct.Hour , _ = strconv.Atoi( cronSegments[1] )
				if nowTime.Hour() > ct.Hour{
					nowTime = nowTime.Add( time.Hour * 24 ) // nowTime.AddDate(0,0,1)
				}
			}
		}
	}else{  //*
		ct.Hour = nowTime.Hour()
	}

	if cronSegments[2] != "*" {
		cronSegments2Slashes := strings.Split( cronSegments[2] , "/" )
		if len(cronSegments2Slashes) > 1 {
			ct.Day = nowTime.Day()
			ct.DayIncrement , _ = strconv.Atoi( cronSegments2Slashes[1] )
		}else{
			cronSegments2Line := strings.Split( cronSegments[2] , "-" )
			if len(cronSegments2Line) > 1 {
				ct.Day , _ = strconv.Atoi( cronSegments2Line[0] )
				ct.DayStart = ct.Day
				ct.DayEnd , _ = strconv.Atoi( cronSegments2Line[1] )
			} else {
				ct.DayIncrement = 0
				ct.MonthIncrement = 1
				ct.Day , _ = strconv.Atoi( cronSegments[2] )
				if nowTime.Day() > ct.Day{
					nowTime = nowTime.AddDate( 0 ,1 , 0 )
				}
			}
		}
	}else{
		ct.Day = nowTime.Day()
	}

	if cronSegments[3] != "*" {
		cronSegments3Slashes := strings.Split( cronSegments[3] , "/" )
		if len(cronSegments3Slashes) > 1 {
			ct.Month = int( nowTime.Month() )
			ct.MonthIncrement , _ = strconv.Atoi( cronSegments3Slashes[1] )
		} else {
			cronSegments3Line := strings.Split( cronSegments[3] , "-" )
			if len(cronSegments3Line) > 1 {
				ct.Month , _ = strconv.Atoi( cronSegments3Line[0] )
				ct.MonthStart = ct.Month
				ct.MonthEnd , _ = strconv.Atoi( cronSegments3Line[1] )
			} else {
				ct.MonthIncrement = 0
				ct.YearIncrement = 1
				ct.Month , _ = strconv.Atoi( cronSegments[3] )
			}
		}
	}else{
		ct.Month = int( nowTime.Month() )
	}
	ct.Year = nowTime.Year()
	return
}

