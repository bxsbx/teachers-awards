package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"teachers-awards/common/util"
	"text/template"
)

const (
	BasePath        = "H:\\dream\\teachers-awards\\"
	ProjectName     = "teachers-awards"
	ControllersPath = BasePath + "controllers\\"
	ServicesPath    = BasePath + "services\\"
	RouterGroupPath = BasePath + "router\\"

	TypePath = BasePath + "model\\req\\"

	Produce = "json"

	ControllersTmpl = BasePath + "generate\\template\\controller.tmpl"
	ServicesTmpl    = BasePath + "generate\\template\\service.tmpl"
	RouterTmpl      = BasePath + "generate\\template\\router.tmpl"
)

var ParamTypeMap = map[string]string{
	"path":  "path",
	"query": "query",
	"form":  "formData",
	"json":  "body",
}

func main() {
	myTemplate := myTemplate{
		ProjectName: ProjectName,
		GroupName:   "User",
		FunName:     "CancelExpertAuth",
		ReqModel:    "CancelExpertAuthReq",
		RespModel:   "",
		DataType:    "",
		Router:      "/v1/user/cancel/expert/auth",
		Method:      http.MethodPost,
	}
	swagger := Swagger{
		Summary: "取消专家授权",
		Tags:    myTemplate.GroupName,
		Produce: Produce,
	}

	myTemplate.SlideGroupName = util.SlideNaming(myTemplate.GroupName)
	myTemplate.LowerFunName = util.FirstLower(myTemplate.FunName)
	myTemplate.LowerGroupName = util.FirstLower(myTemplate.GroupName)

	swagger.Params = GetParamsFromReqType(myTemplate)
	//if swagger.Params == nil {
	//	log.Fatal("请求类型不存在")
	//}
	myTemplate.Swagger = swagger

	ControllersFilePath := ControllersPath + myTemplate.SlideGroupName + ".go"
	isControllerFile, _ := util.FileIsExist(ControllersFilePath)
	myTemplate.IsControllerFile = !isControllerFile

	ServicesFilePath := ServicesPath + myTemplate.SlideGroupName + ".go"
	isServiceFile, _ := util.FileIsExist(ServicesFilePath)
	myTemplate.IsServiceFile = !isServiceFile

	RouterFilePath := RouterGroupPath + myTemplate.SlideGroupName + ".go"
	isRouterFile, _ := util.FileIsExist(RouterFilePath)
	myTemplate.IsRouterFile = !isRouterFile

	funcName := func(line string) bool {
		if len(line) > 4 && line[:4] == "func" && strings.Contains(line, myTemplate.FunName+"(") {
			fmt.Println("方法名已存在：" + myTemplate.FunName)
			return true
		}
		return false
	}

	routerPath := func(line string) bool {
		if strings.Contains(line, myTemplate.Router) {
			fmt.Println("路由已存在：" + myTemplate.Router)
			return true
		}
		return false
	}
	GeneralContentToFile(ControllersTmpl, ControllersFilePath, "", myTemplate, funcName)
	GeneralContentToFile(ServicesTmpl, ServicesFilePath, "", myTemplate, funcName)
	GeneralContentToFile(RouterTmpl, RouterFilePath, "router general tag", myTemplate, routerPath)

	//swag注释格式化
	exec.Command("swag", "fmt", "-d", ControllersFilePath).Run()
	err := exec.Command("swag", "init", "-dir", BasePath, "-output", BasePath+"swagger\\docs", "-exclude", BasePath+"global").Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("执行成功")
}

type Swagger struct {
	Summary string
	Tags    string
	Produce string
	Params  []Params
}
type Params struct {
	Name        string
	ParamType   string
	DataType    string
	IsNeed      bool
	Description string
}

type myTemplate struct {
	ProjectName      string
	GroupName        string
	SlideGroupName   string
	LowerGroupName   string
	FunName          string
	LowerFunName     string
	DataType         string
	ReqModel         string
	RespModel        string
	Swagger          Swagger
	Router           string
	Method           string
	ValidateType     string
	IsRouterFile     bool
	IsControllerFile bool
	IsServiceFile    bool
}

// 从请求类型中获取Swagger信息
func GetParamsFromReqType(myTemplate myTemplate) []Params {
	var list []Params
	if myTemplate.DataType == "body" {
		list = append(list, Params{
			Name:        "data",
			ParamType:   "body",
			DataType:    "req." + myTemplate.ReqModel,
			IsNeed:      true,
			Description: "body请求体",
		})
		return list
	}

	reqTypeStr := "type " + myTemplate.ReqModel + " struct {"
	bytes, err := os.ReadFile(TypePath + myTemplate.SlideGroupName + ".go")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(bytes), "\n")

	flag := false
	for _, line := range lines {
		if strings.Contains(line, reqTypeStr) {
			flag = true
		} else if len(line) > 0 && line[:1] == "}" && flag {
			break
		} else if flag {
			var params Params
			strs := regexp.MustCompile("\\s+").Split(strings.Trim(line, " |\t"), -1)
			params.DataType = strs[1]
			if strings.Contains(params.DataType, "time.Time") {
				params.DataType = "string"
			}
			if strings.Contains(line, "required") {
				params.IsNeed = true
			}

			tags := strings.Trim(regexp.MustCompile("`.*`").FindString(line), "`")
			params.Name = strings.Trim(regexp.MustCompile("\"\\w*\"").FindString(tags), "\"")
			params.ParamType = ParamTypeMap[strings.Split(tags, ":")[0]]
			if myTemplate.Method == http.MethodGet || myTemplate.Method == http.MethodDelete || params.ParamType == "" {
				params.ParamType = "query"
			}
			params.Description = strings.Trim(strings.Trim(regexp.MustCompile("//.*").FindString(line), "//"), " ")
			list = append(list, params)
		}
	}

	return list
}

func GeneralContentToFile(tmpl, filePath, location string, myTemplate myTemplate, f func(l string) bool) {
	exist, err := util.FileIsExist(filePath)
	if err != nil {
		log.Fatal(err)
	}
	if exist {
		content, err := util.ReadFileContent(filePath)
		if err != nil {
			log.Fatal(err)
		}
		//方法已存在则跳过
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if f(line) {
				return
			}
		}
	}

	t := template.Must(template.ParseFiles(tmpl))
	var builder strings.Builder
	err = t.Execute(&builder, myTemplate)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(builder.String())
	if location != "" {
		err = util.AddAtCustomLocation(filePath, builder.String(), location)
	} else {
		err = util.WriteToFile(filePath, builder.String())
	}
	if err != nil {
		log.Fatal(err)
	}
}
