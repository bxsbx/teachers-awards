package services

import (
	"context"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type StatisticsService struct {
	appCtx context.Context
}

func NewStatisticsService(appCtx context.Context) *StatisticsService {
	return &StatisticsService{appCtx: appCtx}
}

func (s *StatisticsService) GetSimpleSumStats(year int, schoolId string) (data *resp.GetSimpleSumStatsResp, err error) {
	data = &resp.GetSimpleSumStatsResp{}
	data.List = make([]resp.SimpleSumStats, 4)
	wg := util.NewWaitGroup(0)

	activityDao := dao.NewActivityDao(s.appCtx)
	lastYear, err := activityDao.GetLastYear(year)
	if err != nil {
		return
	}
	userInfoDao := dao.NewUserInfoDao(s.appCtx)
	//学校
	wg.Go(func() error {
		data.List[0].Type = 1
		if len(schoolId) == 0 {
			count1, err := userInfoDao.GetSchoolNumByYear(year)
			if err != nil {
				return err
			}
			count2, err := userInfoDao.GetSchoolNumByYear(lastYear)
			data.List[0].Num = count1
			data.List[0].YearOnYear = util.OnYear(count1, count2)
			return err
		}
		return nil
	})

	//老师
	wg.Go(func() error {
		data.List[1].Type = 2
		count1, err := userInfoDao.GetTeacherNumByYear(year, schoolId)
		if err != nil {
			return err
		}
		count2, err := userInfoDao.GetTeacherNumByYear(lastYear, schoolId)
		data.List[1].Num = count1
		data.List[1].YearOnYear = util.OnYear(count1, count2)
		return nil
	})
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	userActivity := dao.UserActivity{SchoolId: schoolId}
	//获奖人数
	wg.Go(func() error {
		count1, err := userActivityDao.GetAwardNum(year, userActivity)
		if err != nil {
			return err
		}
		count2, err := userActivityDao.GetAwardNum(lastYear, userActivity)
		data.List[2].Type = 3
		data.List[2].Num = count1
		data.List[2].YearOnYear = util.OnYear(count1, count2)
		return err
	})

	activityIndicatorDao := dao.NewUserActivityIndicatorDao(s.appCtx)
	//申报项目数
	wg.Go(func() error {
		count1, err := activityIndicatorDao.GetDeclareNum(year, schoolId)
		if err != nil {
			return err
		}
		count2, err := activityIndicatorDao.GetDeclareNum(lastYear, schoolId)
		data.List[3].Type = 4
		data.List[3].Num = count1
		data.List[3].YearOnYear = util.OnYear(count1, count2)
		return err
	})
	err = wg.Wait()
	if err != nil {
		return
	}
	if len(schoolId) > 0 {
		data.List = data.List[1:]
	}
	return
}

func (s *StatisticsService) GetAwardRate(year int, schoolId string) (data *resp.GetAwardRateResp, err error) {
	data = &resp.GetAwardRateResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	userActivity := dao.UserActivity{SchoolId: schoolId}
	data.AwardNum, err = userActivityDao.GetAwardNum(year, userActivity)
	if err != nil {
		return
	}
	data.DeclareNum, err = userActivityDao.GetDeclareNum(year, userActivity)
	return
}

func (s *StatisticsService) GetDeclareRate(year int, schoolId string) (data *resp.GetDeclareRateResp, err error) {
	data = &resp.GetDeclareRateResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	data.DeclareNum, err = userActivityDao.GetDeclareNum(year, dao.UserActivity{SchoolId: schoolId})
	if err != nil {
		return
	}
	count, err := dao.NewUserInfoDao(s.appCtx).GetTeacherNumByYear(year, schoolId)
	data.SumNum = count
	return
}

func (s *StatisticsService) GetEveryYearAwardNum(schoolId string) (data *resp.GetEveryYearAwardNumResp, err error) {
	data = &resp.GetEveryYearAwardNumResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	//data.List, err = userActivityDao.GetEveryYearDeclareAwardNum(schoolId)
	list, err := userActivityDao.GetEveryYearDeclareAwardNum(schoolId)
	if err != nil {
		return
	}
	activityDao := dao.NewActivityDao(s.appCtx)
	years, err := activityDao.GetYearsByGroup()
	if err != nil {
		return
	}
	i := 0
	for j := len(years) - 1; j >= 0; j-- {
		if i < len(list) && list[i].Year == years[j] {
			data.List = append(data.List, list[i])
			i++
		} else {
			data.List = append(data.List, dao.YearDeclareAwardNum{Year: years[j]})
		}
	}
	return
}

func (s *StatisticsService) GetEverySchoolAwardNum(year int) (data *resp.GetEverySchoolAwardNumResp, err error) {
	data = &resp.GetEverySchoolAwardNumResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	data.List, err = userActivityDao.GetEverySchoolAwardNum(year)
	if len(data.List) > 5 {
		list := data.List[5:]
		data.List = data.List[:5]
		var item dao.SchoolAwardNum
		item.SchoolId = "-1"
		item.SchoolName = "其他学校"
		for _, v := range list {
			item.AwardNum += v.AwardNum
		}
		data.List = append(data.List, item)
	}
	return
}

func (s *StatisticsService) GetEveryTeacherTypeAwardNum(year int) (data *resp.GetEveryTeacherTypeAwardNumResp, err error) {
	data = &resp.GetEveryTeacherTypeAwardNumResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	data.List, err = userActivityDao.GetEveryTeacherTypeAwardNum(year)
	return
}

func (s *StatisticsService) GetYearDeclareAwardRank(schoolId string) (data *resp.GetYearDeclareAwardRankResp, err error) {
	data = &resp.GetYearDeclareAwardRankResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	list, err := userActivityDao.GetDeclareAwardNumGroupByYear(schoolId)
	if err != nil {
		return
	}
	activityDao := dao.NewActivityDao(s.appCtx)
	years, err := activityDao.GetYearsByGroup()
	if err != nil {
		return
	}

	var schoolNumMap, teacherNumMap map[int]int
	if schoolId != "" {
		// 在岗教师数
		var numMap map[int]int
		teacherNumMap = make(map[int]int)
		numMap, err = dao.NewUserInfoDao(s.appCtx).GetEveryYearTeacherNumByGroup(schoolId)
		for y := range numMap {
			for _, year := range years {
				if y >= year {
					teacherNumMap[y] += numMap[year]
				}
			}
		}
	} else {
		// 参加学校数
		schoolNumMap, err = userActivityDao.GetEveryYearSchoolNumByGroup(years)
	}
	if err != nil {
		return
	}

	for i, j := 0, 0; i < len(years); i++ {
		item := resp.YearDeclareAwardRank{Year: years[i]}
		var cur dao.DeclareAwardRankNum
		if j < len(list) && list[j].Year == years[i] {
			cur = list[j].DeclareAwardRankNum
			j++
		}
		var last dao.DeclareAwardRankNum
		if j < len(list) && i+1 < len(years) && list[j].Year == years[i+1] {
			last = list[j].DeclareAwardRankNum
		}
		item.DeclareOnYear = util.OnYear(cur.DeclareNum, last.DeclareNum)
		item.AwardOnYear = util.OnYear(cur.AwardNum, last.AwardNum)
		item.AwardRate = util.Rate(cur.AwardNum, cur.DeclareNum)
		item.DeclareAwardRankNum = cur
		if schoolId != "" {
			item.Num = teacherNumMap[years[i]]
		} else {
			item.Num = schoolNumMap[years[i]]
		}
		data.List = append(data.List, item)
	}
	return
}

func (s *StatisticsService) GetSchoolDeclareAwardRank(params *req.GetSchoolDeclareAwardRankReq) (data *resp.GetSchoolDeclareAwardRankResp, err error) {
	data = &resp.GetSchoolDeclareAwardRankResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	total, list, err := userActivityDao.GetDeclareAwardRankGroupBySchool(params.Year, params.SchoolName, params.Page, params.Limit)
	if err != nil {
		return
	}
	data.Total = total
	activityDao := dao.NewActivityDao(s.appCtx)
	lastYear, err := activityDao.GetLastYear(params.Year)
	if err != nil {
		return
	}
	data.List = make([]resp.SchoolDeclareAwardRank, len(list))
	wg := util.NewWaitGroup(0)
	for i := range list {
		wg.Add()
		go func(index int) {
			defer wg.Done()
			cur := list[index]
			userActivity := dao.UserActivity{SchoolId: cur.SchoolId}
			lastDeclareNum, err := userActivityDao.GetDeclareNum(lastYear, userActivity)
			if err != nil {
				wg.SetError(err)
				return
			}
			lastAwardNum, err := userActivityDao.GetAwardNum(lastYear, userActivity)
			if err != nil {
				wg.SetError(err)
				return
			}
			data.List[index].DeclareAwardRankNum = cur.DeclareAwardRankNum
			data.List[index].SchoolId = cur.SchoolId
			data.List[index].SchoolName = cur.SchoolName
			data.List[index].DeclareOnYear = util.OnYear(cur.DeclareNum, lastDeclareNum)
			data.List[index].AwardOnYear = util.OnYear(cur.AwardNum, lastAwardNum)
			data.List[index].AwardRate = util.Rate(cur.AwardNum, cur.DeclareNum)
		}(i)
	}
	err = wg.Wait()
	return
}

func (s *StatisticsService) GetTeacherTypeDeclareAwardRank(params *req.GetTeacherTypeDeclareAwardRankReq) (data *resp.GetTeacherTypeDeclareAwardRankResp, err error) {
	data = &resp.GetTeacherTypeDeclareAwardRankResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	list, err := userActivityDao.GetDeclareAwardRankGroupByDeclareType(params.Year)
	if err != nil {
		return
	}

	activityDao := dao.NewActivityDao(s.appCtx)
	lastYear, err := activityDao.GetLastYear(params.Year)
	if err != nil {
		return
	}
	data.List = make([]resp.TeacherTypeDeclareAwardRank, len(list))
	wg := util.NewWaitGroup(0)
	for i := range list {
		wg.Add()
		go func(index int) {
			defer wg.Done()
			cur := list[index]
			userActivity := dao.UserActivity{DeclareType: cur.DeclareType}
			lastDeclareNum, err := userActivityDao.GetDeclareNum(lastYear, userActivity)
			if err != nil {
				wg.SetError(err)
				return
			}
			lastAwardNum, err := userActivityDao.GetAwardNum(lastYear, userActivity)
			if err != nil {
				wg.SetError(err)
				return
			}
			data.List[index].DeclareAwardRankNum = cur.DeclareAwardRankNum
			data.List[index].DeclareType = cur.DeclareType
			data.List[index].DeclareOnYear = util.OnYear(cur.DeclareNum, lastDeclareNum)
			data.List[index].AwardOnYear = util.OnYear(cur.AwardNum, lastAwardNum)
			data.List[index].AwardRate = util.Rate(cur.AwardNum, cur.DeclareNum)
		}(i)
	}
	err = wg.Wait()
	return
}

func (s *StatisticsService) GetResultGroupByDeclareType(activityId int) (data *resp.GetResultGroupByDeclareTypeResp, err error) {
	data = &resp.GetResultGroupByDeclareTypeResp{}
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	list, err := userActivityDao.GetResultGroupByDeclareType(activityId)
	if err != nil {
		return
	}
	for i, v := range list {
		data.SchoolNum += v.SchoolNum
		data.DeclareNum += v.DeclareNum
		list[i].RankNum = 0
		if v.RankNum > 0 {
			list[i].RankNum = 1
		}
	}
	data.List = make([]dao.ResultGroupByDeclareType, 6)
	k := 0
	for i := 0; i < 6; i++ {
		data.List[i].DeclareType = i + 1
		if k < len(list) && i+1 == list[k].DeclareType {
			data.List[i] = list[k]
			k++
		}
	}
	return
}
