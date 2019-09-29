package ging

import (
	"fmt"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * setting util
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * App 域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s DomainOption) AppDomain() string {
	host := glib.FilterHostProtocol(s.AppHost)

	proto := "http"
	if s.IsSsl {
		proto = "https"
	}

	domain := fmt.Sprintf("%s://%s", proto, host)

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Image 域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s DomainOption) ImageDomain() string {
	host := glib.FilterHostProtocol(s.ImageHost)

	proto := "http"
	if s.IsSsl {
		proto = "https"
	}

	domain := fmt.Sprintf("%s://%s", proto, host)

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Audio 域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s DomainOption) AudioDomain() string {
	host := glib.FilterHostProtocol(s.AudioHost)

	proto := "http"
	if s.IsSsl {
		proto = "https"
	}

	domain := fmt.Sprintf("%s://%s", proto, host)

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Video 域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s DomainOption) VideoDomain() string {
	host := glib.FilterHostProtocol(s.VideoHost)

	proto := "http"
	if s.IsSsl {
		proto = "https"
	}

	domain := fmt.Sprintf("%s://%s", proto, host)

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * File 域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s DomainOption) FileDomain() string {
	host := glib.FilterHostProtocol(s.FileHost)

	proto := "http"
	if s.IsSsl {
		proto = "https"
	}

	domain := fmt.Sprintf("%s://%s", proto, host)

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取阿里云网关域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *AliyunSms) Domain() string {
	domain := fmt.Sprintf("%s://%s", "http", s.Gateway)

	if s.IsSsl {
		domain = fmt.Sprintf("%s://%s", "https", s.Gateway)
	}

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取阿里大鱼网关域名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *AlidayuSms) Domain() string {
	domain := fmt.Sprintf("%s://%s", "http", s.Gateway)

	if s.IsSsl {
		domain = fmt.Sprintf("%s://%s", "https", s.Gateway)
	}

	return domain
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 判断用户id是否存在
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s Blackwhite) IsInUsers(userId string) bool {
	isInUsers := false

	for _, _userId := range s.Users {
		if _userId == userId {
			isInUsers = true
			break
		}
	}

	return isInUsers
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 判断ip是否存在
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s Blackwhite) IsInIps(ip string) bool {
	isInIps := false

	if ips := glib.StringToStringSlice(ip, ":"); len(ips) > 1 {
		ip = ips[0]
	}

	for _, _ip := range s.Ips {
		if _ip == ip {
			isInIps = true
			break
		}
	}

	return isInIps
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 判断当前时刻是否属于限制时段中
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s TimeLimit) IsTime() bool {
	isTime := false
	if !s.IsDisabled {
		if s.Opening >= 0 && s.Closing >= 0 {
			hour := glib.GetCurrentHour()
			minute := glib.GetCurrentMinute()

			totalMinute := hour*60 + minute

			if totalMinute >= int32(s.Opening) && totalMinute <= int32(s.Closing) {
				isTime = true
			}
		}
	}

	return isTime
}
