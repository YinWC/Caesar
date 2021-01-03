package cores

import (
	"net/http"
	"net/url"
	"strings"

	"Caesar/internal/library/director"
	"Caesar/internal/relation"
	"Caesar/pkg/utils"
)

func CheckWaf(urlAddress string) bool {
	/*
	   该函数用来判断目标网站是否受WAF保护
	*/

	host := utils.GetNewHost(urlAddress)
	_, _, safeBody, _ := director.GenerateNormalGet(host + "?id=1")

	IpsWafCheckPayload := "AND 1=1 UNION ALL SELECT 1,NULL,'<script>alert(\"XSS\")</script>',table_name FROM information_schema.tables WHERE 2>1--/**/; EXEC xp_cmdshell('cat ../../../etc/passwd')#"

	AllPayLoad := url.PathEscape(IpsWafCheckPayload)

	var build strings.Builder
	build.WriteString(urlAddress)
	build.WriteString("?id=1&")
	build.WriteString(utils.GenRandString(4))
	build.WriteString("=")
	build.WriteString(AllPayLoad)
	s3 := build.String()

	//_, _, safeBody2, err := director.GenerateNormalGet(s3)
	_, _, safeBody2, err := director.GenerateNormalGet(s3)
	if err != nil {
		return true
	}

	if string(safeBody) == string(safeBody2) {
		// 如果两个页面相等，证明不存在waf
		return false
	}

	if utils.ComputeLevenshteinPercentage(string(safeBody), string(safeBody2)) > relation.Engine.UpperRatioBound {
		return false
	} else {
		return true
	}

}

func CheckCDN(urlAddress string) bool {
	/*
		检查地址是否是CDN
	*/

	return false
}

func CheckMVC(header http.Header, body []byte) (isMVC bool, frame string) {
	/*
	   该函数用来判断目标网站是否是MVC框架
	*/

	return true, ""

}
