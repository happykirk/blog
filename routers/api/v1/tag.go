package v1

import (
	"github.com/unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/happykirk/blog/models"
	"github.com/happykirk/blog/pkg/e"
	"github.com/happykirk/blog/pkg/setting"
	"github.com/happykirk/blog/pkg/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"


)
type tagInfo struct {
	Name      string `json:"name"`
	State   int    `json:"state"`
	Created_by   string  `json:"created_by"`
}

//获取多个文章标签
func GetTags(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS
	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})

}

//新增文章标签
func AddTag(c *gin.Context) {

	var name =""
	var state  = 0
	var createBy = ""

	//post json格式，注意json的字符串和数字
	if c.Request.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(c.Request.Body)
		code := e.INVALID_PARAMS
		if err != nil {
			log.Println(err, "read http body failed！error msg:"+err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]string),
			})
			return
		}
		log.Println("INFO|OrderAdd body :" + string(body))

		var params tagInfo
		err = json.Unmarshal(body, &params)
		if err != nil {
			log.Println(err, "read http body json failed! err :"+err.Error())
			code := e.ERROR_READ_HTTP_BODY_JSON_FAILED
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": make(map[string]string),
			})
			return
		}
		name = params.Name
		state = com.StrTo(params.State).MustInt()
		createBy = params.Created_by


	} else {
		//Query的方式是：http://127.0.0.1:8000/api/v1/tags?name=1&state=1&created_by=test
		name = c.Query("name")
		state = com.StrTo(c.DefaultQuery("state", "0")).MustInt()
		createBy = c.Query("created_by")
	}




	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或者1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}

	} else {
		fmt.Println(valid.Errors)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

//修改文章标签
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy

			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}

//删除文章标签
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})

}


