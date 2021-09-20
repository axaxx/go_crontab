/********************************************************************
***                 文件概述 --- 文件概述
*********************************************************************
*
*                   @Project Name : go_crontab
*
*                   @File Name    : admin_web.go
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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"go_crontab/config"
	"go_crontab/db"
	"strings"
)

type Web struct{}
type PluginsStruct struct {
	Id         int64
	Name       string
	Plugin     string
	PluginPath string
	Status     string
	CreateTime string
	UpdateTime string
	ExecTime   string
}

var DB *db.DB

func init(){
	DB = new(db.DB)
	DB.Connect()
}

func (W *Web) post( r *http.Request , key string)string{
	return r.PostFormValue(key)
}

func (W *Web) AddPlugins( w http.ResponseWriter, r *http.Request )  {

	/*if W.checkLogin(r) != nil {
		_,_ = fmt.Fprintf(w, "未登录")
		return
	}*/

	plugins := W.post( r ,"plugin" )
	name := W.post( r ,"plugin_name" )
	execTime := W.post( r , "exec_time")
	pluginPath := "./plugin/" + plugins

	inSql := "insert into crontab_plugins_config ( name , plugin , plugin_path , exec_time ) values ( ? , ? , ? , ? )"
	_ , err := DB.Add( inSql , name , plugins , pluginPath , execTime )
	if err != nil{
		_ , _ = fmt.Fprintf( w , "添加失败" )
		return
	}

	config.AddPluginContent<- plugins + "," + execTime
	_ , _ = fmt.Fprintf( w , "添加成功" )
	return
}

func (W *Web) DeletePlugins( w http.ResponseWriter, r *http.Request )  {

	/*if W.checkLogin(r) != nil {
		_,_ = fmt.Fprintf(w, "未登录")
		return
	}*/

	inSql := "delete from crontab_plugins_config where id = ? "
	_ , err := DB.Delete(inSql ,  W.post( r ,"id" ) )
	if err != nil{
		_,_ = fmt.Fprintf( w, "删除失败" )
		return
	}
	_,_ = fmt.Fprintf( w, "删除成功" )
	return

}

func (W *Web) InitPlugins(){
	inSql := "select * from crontab_plugins_config"
	ret , err := DB.Select( inSql , PluginsStruct{} )
	if err != nil {
		//todo
		config.Plugins = ""
	}
	p := ""
	for _ , pluginsValue := range ret{
		if pluginsValue.(PluginsStruct).Plugin != "" {
			p = p + pluginsValue.(PluginsStruct).Plugin + ","
			config.PluginsArray[ pluginsValue.(PluginsStruct).Plugin ] = pluginsValue.(PluginsStruct).ExecTime
		}
	}
	config.Plugins = strings.Trim( p ,"," )

}

func (W *Web) ListPlugins( w http.ResponseWriter, r *http.Request ){

	/*if W.checkLogin(r) != nil {
		_,_ = fmt.Fprintf(w, "未登录")
		return
	}*/

	inSql := "select * from crontab_plugins_config"
	ret , err := DB.Select( inSql , PluginsStruct{} )
	if err != nil {
		_,_ = fmt.Fprintf( w, "获取列表失败" )
		return
	}

	retData , err := json.Marshal( ret )
	if err != nil {
		_,_ = fmt.Println( "error:" , err )
	}

	_,_ = fmt.Fprintf( w, string(retData) )

	return
}

func (W *Web) UpdatePlugins( _ http.ResponseWriter, _ *http.Request ) {

}

func (W *Web) Build(w http.ResponseWriter, r *http.Request){

	/*if W.checkLogin(r) != nil {
		_ , _ = fmt.Fprintf( w , "未登录" )
		return
	}*/

	p := W.post( r ,"plugin" )
	var out bytes.Buffer

	c := exec.Command( "/bin/bash", "-c", "pwd" )
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		_,_ = fmt.Fprintf( w, "获取路径失败" )
		return
	}
	root := strings.Replace( out.String() , "\n" , "" , -1 )

	pluginSrcPath := root + "/src/" + p + "/" + p + ".go"
	ce := "cd " + root + "/plugins && go build -buildmode=plugin -o " + config.FirstUpper(p) + ".so " + pluginSrcPath
	c = exec.Command("/bin/bash", "-c", ce)
	err = c.Run()
	if err != nil {
		_ , _ = fmt.Fprintf( w, "编译失败:"+ce )
		return
	}
	_ , _ = fmt.Fprintf( w, "编译成功:"+ce )
	return
}

