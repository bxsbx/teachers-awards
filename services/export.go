package services

import (
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
	"strings"
	"teachers-awards/common/util"
	"teachers-awards/global"
	"teachers-awards/model/req"
)

type ExportService struct {
	appCtx context.Context
}

func NewExportService(appCtx context.Context) *ExportService {
	return &ExportService{appCtx: appCtx}
}

func (s *ExportService) mergeCellWithStyle(f *excelize.File, sheetName, hcell, vcell, value string) {
	f.MergeCell(sheetName, hcell, vcell)
	f.SetCellStr(sheetName, hcell, value)
}

func (s *ExportService) ExportReviewList(params *req.ExportReviewListReq, f *excelize.File) (err error) {
	reviewDeclareService := NewReviewDeclareService(s.appCtx)
	params2 := req.GetReviewListReq{}
	util.ObjToObjByReflect(params, &params2)
	data, err := reviewDeclareService.GetReviewList(&params2)
	if err != nil {
		return
	}
	sheetName := "审核列表"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	styleText, _ := f.NewStyle(`{"font":{"family":"微软雅黑"},"alignment":{"horizontal":"center","vertical":"center"}}`)
	subjectMap := global.SubjectEnumMap.SubjectMap
	if params.ReviewProcess == global.ProcessSchool {
		f.SetColStyle(sheetName, "A:J", styleText)
		f.SetColWidth(sheetName, "A", "J", 20)
		f.SetSheetRow(sheetName, "A1", &[]string{"序号", "年度", "一级指标", "二级指标", "分值", "姓名", "性别", "学科", "学校审核", "申报时间"})

		for i, row := range data.List {
			startCell := "A" + strconv.Itoa(i+2)
			f.SetSheetRow(sheetName, startCell,
				&[]interface{}{i + 1, row.Year, row.OneIndicatorName, row.TwoIndicatorName, row.Score, row.UserName, global.SexNameMap[row.UserSex],
					subjectMap[row.SubjectCode], global.ReviewStatusNameMap[row.ReviewStatus], row.CreateTime.Format(util.YMD_HMS)})
		}
	} else if params.ReviewProcess == global.ProcessExpert {
		f.SetColStyle(sheetName, "A:J", styleText)
		f.SetColWidth(sheetName, "A", "J", 20)
		f.SetSheetRow(sheetName, "A1", &[]string{"序号", "一级指标", "二级指标", "分值", "姓名", "性别", "学科", "教师类型", "专家审核", "申报时间"})
		for i, row := range data.List {
			startCell := "A" + strconv.Itoa(i+2)
			f.SetSheetRow(sheetName, startCell,
				&[]interface{}{i + 1, row.OneIndicatorName, row.TwoIndicatorName, row.Score, row.UserName, global.SexNameMap[row.UserSex],
					subjectMap[row.SubjectCode], global.DeclareNameMap[row.DeclareType], global.ReviewStatusNameMap[row.ReviewStatus], row.CreateTime.Format(util.YMD_HMS)})
		}
	} else {
		f.SetColStyle(sheetName, "A:K", styleText)
		f.SetColWidth(sheetName, "A", "K", 20)
		f.SetSheetRow(sheetName, "A1", &[]string{"序号", "年度", "一级指标", "二级指标", "分值", "姓名", "性别", "学科", "教师类型", "教育局审核", "申报时间"})
		for i, row := range data.List {
			startCell := "A" + strconv.Itoa(i+2)
			f.SetSheetRow(sheetName, startCell,
				&[]interface{}{i + 1, row.Year, row.OneIndicatorName, row.TwoIndicatorName, row.Score, row.UserName, global.SexNameMap[row.UserSex],
					subjectMap[row.SubjectCode], global.DeclareNameMap[row.DeclareType], global.ReviewStatusNameMap[row.ReviewStatus], row.CreateTime.Format(util.YMD_HMS)})
		}
	}
	return
}

func (s *ExportService) ExportHistoryActivityList(params *req.ExportHistoryActivityListReq, f *excelize.File) (err error) {
	reviewDeclareService := NewReviewDeclareService(s.appCtx)
	params2 := req.GetHistoryActivityListReq{}
	util.ObjToObjByReflect(params, &params2)
	data, err := reviewDeclareService.GetHistoryActivityList(&params2)
	if err != nil {
		return
	}
	sheetName := "历史活动结果"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	styleText, _ := f.NewStyle(`{"font":{"family":"微软雅黑"},"alignment":{"horizontal":"center","vertical":"center"}}`)

	f.SetColStyle(sheetName, "A:J", styleText)
	f.SetColWidth(sheetName, "A", "J", 20)
	f.SetSheetRow(sheetName, "A1", &[]string{"排名", "年度", "活动名称", "姓名", "身份证号", "性别", "所属学校、学科", "申报类型", "分数", "名次"})
	subjectMap := global.SubjectEnumMap.SubjectMap
	for i, row := range data.List {
		startCell := "A" + strconv.Itoa(i+2)
		f.SetSheetRow(sheetName, startCell,
			&[]interface{}{row.Rank, row.Year, row.ActivityName, row.UserName, row.IdentityCard, global.SexNameMap[row.UserSex], row.SchoolName + "——" + subjectMap[row.SubjectCode],
				global.DeclareNameMap[row.DeclareType], row.FinalScore, fmt.Sprintf("%v\n奖金：%v元", global.RankPrizeNameMap[*row.RankPrize], row.Prize)})
	}
	return
}

func (s *ExportService) ExportUserDeclareRecordListByYear(params *req.ExportUserDeclareRecordListByYearReq, f *excelize.File) (err error) {
	userActivityService := NewUserActivityService(s.appCtx)
	params2 := req.GetUserDeclareRecordListByYearReq{}
	util.ObjToObjByReflect(params, &params2)
	data, err := userActivityService.GetUserDeclareRecordListByYear(&params2)
	if err != nil {
		return
	}
	sheetName := "用户的申报记录"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	styleText, _ := f.NewStyle(`{"font":{"family":"微软雅黑"},"alignment":{"horizontal":"center","vertical":"center"}}`)
	//教育局
	if params.Role == global.RoleEdb {
		f.SetColStyle(sheetName, "A:G", styleText)
		f.SetColWidth(sheetName, "A", "G", 20)
		f.SetSheetRow(sheetName, "A1", &[]string{"一级指标", "二级指标", "教师自评", "学校评分", "专家评分", "教育局评分", "申报时间"})
		for i, row := range data.List {
			startCell := "A" + strconv.Itoa(i+2)
			var zjPass []string
			process := 0
			for j := 0; j < len(row.JudgesVerifyList); j++ {
				isPass := row.JudgesVerifyList[j].IsPass
				score := 0
				if isPass == global.PassYes {
					score = row.Score
				}
				zjPass = append(zjPass, fmt.Sprintf("%v(%v)分", global.PassNameMap[isPass], score))
				process = row.JudgesVerifyList[j].JudgesType
			}
			status := global.StatusNameMap[row.Status]
			switch process {
			case 0:
				zjPass = append(zjPass, status, status, status)
			case 1:
				zjPass = append(zjPass, status, status)
			case 2:
				zjPass = append(zjPass, status)
			}
			f.SetSheetRow(sheetName, startCell,
				&[]interface{}{row.OneIndicatorName, row.TwoIndicatorName, row.Score, zjPass[0],
					strings.Join(zjPass[1:len(zjPass)-1], ";"), zjPass[len(zjPass)-1], row.CreateTime.Format(util.YMD_HMS)})
		}
	} else {
		f.SetColStyle(sheetName, "A:E", styleText)
		f.SetColWidth(sheetName, "A", "E", 20)
		f.SetSheetRow(sheetName, "A1", &[]string{"一级指标", "二级指标", "分值", "状态", "申报时间"})
		for i, row := range data.List {
			startCell := "A" + strconv.Itoa(i+2)
			f.SetSheetRow(sheetName, startCell,
				&[]interface{}{row.OneIndicatorName, row.TwoIndicatorName, row.Score, global.StatusNameMap[row.Status], row.CreateTime.Format(util.YMD_HMS)})
		}
	}
	return
}

func (s *ExportService) ExportEdbReviewList(params *req.ExportEdbReviewListReq, f *excelize.File) (err error) {
	reviewDeclareService := NewReviewDeclareService(s.appCtx)
	params2 := req.GetEdbReviewListReq{}
	util.ObjToObjByReflect(params, &params2)
	data, err := reviewDeclareService.GetEdbReviewList(&params2)
	if err != nil {
		return
	}
	sheetName := "审核列表"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	styleText, _ := f.NewStyle(`{"font":{"family":"微软雅黑"},"alignment":{"horizontal":"center","vertical":"center"}}`)

	f.SetColStyle(sheetName, "A:H", styleText)
	f.SetColWidth(sheetName, "A", "H", 20)

	f.SetSheetRow(sheetName, "A1", &[]string{"姓名", "性别", "身份证号", "类型", "学校", "学科", "总分", "状态"})
	subjectMap := global.SubjectEnumMap.SubjectMap
	for i, row := range data.List {
		startCell := "A" + strconv.Itoa(i+2)
		status := "待审核"
		if row.Status == 2 {
			status = "已审核"
		}
		f.SetSheetRow(sheetName, startCell,
			&[]interface{}{row.UserName, global.SexNameMap[row.UserSex], row.IdentityCard, global.DeclareNameMap[row.DeclareType], row.SchoolName, subjectMap[row.SubjectCode],
				row.FinalScore, status})
	}
	return
}

func (s *ExportService) ExportAwardsSetList(params *req.ExportAwardsSetListReq, f *excelize.File) (err error) {
	reviewDeclareService := NewReviewDeclareService(s.appCtx)
	params2 := req.GetAwardsSetListReq{}
	util.ObjToObjByReflect(params, &params2)
	data, err := reviewDeclareService.GetAwardsSetList(&params2)
	if err != nil {
		return
	}
	sheetName := "获奖列表"
	f.SetSheetName(f.GetSheetName(0), sheetName)
	styleText, _ := f.NewStyle(`{"font":{"family":"微软雅黑"},"alignment":{"horizontal":"center","vertical":"center"}}`)

	f.SetColStyle(sheetName, "A:H", styleText)
	f.SetColWidth(sheetName, "A", "H", 20)

	s.mergeCellWithStyle(f, sheetName, "A1", "D1", fmt.Sprintf("%v申报结果-%v类(%v)", data.ActivityName, global.DeclareNameMap[params.DeclareType], data.Total))
	s.mergeCellWithStyle(f, sheetName, "E1", "H1", fmt.Sprintf("获奖人数：%v人，一等奖%v人；二等奖%v人；三等奖%v人", data.AwardNum, data.FirstPrizeNum, data.SecondPrizeNum, data.ThirdPrizeNum))

	f.SetSheetRow(sheetName, "A2", &[]string{"排名", "姓名", "身份证号", "性别", "所属学校、学科", "申报类型", "分数", "名次"})

	subjectMap := global.SubjectEnumMap.SubjectMap
	for i, row := range data.List {
		startCell := "A" + strconv.Itoa(i+3)
		f.SetSheetRow(sheetName, startCell,
			&[]interface{}{row.Rank, row.UserName, row.IdentityCard, global.SexNameMap[row.UserSex], row.SchoolName + "——" + subjectMap[row.SubjectCode],
				global.DeclareNameMap[row.DeclareType], row.FinalScore, fmt.Sprintf("%v\n奖金：%v元", global.RankPrizeNameMap[*row.RankPrize], row.Prize)})
	}
	return
}
