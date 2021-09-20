/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : admin_user.go
*
*                   @Programmer   : 刘泽奇
*
*                   @Start Date   : 2021/8/2 18:19
*
*                   @Last Update  : 2021/8/2 18:19
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

package web

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserStruct struct {
	Id         int64
	Name       string
	Passwd     string
	Uid        string
	Status     string
	CreateTime string
	UpdateTime string
}

func (W *Web) AddUser( w http.ResponseWriter, r *http.Request )  {

	if W.checkLogin(r) != nil {
		_,_ = fmt.Fprintf(w, "未登录")
		return
	}

	name := W.post( r ,"user_name" )
	passWd := W.post( r , "passwd")
	Md5 := md5.New()
	_ , _ = io.WriteString(Md5, passWd)
	rst := fmt.Sprintf("%x", Md5.Sum(nil))

	inSql := "insert into crontab_user ( name , uid , passwd ) values ( ? , ? , ?)"
	_ , err := DB.Add(inSql , name , time.Now().Unix() , rst )
	if err != nil{
		_ , _ = fmt.Fprintf(w, "添加失败")
		return
	}
	_ , _ = fmt.Fprintf(w, "添加成功")
	return
}

func (W *Web) DeleteUser( w http.ResponseWriter, r *http.Request )  {

	if W.checkLogin(r) != nil {
		_,_ = fmt.Fprintf(w, "未登录")
		return
	}

	inSql := "delete from crontab_user where id = ? "
	_ , err := DB.Delete(inSql ,  W.post( r ,"id" ) )
	if err != nil{
		_ , _ = fmt.Fprintf( w, "删除失败" )
		return
	}
	_,_ = fmt.Fprintf( w, "删除成功" )
	return

}

func (W *Web) ListUsers( w http.ResponseWriter, r *http.Request ){

	if W.checkLogin(r) != nil {
		_ , _ = fmt.Fprintf(w, "未登录")
		return
	}

	inSql := "select * from crontab_user"
	ret , err := DB.Select( inSql , UserStruct{} )
	if err != nil {
		_,_ = fmt.Fprintf( w, "获取列表失败" )
		return
	}

	retData , err := json.Marshal( ret )
	if err != nil {
		_ , _ = fmt.Println( "error:" , err )
	}

	_,_ = fmt.Fprintf( w, string(retData) )
	return
}

func (W *Web) UpdateUser( _ http.ResponseWriter, _ *http.Request ) {

}

func (W *Web) checkLogin( r *http.Request ) error {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return err
	}
	if cookie.Value == "" {
		return errors.New("未登录")
	}
	return err
}

func (W *Web) Login(w http.ResponseWriter, r *http.Request){

	name := W.post( r ,"user_name" )
	passWd := W.post( r , "passwd")
	Md5 := md5.New()
	_ , _ = io.WriteString(Md5, passWd)
	rst := fmt.Sprintf("%x", Md5.Sum(nil))
	inSql := "select * from crontab_user where name = ? and passwd = ?"

	ret , err := DB.Select( inSql , UserStruct{} , name , rst)

	if err != nil {
		_,_ = fmt.Fprintf( w, "登录失败,请检查账号或者密码" )
		return
	}

	cookie := http.Cookie{Name: "sessionId", Value: ret[0].(UserStruct).Uid, Expires: time.Now().AddDate( 1 , 0 , 0 )}
	http.SetCookie( w , &cookie )
	return

}