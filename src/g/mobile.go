// 手机号码 号段 记录 

package g

import (
	//"encoding/json"
	//"log"
	"regexp"
)

const (
	ChinaMobile		= 1   // 移动
	ChinaUnicom		= 2   // 联通
	ChinaTelecom	= 3   // 电信
	ChinaOther		= 4   // 虚拟   170 
	Unknow 			= 0   // 不支持的号码 
)


// 号码归属地判断 
func MobileBelongTo(mobile string) int {
	mob := mobile[0:3]

	var r1 *regexp.Regexp
	var r2 *regexp.Regexp
	var r3 *regexp.Regexp
	var r4 *regexp.Regexp

	if Config().Haoduan.ChinaMobile != "" {
		r1 = regexp.MustCompile(Config().Haoduan.ChinaMobile)
	}

	if Config().Haoduan.ChinaUnicom != "" {
		r2 = regexp.MustCompile(Config().Haoduan.ChinaUnicom)
	}

	if Config().Haoduan.ChinaTelecom != "" {
		r3 = regexp.MustCompile(Config().Haoduan.ChinaTelecom)
	}

	if Config().Haoduan.ChinaOther != "" {
		r4 = regexp.MustCompile(Config().Haoduan.ChinaOther)
	}


	if r1 != nil && r1.MatchString(mob) == true {
		return ChinaMobile
	}

	if r2 != nil && r2.MatchString(mob) == true {
		return ChinaUnicom
	}

	if r3 != nil && r3.MatchString(mob) == true {
		return ChinaTelecom
	}

	if r4 != nil && r4.MatchString(mob) == true {
		return ChinaOther
	}

	return Unknow
}
