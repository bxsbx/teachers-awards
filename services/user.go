package services

import (
	"context"
	"gorm.io/gorm"
	"teachers-awards/client"
	"teachers-awards/common/errorz"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	clientMode "teachers-awards/model/client"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type UserService struct {
	appCtx context.Context
}

func NewUserService(appCtx context.Context) *UserService {
	return &UserService{appCtx: appCtx}
}

func (s *UserService) GetUserInfo(userId string) (data *resp.GetUserInfoResp, err error) {
	data = &resp.GetUserInfoResp{}
	userInfoDao := dao.NewUserInfoDao(s.appCtx)
	err = userInfoDao.First(dao.UserInfo{UserId: userId}, &data.UserInfo)
	if err != nil {
		//找不到用户信息，返回中台录入的用户信息
		myError := err.(*errorz.MyError)
		if myError.Unwrap() == gorm.ErrRecordNotFound {
			userInfo := global.GetUserInfo(s.appCtx)
			platformClient := client.NewMiddlePlatformClientWithFrom(s.appCtx, userInfo.From)
			info, err := platformClient.GetTeacherDetailInfoById(userId)
			if err != nil {
				return nil, err
			}
			userInfo.UserName = info.Username

			data.UserSex = global.ZtSexToCurSex(info.SexCode)

			data.Phone = info.Phone
			data.Birthday = global.ZtBirthdayToCurBirthday(info.Birthday)

			data.SchoolId = info.OrganId
			data.SchoolName = info.OrganName
			data.Avatar = info.Avatar
			data.IdentityCard = global.ZtCardNumberToIdentityCard(info.CardTypeCode, info.CardNumber)

			data.UserId = userInfo.UserId
			data.UserName = userInfo.UserName
			return data, nil
		}
	}
	return
}

func (s *UserService) SaveUserInfo(params *req.SaveUserInfoReq) (err error) {
	userInfoDao := dao.NewUserInfoDao(s.appCtx)
	userInfo := dao.UserInfo{}
	util.ObjToObjByReflect(params, &userInfo)
	where := dao.UserInfo{UserId: params.UserId}
	count, err := userInfoDao.Count(where)
	if err != nil {
		return
	}
	nowTime := util.NowTime()
	if count > 0 { //更新
		userInfo.UpdateTime = &nowTime
		err = userInfoDao.UpdateByWhere(where, userInfo)
	} else { //添加
		userInfo.CreateTime = &nowTime
		userInfo.Year = nowTime.Year()
		userInfo.Role = global.Role2Teacher //开始时默认全部教师
		err = userInfoDao.Create(&userInfo)
	}
	return
}

func (s *UserService) GetUsersByName(params *req.GetUsersByNameReq) (data *resp.GetUsersByNameResp, err error) {
	data = &resp.GetUsersByNameResp{}
	platformClient := client.NewMiddlePlatformClient(s.appCtx)
	userList, total, err := platformClient.GetTeachersByName(params.UserName, params.Page, params.Limit)
	userIds := util.ListObjToListObj(userList, func(obj clientMode.TeacherDetailInfo) string {
		return obj.PersonId
	})
	var users []dao.UserInfo
	err = dao.NewUserInfoDaoWithDB(global.GormDB.Where("user_id in (?)", userIds), s.appCtx).
		Find(nil, &users, "user_id,role")
	userRoleMap := util.ListObjToMap(users, func(obj dao.UserInfo) (string, int) {
		return obj.UserId, obj.Role
	})
	data.UserList = make([]resp.UserInfo, len(userList))
	for i, v := range userList {
		data.UserList[i].UserId = v.PersonId
		data.UserList[i].UserName = v.Username
		data.UserList[i].UserSex = global.ZtSexToCurSex(v.SexCode)
		data.UserList[i].IdentityCard = global.ZtCardNumberToIdentityCard(v.CardTypeCode, v.CardNumber)
		data.UserList[i].Phone = v.Phone
		data.UserList[i].SchoolId = v.OrganId
		data.UserList[i].SchoolName = v.OrganName
		if roleValue, ok := userRoleMap[v.PersonId]; ok {
			data.UserList[i].Roles = global.RoleValToRoles(roleValue)
		} else {
			data.UserList[i].Roles = []int{global.RoleTeacher}
		}

	}
	data.Total = total
	return
}

func (s *UserService) SetRoleToUser(params *req.SetRoleToUserReq) (err error) {
	userInfoDao := dao.NewUserInfoDaoWithDB(global.GormDB.Where("user_id in (?)", params.UserIds), s.appCtx)
	var existUserIds []string
	err = userInfoDao.Pluck(nil, &existUserIds, "user_id")
	if err != nil {
		return
	}
	userIdsSet := util.ListObjToMap(existUserIds, func(obj string) (string, struct{}) {
		return obj, struct{}{}
	})
	var noExistUserIds []string
	for _, userId := range params.UserIds {
		if _, ok := userIdsSet[userId]; !ok {
			noExistUserIds = append(noExistUserIds, userId)
		}
	}
	//过滤不合法的角色
	var roles []int
	for _, v := range params.Roles {
		if _, ok := global.RoleNameMap[v]; ok {
			roles = append(roles, v)
		}
	}
	roleVal := global.RolesToRoleVal(params.Roles)
	//默认加上教师角色
	if roleVal&global.Role2Teacher == 0 {
		roleVal += global.Role2Teacher
	}
	platformClient := client.NewMiddlePlatformClient(s.appCtx)
	createUserList := make([]dao.UserInfo, 0)
	nowTime := util.NowTime()
	exportAuth := 0
	if roleVal&global.Role2Expert > 0 {
		exportAuth = global.ExportAuthNo
	}
	if len(noExistUserIds) > 0 {
		wg := util.NewWaitGroup(20)
		for _, userId := range noExistUserIds {
			wg.Add()
			go func(teacherId string) {
				defer wg.Done()
				info, err := platformClient.GetTeacherDetailInfoById(teacherId)
				if err != nil {
					wg.SetError(err)
					return
				}
				data := dao.UserInfo{}
				data.UserId = info.PersonId
				data.UserName = info.Username
				data.UserSex = global.ZtSexToCurSex(info.SexCode)
				data.IdentityCard = global.ZtCardNumberToIdentityCard(info.CardTypeCode, info.CardNumber)
				data.Birthday = global.ZtBirthdayToCurBirthday(info.Birthday)
				data.Phone = info.Phone
				data.SchoolId = info.OrganId
				data.SchoolName = info.OrganName
				data.Avatar = info.Avatar
				data.Role = roleVal
				data.CreateTime = &nowTime
				data.ExportAuth = exportAuth
				data.Year = nowTime.Year()
				wg.Lock()
				createUserList = append(createUserList, data)
				wg.Unlock()
			}(userId)
		}
		err = wg.Wait()
		if err != nil {
			return
		}
	}

	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		//存在的更新
		if len(existUserIds) > 0 {
			updateMap := make(map[string]interface{})
			updateMap["role"] = roleVal
			err := dao.NewUserInfoDaoWithDB(tx.Where("user_id in (?)", existUserIds), s.appCtx).UpdateByWhere(nil, updateMap)
			if err != nil {
				return err
			}
		}
		//不存在则插入
		if len(createUserList) > 0 {
			err := dao.NewUserInfoDaoWithDB(tx, s.appCtx).BatchInsert(createUserList)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *UserService) GetUserListByWhere(params *req.GetUserListByWhereReq) (data *resp.GetUserListByWhereResp, err error) {
	data = &resp.GetUserListByWhereResp{}
	db := global.GormDB
	if params.UserName != "" {
		db = db.Where("user_name like ?", "%"+params.UserName+"%")
	}
	if params.SchoolName != "" {
		db = db.Where("school_name like ?", "%"+params.SchoolName+"%")
	}
	if params.Role == 0 {
		roleVal := global.RolesToRoleVal([]int{global.RoleSchool, global.RoleExpert, global.RoleEdb, global.RoleAdmin})
		db = db.Where("role & ? > 0", roleVal)
	} else {
		roleVal := global.RolesToRoleVal([]int{params.Role})
		db = db.Where("role & ? = ?", roleVal, roleVal)
	}
	userInfoDao := dao.NewUserInfoDaoWithDB(db, s.appCtx)

	where := dao.UserInfo{
		UserSex:      params.UserSex,
		IdentityCard: params.IdentityCard,
		Phone:        params.Phone,
		SchoolId:     params.SchoolId,
	}
	data.Total, err = userInfoDao.FindAndCountWithPageOrder(where, &data.List, params.Page, params.Limit, "")
	for i, v := range data.List {
		data.List[i].Roles = global.RoleValToRoles(v.Role)
	}
	return
}

func (s *UserService) GetExpertAuthListByWhere(params *req.GetExpertAuthListByWhereReq, twoIds []int) (data *resp.GetExpertAuthListByWhereResp, err error) {
	data = &resp.GetExpertAuthListByWhereResp{}
	db := global.GormDB
	if len(twoIds) > 0 {
		db = db.Where("user_id in (select user_id from expert_auth_indicator where two_indicator_id in (?) group by user_id)", twoIds)
	}
	if params.UserName != "" {
		db = db.Where("user_name like ?", "%"+params.UserName+"%")
	}
	//2:专家，过滤专家
	db = db.Where("role & ? > 0", global.RoleExpert).Preload("ExpertAuthIndicatorList")
	where := dao.UserInfo{
		UserSex:      params.UserSex,
		IdentityCard: params.IdentityCard,
		Phone:        params.Phone,
		ExportAuth:   params.ExportAuth,
		AuthDay:      params.AuthDay,
	}
	userInfoDao := dao.NewUserInfoDaoWithDB(db, s.appCtx)
	var list []dao.UserInfo
	data.Total, err = userInfoDao.FindAndCountWithPageOrder(where, &list, params.Page, params.Limit, "export_auth asc,auth_day desc")
	if err != nil {
		return
	}
	//组装指标数据
	var twoIdsFromList []int
	for _, v1 := range list {
		for _, v2 := range v1.ExpertAuthIndicatorList {
			twoIdsFromList = append(twoIdsFromList, v2.TwoIndicatorId)
		}
	}
	twoIdsFromList = util.RemoveRepeatFromList(twoIdsFromList)
	twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(global.GormDB.Where("two_indicator_id in (?)", twoIdsFromList), s.appCtx)
	var twoList []dao.TwoIndicator
	err = twoIndicatorDao.Find(nil, &twoList)
	if err != nil {
		return
	}
	oneIds := util.ListToDeduplicationList(twoList, func(t dao.TwoIndicator) (int, int) {
		return t.OneIndicatorId, t.OneIndicatorId
	})
	oneIndicatorDao := dao.NewOneIndicatorDaoWithDB(global.GormDB.Where("one_indicator_id in (?)", oneIds), s.appCtx)
	var oneList []dao.OneIndicator
	err = oneIndicatorDao.Find(nil, &oneList, "one_indicator_id,one_indicator_name")
	if err != nil {
		return
	}
	twoMap := util.ListObjToMap(twoList, func(obj dao.TwoIndicator) (int, dao.TwoIndicator) {
		return obj.TwoIndicatorId, obj
	})
	oneMap := util.ListObjToMap(oneList, func(obj dao.OneIndicator) (int, dao.OneIndicator) {
		return obj.OneIndicatorId, obj
	})

	//组装数据
	data.List = make([]resp.ExpertAuthList, len(list))
	for i, v := range list {
		data.List[i] = resp.ExpertAuthList{
			UserId:       v.UserId,
			UserName:     v.UserName,
			UserSex:      v.UserSex,
			IdentityCard: v.IdentityCard,
			Phone:        v.Phone,
			ExportAuth:   v.ExportAuth,
			AuthDay:      v.AuthDay,
		}
		oneIndexMap := make(map[int]int)
		var ones []resp.OneIndicatorOnlyName
		for _, w := range v.ExpertAuthIndicatorList {
			two := twoMap[w.TwoIndicatorId]
			oneIndex := 0
			if index, ok := oneIndexMap[two.OneIndicatorId]; ok {
				oneIndex = index
			} else {
				oneIndexMap[two.OneIndicatorId] = len(ones)
				oneIndex = len(ones)
				one := oneMap[two.OneIndicatorId]
				ones = append(ones, resp.OneIndicatorOnlyName{OneIndicatorId: one.OneIndicatorId, OneIndicatorName: one.OneIndicatorName})
			}
			ones[oneIndex].TwoList = append(ones[oneIndex].TwoList, resp.TwoIndicatorOnlyName{TwoIndicatorId: two.TwoIndicatorId, TwoIndicatorName: two.TwoIndicatorName})
		}
		data.List[i].ExpertAuthIndicatorList = ones
	}
	return
}

func (s *UserService) SetExpertAuth(params *req.SetExpertAuthReq) (err error) {
	//过滤掉数据库中不存在的twoId
	twoIndicatorDao := dao.NewTwoIndicatorDaoWithDB(global.GormDB.Where("two_indicator_id in (?)", params.TwoIds), s.appCtx)
	var twoIds []int
	err = twoIndicatorDao.Pluck(nil, &twoIds, "two_indicator_id")
	if err != nil {
		return
	}
	if len(twoIds) <= 0 {
		return
	}
	//获取用户存在的twoId
	var list []dao.ExpertAuthIndicator
	expertAuthIndicatorDao := dao.NewExpertAuthIndicatorDao(s.appCtx)
	err = expertAuthIndicatorDao.Find(dao.ExpertAuthIndicator{UserId: params.UserId}, &list)
	if err != nil {
		return
	}
	twoIdSet := util.ListObjToMap(list, func(obj dao.ExpertAuthIndicator) (int, struct{}) {
		return obj.TwoIndicatorId, struct{}{}
	})
	if len(twoIdSet) <= 0 {
		nowDay := util.NowDate()
		err = dao.NewUserInfoDao(s.appCtx).UpdateByWhere(dao.UserInfo{UserId: params.UserId}, dao.UserInfo{AuthDay: &nowDay, ExportAuth: global.ExportAuthYes})
		if err != nil {
			return
		}
	}
	var addList []dao.ExpertAuthIndicator
	existsIdMap := make(map[int]interface{})
	for _, twoId := range twoIds {
		if _, ok := twoIdSet[twoId]; !ok {
			addList = append(addList, dao.ExpertAuthIndicator{UserId: params.UserId, TwoIndicatorId: twoId})
		} else {
			existsIdMap[twoId] = struct{}{}
		}
	}
	err = expertAuthIndicatorDao.BatchInsert(addList)
	if err != nil {
		return
	}
	var delIds []int
	for twoId := range twoIdSet {
		if _, ok := existsIdMap[twoId]; !ok {
			delIds = append(delIds, twoId)
		}
	}
	err = expertAuthIndicatorDao.DeleteByWhere("user_id = ? and two_indicator_id in (?)", params.UserId, delIds)
	return
}

func (s *UserService) CancelExpertAuth(userId string) (err error) {
	err = global.GormDB.Transaction(func(tx *gorm.DB) error {
		err := dao.NewExpertAuthIndicatorDaoWithDB(tx, s.appCtx).DeleteByWhere(dao.ExpertAuthIndicator{UserId: userId})
		if err != nil {
			return err
		}
		userInfoDao := dao.NewUserInfoDaoWithDB(tx, s.appCtx)
		updateMap := make(map[string]interface{})
		updateMap["auth_day"] = nil
		updateMap["export_auth"] = global.ExportAuthNo
		err = userInfoDao.UpdateByWhere(dao.UserInfo{UserId: userId}, updateMap)
		return err
	})
	return
}
