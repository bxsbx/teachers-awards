package services

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"teachers-awards/client"
	"teachers-awards/common/util"
	"teachers-awards/dao"
	"teachers-awards/global"
	"teachers-awards/model/req"
	"teachers-awards/model/resp"
)

type OtherService struct {
	appCtx context.Context
}

func NewOtherService(appCtx context.Context) *OtherService {
	return &OtherService{appCtx: appCtx}
}

func (s *OtherService) GetUploadToken(userId string, fileName string) (data *resp.GetUploadTokenResp, err error) {
	nowTime := util.NowTime()
	data = &resp.GetUploadTokenResp{}
	// 构建资源存储路径(存储key)
	storagePath := fmt.Sprintf("teachers-awards/%v/%v/%v/%v", userId, nowTime.Format("2006/01/02"), nowTime.Unix(), fileName)

	qiNiuConfig := global.ServeCfg.QiNiu
	putPolicy := storage.PutPolicy{ //*七牛云上传策略*
		Scope:   qiNiuConfig.Bucket, // 指定上传的目标资源空间 Bucket
		Expires: 300,                // 5分钟有效期
		//InsertOnly: 1,                                                                    // 仅能以新增模式上传文件
		ReturnBody: `{"F_url":"` + qiNiuConfig.Domain + `/` + storagePath + `"}`, // 自定义七牛云最终返回給上传端的数据
		SaveKey:    storagePath,                                                  // 自定义资源名
		MimeLimit:  "image/*;application/*",                                      //限制上传类型
		FsizeLimit: 10 * 1024 * 1024,                                             //限制上传大小
	}
	mac := qbox.NewMac(qiNiuConfig.AccessKey, qiNiuConfig.SecretKey) //*accessKey和secretKey*
	data.Token = putPolicy.UploadToken(mac)                          //*生成上传token*
	return
}

// 记录uai的操作日记
func (s *OtherService) OperationUaiRecord(uaiId int64, operationType, operationRole int) (err error) {
	twoIndicatorInfo, err := dao.NewUserActivityIndicatorDao(s.appCtx).GetIndicatorInfo(uaiId)
	if err != nil {
		return err
	}
	description := fmt.Sprintf("<%v-%v>", twoIndicatorInfo.OneIndicatorName, twoIndicatorInfo.TwoIndicatorName)
	switch operationType {
	case 1:
		description = "添加申报" + description
	case 2:
		description = "修改申报" + description
	case 3:
		description = "撤销申报" + description
	}
	err = dao.NewOperationRecordDao(s.appCtx).InsertOperationRecordToUAI(uaiId, 1, operationType, operationRole, description)
	return
}

// 记录uai的操作日记（是否通过）
func (s *OtherService) OperationUaiPassRecord(uaiId int64, isPass int, role int) (err error) {
	twoIndicatorInfo, err := dao.NewUserActivityIndicatorDao(s.appCtx).GetIndicatorInfo(uaiId)
	if err != nil {
		return err
	}
	description := fmt.Sprintf("<%v-%v>", twoIndicatorInfo.OneIndicatorName, twoIndicatorInfo.TwoIndicatorName)
	switch isPass {
	case 0:
		description = "不通过申报" + description
	case 1:
		description = "通过申报" + description
	}
	err = dao.NewOperationRecordDao(s.appCtx).InsertOperationRecordToUAI(uaiId, 1, 2, role, description)
	return
}

// 记录uai更改指标的记录
func (s *OtherService) UpdateUaiRecord(uaiId, uaId int64, oldTwoId, newTwoId int) (err error) {
	userActivityDao := dao.NewUserActivityDao(s.appCtx)
	oldInfo, err := userActivityDao.GetIndicatorInfo(uaId, oldTwoId)
	if err != nil {
		return err
	}
	newInfo, err := userActivityDao.GetIndicatorInfo(uaId, newTwoId)
	if err != nil {
		return err
	}
	description := fmt.Sprintf("更改申报<%v-%v>为<%v-%v>", oldInfo.OneIndicatorName, oldInfo.TwoIndicatorName, newInfo.OneIndicatorName, newInfo.TwoIndicatorName)
	err = dao.NewOperationRecordDao(s.appCtx).InsertOperationRecordToUAI(uaiId, 1, 2, 3, description)
	return
}

func (s *OtherService) GetOperationRecordFromUai(userActivityIndicatorId int64) (data *resp.GetOperationRecordFromUaiResp, err error) {
	data = &resp.GetOperationRecordFromUaiResp{}
	operationRecordDao := dao.NewOperationRecordDao(s.appCtx)
	err = operationRecordDao.Find(dao.OperationRecord{RelationalId: userActivityIndicatorId, RelationalType: 1}, &data.List)
	return
}

// 获取学科枚举
func (s *OtherService) GetSubjectEnum(refresh int) (subjectEnumMap map[string]string, err error) {
	f := func() (map[string]string, error) {
		platformClient := client.NewMiddlePlatformClient(s.appCtx)
		enums, err := platformClient.GetEnumInfo("discipline")
		if err != nil {
			return nil, err
		}
		enumMap := make(map[string]string)
		for _, v := range enums {
			enumMap[v.Code] = v.Name
		}
		global.SubjectEnumMap.SubjectMap = enumMap
		return enumMap, nil
	}

	if refresh == 1 {
		return f()
	}

	// 使用双重判断加锁来获取（动态获取）
	enum := global.SubjectEnumMap
	if enum.SubjectMap == nil {
		enum.Lock()
		defer enum.Unlock()
		if enum.SubjectMap == nil {
			return f()
		} else {
			subjectEnumMap = enum.SubjectMap
		}
	} else {
		subjectEnumMap = enum.SubjectMap
	}
	return
}

func (s *OtherService) GetSchoolListByPage(params *req.GetSchoolListByPageReq) (data *resp.GetSchoolListByPageResp, err error) {
	data = &resp.GetSchoolListByPageResp{}
	platformClient := client.NewMiddlePlatformClient(s.appCtx)
	schoolList, err := platformClient.GetSchoolListByPage(params.Page, params.Limit)
	if err != nil {
		return
	}
	data.Total = schoolList.Total
	data.SchoolList = make([]resp.School, len(schoolList.Data))
	for i, v := range schoolList.Data {
		data.SchoolList[i].SchoolId = v.Id
		data.SchoolList[i].SchoolName = v.Name
	}
	return
}

func (s *OtherService) GetSchoolListByName(schoolName string) (data *resp.GetSchoolListByNameResp, err error) {
	data = &resp.GetSchoolListByNameResp{}
	platformClient := client.NewMiddlePlatformClient(s.appCtx)
	schoolList, err := platformClient.GetSchoolListByName(schoolName)
	data.SchoolList = make([]resp.School, len(schoolList))
	for i, v := range schoolList {
		data.SchoolList[i].SchoolId = v.SchoolId
		data.SchoolList[i].SchoolName = v.SchoolName
	}
	return
}
