/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : db.go
*
*                   @Programmer   : 刘泽奇
*
*                   @Start Date   : 2021/8/6 11:13
*
*                   @Last Update  : 2021/8/6 11:13
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

package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"go_crontab/config"
	"reflect"
	"strconv"
	"strings"
)
type DB struct{}
type Ret struct{
	ErrorCode string
	ErrorMessage string
}
func ( D *DB )Connect(){

	dbDSN := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local&timeout=3600s&interpolateParams=true",
		config.MysqlConfig.User,
		config.MysqlConfig.Password,
		config.MysqlConfig.Host,
		config.MysqlConfig.Port,
		config.MysqlConfig.DbName)

	var err error
	config.MysqlDb, err = sql.Open("mysql", dbDSN)
	if err != nil {
		log.Println("dbDSN: " + dbDSN)
		panic("数据源配置不正确: " + err.Error())
	} else {
		log.Println("连接成功 ")
	}

}

func ( D *DB ) Exec( inSql string , retData interface{} ,inType string , inData... interface{} )  ( map[int]interface{} , error ) {

	stmt , err := config.MysqlDb.Prepare( inSql )
	if err != nil{
		return map[int]interface{}{ 0:map[string]string{"errorCode":"2","errorMessage":"Sql预处理失败"}} , err
	}
	rows , err := stmt.Query( inData... )

	switch inType {
	case "update":
		if err == nil{
			return map[int]interface{}{ 0:map[string]string{"errorCode":"1","errorMessage":"更新成功"}} , err
		}else{
			return map[int]interface{}{ 0:map[string]string{"errorCode":"2","errorMessage":"更新失败"}} , err
		}
	case "delete":
		if err == nil{
			return map[int]interface{}{ 0:map[string]string{"errorCode":"1","errorMessage":"删除成功"}} , err
		}else{
			return map[int]interface{}{ 0:map[string]string{"errorCode":"2","errorMessage":"删除失败"}} , err
		}
	case "insert":
		if err == nil{
			return map[int]interface{}{ 0:map[string]string{"errorCode":"1","errorMessage":"写入成功"} } , err
		}else{
			return map[int]interface{}{ 0:map[string]string{"errorCode":"2","errorMessage":"写入失败"} } , err
		}
	}

	if err != nil{
		return map[int]interface{}{ 0:map[string]string{"errorCode":"2","errorMessage":"查询失败"} } , err
	}

	cols, _ := rows.Columns()
	values := make( [][]byte , len(cols) )
	scans := make( []interface{} , len(cols) )
	for i := range values {
		scans[i] = &values[i]
	}
	retDataType := reflect.TypeOf( retData )
	retDataReflect := reflect.New( retDataType ).Elem()
	ret := make( map[int]interface{} )

	i := 0
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return ret, err
		}

		for k, v := range values {

			if retDataReflect.FieldByName(SplitString(cols[k])).IsValid() {
				r , _ := retDataType.FieldByName(SplitString(cols[k]))
				switch r.Type.String() {
				case "string":
					retDataReflect.FieldByName(SplitString(cols[k])).SetString(string(v))
				case "int64":
					t , int64Err := strconv.ParseInt(string(v), 10, 64)
					if int64Err != nil{
						fmt.Println("数据库读取失败",int64Err)
					}
					retDataReflect.FieldByName( SplitString(cols[k]) ).SetInt( t )
				}

			}
		}
		ret[i] = retDataReflect.Interface()
		i++
	}

	return ret, err
}

func SplitString( s string ) string {
	segmentArray := strings.Split(s, "_")
	s = ""
	for _ , EachSegment := range segmentArray {
		EachSegmentII := []rune(EachSegment)
		if len(EachSegmentII) > 0 {
			if EachSegmentII[0] >= 'a' && EachSegmentII[0] <= 'z' {
				EachSegmentII[0] -= 32
			}
			s += string(EachSegmentII)
		}
	}
	return s
}

func ( D *DB ) Add( inSql string ,  inData... interface{} ) ( map[int]interface{} , error ){

	return D.Exec( inSql , "retData" , "update" , inData...)
}

func ( D *DB ) Delete( inSql string  ,  inData... interface{} ) ( map[int]interface{} , error ){

	return D.Exec( inSql , "retData" , "update" , inData...)
}

func ( D *DB ) Update( inSql string  ,  inData... interface{} ) ( map[int]interface{} , error ){

	return D.Exec( inSql , "retData" , "update" , inData...)
}

func ( D *DB ) Select( inSql string  ,  retData interface{} ,  inData... interface{} ) ( map[int]interface{} , error ){

	return D.Exec( inSql , retData , "select" , inData...)
}
