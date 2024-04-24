package controllers

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"teachers-awards/common/errorz"
	"teachers-awards/model/resp"
)

const (
	OK       = 10000
	OKMSG    = "成功"
	ERRSTACK = "errStack"
)

func OutputError(c *gin.Context, err error) {
	errStack := errorz.GetErrorCallerList(err)
	c.Set(ERRSTACK, errStack)
	code, msg := errorz.GlobalError(err)
	c.JSON(http.StatusOK, resp.Response{Code: code, Msg: msg, ErrStack: errStack})
}

func OutputSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, resp.Response{Code: OK, Msg: OKMSG, Data: data})
}

func Output(c *gin.Context, code int, msg string, data interface{}, err error) {
	errStack := errorz.GetErrorCallerList(err)
	response := resp.Response{Code: code, Msg: msg}
	if data != nil {
		response.Data = data
	}
	if errStack != nil {
		c.Set(ERRSTACK, errStack)
		response.ErrStack = errStack
	}
	c.JSON(http.StatusOK, response)
}

func ExportExcelFile(c *gin.Context, f *excelize.File, fileName string, err error) {
	if err != nil {
		errStack := errorz.GetErrorCallerList(err)
		code, msg := errorz.GlobalError(err)
		c.JSON(http.StatusInternalServerError, resp.Response{Code: code, Msg: msg, ErrStack: errStack})
	} else {
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+fileName+".xlsx")
		c.Writer.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		f.Write(c.Writer)
	}
}
