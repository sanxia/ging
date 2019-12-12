package ging

import (
	"fmt"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * setting
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Setting struct {
		WorkerNodeId int64
		Server       ServerOption
		Domain       DomainOption
		Image        ImageOption
		Register     RegisterOption
		Logon        LogonOption
		Session      SessionOption
		Cache        CacheOption
		MessageQueue MessageQueueOption
		Database     DatabaseOption
		Security     SecurityOption
		Oauth        OauthOption
		Pay          PayOption
		Mail         MailOption
		Sms          SmsOption
		Storage      StorageOption
		Search       SearchOption
		Log          LogOption
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * server option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	ServerOption struct {
		Host  string
		Ports []int
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * domain option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DomainOption struct {
		AppHost    string
		ImageHost  string
		AudioHost  string
		VideoHost  string
		FileHost   string
		Logo       string
		Seo        DomainSeo
		Icp        DomainIcp
		Template   DomainTemplate
		Copyright  string
		Version    string
		StatusCode string
		IsSsl      bool
		IsTest     bool
		IsDebug    bool
		IsOnline   bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * domain seo
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DomainSeo struct {
		Author      string `form:"author" json:"author"`
		Title       string `form:"title" json:"title"`
		Keywords    string `form:"keywords" json:"keywords"`
		Description string `form:"description" json:"description"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * domain icp
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DomainIcp struct {
		IcpNo     string `form:"icp_no" json:"icp_no"`
		RecordNo  string `form:"record_no" json:"record_no"`
		RecordUrl string `form:"record_url" json:"record_url"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * html template
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DomainTemplate struct {
		Path       string   `json:"path"`       //html template path
		Extensions []string `json:"extensions"` //extension collection
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * image option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	ImageOption struct {
		Avatar  Avatar //user profile picture
		Default string //default picture
		Font    string //font path
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * user profile picture
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	Avatar struct {
		Male   string `json:"male"`
		Female string `json:"female"`
		Other  string `json:"other"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * sign up option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	RegisterOption struct {
		Username          RegisterRule       //username rules
		Password          RegisterRule       //password rules
		Invitation        RegisterInvitation //invitation
		Ip                IpLimit            //same ip limit
		Time              TimeLimit          //time limit
		IsConfirmPassword bool               //need confirm password
		IsCaptcha         bool               //human-machine verification code
		IsApprove         bool               //need approve
		IsUsername        bool               //allow your username to register
		IsMobile          bool               //allow your phone to register
		IsDisabled        bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * sign up rule
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	RegisterRule struct {
		Rule       string   `json:"rule"`       //regular
		Forbiddens []string `json:"forbiddens"` //prohibited collections
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * sign up for an invitation
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	RegisterInvitation struct {
		Count      int32 `json:"count"`       //maximum number of invitations
		IsDisabled bool  `json:"is_disabled"` //is disabled
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * same Ip limit
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	IpLimit struct {
		Minute     int  `json:"minute"`      //same ip（minutes）
		Count      int  `json:"count"`       //same ip max count
		IsDisabled bool `json:"is_disabled"` //is disabled
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * time limit
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	TimeLimit struct {
		Opening    int  `json:"opening"` //start（minutes）
		Closing    int  `json:"closing"` //end（minutes）
		IsDisabled bool `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * sign-in option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	LogonOption struct {
		AuthorizeUrl string      //certified url
		DefaultUrl   string      //default URL after successful certification
		PassUrls     []string    //skip the certified URL collection
		Cookie       Cookie      //cookie
		Ip           IpLimit     //same ip limit
		Time         TimeLimit   //time limit
		Test         TestAccount //test account
		ErrorCount   int32       //Number of login errors
		LockMinute   int32       //Number of lock minutes after error symup
		IsCookie     bool        //After the login is successful, the response header returns the token string, otherwise the client cookie is set
		IsSliding    bool        //sliding expiration date
		IsCaptcha    bool        //human-machine verification code
		IsMobile     bool        //allow your phone to sign in
		IsUsername   bool        //allow user name login
		IsDisabled   bool        //logon is closed
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * test account
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	TestAccount struct {
		Mobiles    []string `json:"mobiles"`
		MobileCode string   `json:"mobile_code"`
		IsDisabled bool     `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * session option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SessionOption struct {
		Cookie     Cookie      //cookie
		Redis      RedisServer //redis server
		IsRedis    bool        //false:cookie | true:redis
		IsDisabled bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * cache option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	CacheOption struct {
		Prefix     string
		Redis      RedisServer //redis server
		Expire     int         //default expiration date（seconds）
		IsRedis    bool        //false:memory | true:redis
		IsDisabled bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * redis server
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	RedisServer struct {
		Ip       string
		Port     int
		Password string
		Timeout  int
		Db       int //redis db index
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * message queue option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	MessageQueueOption struct {
		Server MessageQueueServerOption
		Vhost  string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * message queue server
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	MessageQueueServerOption struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * database option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DatabaseOption struct {
		Connections []DatabaseConnection
		IsLog       bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * database connection
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DatabaseConnection struct {
		Key           string
		Database      string
		Dialect       string
		ShardingCount int32
		Servers       []DatabaseServer
	}

	DatabaseServer struct {
		Index    int32
		Username string
		Password string
		Host     string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * security option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SecurityOption struct {
		ApiSecret     string     //api key
		EncryptSecret string     //secure encryption Key
		CaptchaExpire int        //Man-machine check expiration time（minutes）
		Black         Blackwhite //blacklist
		White         Blackwhite //whitelist
		Ip            IpLimit
		InTime        TimeLimit //write time
		OutTime       TimeLimit //out time
		Cors          CorsDomain
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * blackwhite list
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	Blackwhite struct {
		Users      []string `json:"users"` //user id
		Ips        []string `json:"ips"`   //ip
		IsDisabled bool     `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * security domain
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	CorsDomain struct {
		Domains    []string `json:"domains"`
		IsAllowAll bool     `json:"is_allow_all"`
		IsDisabled bool     `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * third-party login option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	OauthOption struct {
		Wechat OauthPlatform
		Qq     OauthPlatform
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * third-party login platform
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	OauthPlatform struct {
		Ios     OauthApp
		Android OauthApp
		Pc      OauthApp
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * third-party login app
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	OauthApp struct {
		AppId            string `json:"app_id"`
		AppSecret        string `json:"app_secret"`
		CallbackUrl      string `json:"callback_url"`
		AuthorizeCodeUrl string `json:"authorize_code_url"`
		AccessTokenUrl   string `json:"access_token_url"`
		RefreshTokenUrl  string `json:"refresh_token_url"`
		OpenIdUrl        string `json:"open_id_url"`
		UserInfoUrl      string `json:"user_info_url"`
		IsDisabled       bool   `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * payment option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	PayOption struct {
		AliPay    AliPay
		WechatPay WechatPay
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * alipay payment
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	AliPay struct {
		AppId         string    `json:"app_id"`
		SellerId      string    `json:"seller_id"`
		PublicKey     string    `json:"public_key"`
		AppPublicKey  string    `json:"app_public_key"`
		AppPrivateKey string    `json:"app_private_key"`
		AppPrivatePem string    `json:"app_private_pem"`
		GatewayUrl    string    `json:"gateway_url"`
		ReturnUrl     string    `json:"return_url"`
		NotifyUrl     string    `json:"notify_url"`
		Api           AliPayApi `json:"api"`
		SignType      string    `json:"sign_type"`
		Format        string    `json:"format"`
		IsDisabled    bool      `json:"is_disabled"`
	}

	AliPayApi struct {
		AppPay           string `json:"app_pay"`
		PreCreate        string `json:"pre_create"`
		TradeQuery       string `json:"trade_query"`
		TradeRefund      string `json:"trade_refund"`
		TradeRefundQuery string `json:"trade_refund_query"`
		TradeClose       string `json:"trade_close"`
		BillQuery        string `json:"bill_query"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * wechat payment
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	WechatPay struct {
		AppId      string       `json:"app_id"`
		PartnerId  string       `json:"partner_id"`
		AppSecret  string       `json:"app_secret"`
		ApiSecret  string       `json:"api_secret"`
		NotifyUrl  string       `json:"notify_url"`
		Api        WechatPayApi `json:"api"`
		FeeType    string       `json:"fee_type"`
		SignType   string       `json:"sign_type"`
		IsDisabled bool         `json:"is_disabled"`
	}

	WechatPayApi struct {
		UnifiedOrder string `json:"unified_order"`
		OrderQuery   string `json:"order_query"`
		Refund       string `json:"refund"`
		RefundQuery  string `json:"refund_query"`
		CloseOrder   string `json:"close_order"`
		DownloadBill string `json:"download_bill"`
		Report       string `json:"report"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * sms option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SmsOption struct {
		Aliyun  AliyunSms
		Alidayu AlidayuSms
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * aliyun sms
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	AliyunSms struct {
		AppKey      string      `json:"app_key"`
		AppSecret   string      `json:"app_secret"`
		Gateway     string      `json:"gateway"`
		Api         SmsApi      `json:"api"`
		Template    SmsTemplate `json:"template"`
		RegionId    string      `json:"region_id"`
		Format      string      `json:"format"`
		Type        string      `json:"type"`
		SignName    string      `json:"sign_name"`
		SignMethod  string      `json:"sign_method"`
		Version     string      `json:"version"`
		MobileCount int         `json:"mobile_count"`
		CodeLength  int         `json:"code_length"`
		IsSsl       bool        `json:"is_ssl"`
		IsDisabled  bool        `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * alidayu sms
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	AlidayuSms struct {
		AppKey      string      `json:"app_key"`
		AppSecret   string      `json:"app_secret"`
		Gateway     string      `json:"gateway"`
		Api         SmsApi      `json:"api"`
		Template    SmsTemplate `json:"template"`
		Format      string      `json:"format"`
		Type        string      `json:"type"`
		SignName    string      `json:"sign_name"`
		SignMethod  string      `json:"sign_method"`
		Version     string      `json:"version"`
		MobileCount int         `json:"mobile_count"`
		CodeLength  int         `json:"code_length"`
		IsSsl       bool        `json:"is_ssl"`
		IsDisabled  bool        `json:"is_disabled"`
	}

	SmsApi struct {
		Send      string `json:"send"`
		SendQuery string `json:"send_query"`
	}

	SmsTemplate struct {
		Logon        string `json:"logon"`
		Register     string `json:"register"`
		Auth         string `json:"auth"`
		Info         string `json:"info"`
		FindPassword string `json:"find_password"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * mail option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	MailOption struct {
		Smtp MailServer
		Pop3 MailServer
	}

	MailServer struct {
		Host       string `json:"host"`
		Port       int32  `json:"port"`
		Username   string `json:"username"`
		Password   string `json:"password"`
		IsSsl      bool   `json:"is_ssl"`
		IsDisabled bool   `json:"is_disabled"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * storage option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	StorageOption struct {
		Static     StaticRoute
		Upload     UploadFile
		Server     StorageServer
		IsFdfs     bool
		IsDisabled bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * static route
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	StaticRoute struct {
		File string `json:"file"` //static file path identity
		Fdfs string `json:"fdfs"` //fdfs file path identity
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * upload file
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	UploadFile struct {
		Root  string         `json:"root"` //file upload root
		Temp  string         `json:"temp"` //file upload temporary directory
		Image UploadFileItem `json:"image"`
		Audio UploadFileItem `json:"audio"`
		Video UploadFileItem `json:"video"`
		File  UploadFileItem `json:"file"`
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * upload file option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	UploadFileItem struct {
		Formats []string `json:"formats"` //format collection
		Size    int32    `json:"size"`    //size bytes
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * fdfs server
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	StorageServer struct {
		Trackers []string `json:"trackers"` //Tracker Host collection
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * search option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SearchOption struct {
		Hosts               []string // Host collection
		DefaultAnalyzer     string   // Default word breaker name
		NumberOfShards      int      // Number of shards
		NumberOfReplicas    int      // Number of copies
		HealthcheckInterval int      // Health test interval, seconds
		IsDisabled          bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * log option
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	LogOption struct {
		Level      string //INFO | WARN | DEBUG | ERROR
		LogonCount int    //Number of logon
		IsDisabled bool
	}
)

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
