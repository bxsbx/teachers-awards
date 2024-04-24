package controllers

import (
	"github.com/gin-gonic/gin"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/services"
)

type IndicatorController struct {
}

// @Summary	创建或者更新一级指标
// @Tags		Indicator
// @Produce	json
// @Param		one_indicator_id	formData	int		false	"一级指标id"
// @Param		one_indicator_name	formData	string	true	"一级指标名称"
// @Param		content				formData	string	true	"评分标准说明"
// @Response	200					{object}	resp.Response
// @Router		/v1/indicator/one/save [POST]
func (u *IndicatorController) CreateOrUpdateOneIndicator(c *gin.Context) {
	var params req.CreateOrUpdateOneIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	indicatorService := services.NewIndicatorService(global.GetContext(c))
	err := indicatorService.CreateOrUpdateOneIndicator(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	获取一级指标列表
// @Tags		Indicator
// @Produce	json
// @Param		input_name		query		string	false	"输入的一级指标名称"
// @Param		input_content	query		string	false	"输入的评分标准说明"
// @Param		with_two		query		bool	false	"是否获取二级指标"
// @Param		page			query		int		false	"页数"
// @Param		limit			query		int		false	"每页大小"
// @Response	200				{object}	resp.Response{data=resp.GetOneIndicatorListResp}
// @Router		/v1/indicator/one/list [GET]
func (u *IndicatorController) GetOneIndicatorList(c *gin.Context) {
	var params req.GetOneIndicatorListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	indicatorService := services.NewIndicatorService(global.GetContext(c))
	data, err := indicatorService.GetOneIndicatorList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	删除一级指标
// @Tags		Indicator
// @Produce	json
// @Param		one_indicator_ids	query		string	true	"一级指标ids"
// @Response	200					{object}	resp.Response
// @Router		/v1/indicator/one/delete/ids [DELETE]
func (u *IndicatorController) DeleteOneIndicatorByIds(c *gin.Context) {
	var params req.DeleteOneIndicatorByIdsReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	intArray, err := util.StrToIntArray(params.OneIndicatorIds, ",")
	if err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	indicatorService := services.NewIndicatorService(global.GetContext(c))
	err = indicatorService.DeleteOneIndicatorByIds(intArray)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	创建或者更新二级指标
// @Tags		Indicator
// @Produce	json
// @Param		two_indicator_id	formData	int		false	"二级指标id"
// @Param		two_indicator_name	formData	string	true	"二级指标名称"
// @Param		score				formData	int		true	"分值"
// @Param		one_indicator_id	formData	int		true	"所属一级指标"
// @Response	200					{object}	resp.Response{data=resp.CreateOrUpdateTwoIndicatorResp}
// @Router		/v1/indicator/two/save [POST]
func (u *IndicatorController) CreateOrUpdateTwoIndicator(c *gin.Context) {
	var params req.CreateOrUpdateTwoIndicatorReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	indicatorService := services.NewIndicatorService(global.GetContext(c))
	data, err := indicatorService.CreateOrUpdateTwoIndicator(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}

// @Summary	删除二级指标
// @Tags		Indicator
// @Produce	json
// @Param		two_indicator_ids	query		string	true	"二级指标ids"
// @Response	200					{object}	resp.Response
// @Router		/v1/indicator/two/delete/ids [DELETE]
func (u *IndicatorController) DeleteTwoIndicatorByIds(c *gin.Context) {
	var params req.DeleteTwoIndicatorByIdsReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	intArray, err := util.StrToIntArray(params.TwoIndicatorIds, ",")
	if err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	indicatorService := services.NewIndicatorService(global.GetContext(c))
	err = indicatorService.DeleteTwoIndicatorByIds(intArray)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, nil)
}

// @Summary	获取二级指标列表
// @Tags		Indicator
// @Produce	json
// @Param		one_indicator_id	query		int		false	"所属一级指标"
// @Param		input_name			query		string	false	"输入的二级指标名称"
// @Param		input_score			query		int		false	"输入的分值"
// @Param		page				query		int		false	"页数"
// @Param		limit				query		int		false	"每页大小"
// @Response	200					{object}	resp.Response{data=resp.GetTwoIndicatorListResp}
// @Router		/v1/indicator/two/list [GET]
func (u *IndicatorController) GetTwoIndicatorList(c *gin.Context) {
	var params req.GetTwoIndicatorListReq
	if err := c.ShouldBind(&params); err != nil {
		OutputError(c, errorz.CodeError(errorz.RESP_PARAM_ERR, err))
		return
	}
	indicatorService := services.NewIndicatorService(global.GetContext(c))
	data, err := indicatorService.GetTwoIndicatorList(&params)
	if err != nil {
		OutputError(c, err)
		return
	}
	OutputSuccess(c, data)
}
