package desens

import "strings"

// DesensitizationEmail 邮箱脱敏
func DesensitizationEmail(email string) string {
	// 2933756974@qq.com  2****@qq.com
	// yubiniom@yaho.com y****@yaho.com
	eList := strings.Split(email, "@")
	if len(eList) != 2 {
		return ""
	}
	return eList[0][:1] + "****@" + eList[1]
}

// DesensitizationTel 手机号脱敏
func DesensitizationTel(tel string) string {
	// 150 7162 1639
	// 150 **** 1639
	if len(tel) != 11 {
		return ""
	}

	return tel[:3] + "****" + tel[7:] // 每一数的后面切分别在第三个和第七个数字后面切
}
