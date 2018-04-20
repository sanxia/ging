package ging

/* ================================================================================
 * 设置数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Settings struct {
		AppName      string             //应用名称
		Server       ServerOption       //服务器
		Domain       DomainOption       //域名
		Image        ImageOption        //图像
		Font         FontOption         //字体
		Forms        FormsOption        //表单
		Session      SessionOption      //会话
		Redis        RedisOption        //Redis
		MessageQueue MessageQueueOption //消息队列
		Database     DatabaseOption     //数据库
		Security     SecurityOption     //安全
		Pay          PayOption          //支付
		ValidateCode ValidateCodeOption //验证码
		Mail         MailOption         //邮件
		Sms          SmsOption          //短信
		AliyunSms    SmsOption          //阿里云短信
		Im           ImOption           //即时通信
		Storage      StorageOption      //存储
		Cors         CorsOption         //跨域
		Test         TestOption         //测试
		Log          LogOption          //日志
		Version      VersionOption      //版本
		IsCache      bool               //是否启用缓存
		IsHttps      bool               //是否https
		IsTesting    bool               //是否测试
		IsOnline     bool               //是否线上
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 服务器设置
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	ServerOption struct {
		Host  string
		Ports []int
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 域名设置
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DomainOption struct {
		Res   string
		Image string
		Audio string
		Video string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 图片配置结构
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	ImageOption struct {
		Avatar  AvatarOption
		Default DefaultImageOption
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 字体配置结构
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FontOption struct {
		Path string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 默认用户头像(未知, 男, 女)
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	AvatarOption struct {
		Default string
		Male    string
		Female  string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 默认图片
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DefaultImageOption struct {
		Image string
		Audio string
		Video string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 表单选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FormsOption struct {
		Register       FormsRegisterOption
		Authentication FormsAuthenticationOption
		Identity       string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 表单注册
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FormsRegisterOption struct {
		UsernameRule       string   //用户名正则表达式
		PasswordRule       string   //密码正则表达式
		ForbiddenUsers     []string //禁止注册的用户名集合
		ForbiddenPasswords []string //禁止使用的密码集合
		IsApproved         bool     //是否需要审核
		IsActived          bool     //是否需要激活
		IsEnabled          bool     //是否开启新用户注册
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 表单认证
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FormsAuthenticationOption struct {
		LogonUrl        string             //登陆认证Url
		DefaultUrl      string             //默认返回地址
		PassUrls        []string           //无需认证Url集合
		Cookie          *FormsCookieOption //form cookie
		LoginErrorCount int32              //最大连续错误次数
		LockMinutes     int32              //锁定多少分钟
		IsPersistence   bool               //是否持久会话
		IsEnabled       bool               //是否开启登入
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 表单Cookie
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FormsCookieOption struct {
		Name     string //客户端Cookie名称
		Path     string //客户端Cookie路径
		Domain   string //客户端Cookie域
		MaxAge   int    //过期时长（单位：秒）
		HttpOnly bool   //是否只能http读取
		Secure   bool   //是否https
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 会话选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SessionOption struct {
		Cookie     *FormsCookieOption //form cookie
		RedisStore RedisStoreOption   //Redis存储会话
		StoreType  string             //会话存储类型码
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * Redis选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	RedisOption struct {
		Prefix string //key前缀
		Stores []RedisStoreOption
		Expire int //多少秒后过期
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * Redis存储
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	RedisStoreOption struct {
		Host     string
		Port     int
		Password string
		Timeout  int
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * MessageQueue选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	MessageQueueOption struct {
		Server MessageQueueServerOption
		Vhost  string //主机
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 消息服务器选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	MessageQueueServerOption struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 数据库选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	DatabaseOption struct {
		Connections []*DatabaseConnectionOption
		IsLog       bool
	}

	DatabaseConnectionOption struct {
		Key      string
		Username string
		Password string
		Host     string
		Database string
		Dialect  string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 安全选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SecurityOption struct {
		ApiSecret            string //Api密匙
		AccountSecret        string //账户密匙
		EncryptKey           string //安全加密Key
		VerifyWaitingMinutes int64  //验证等待多少分钟
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 支付设置
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	PayOption struct {
		WeChatPay WeChatPayOption
		AliPay    AlipayPayOption
	}

	WeChatPayOption struct {
		AppId     string
		AppSecret string
		PartnerId string
		ApiSecret string
		NotifyUrl string
		Apis      map[string]string
		FeeType   string
		SingType  string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 支付宝选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	AlipayPayOption struct {
		AppId         string
		SellerId      string
		GatewayUrl    string
		ReturnUrl     string
		NotifyUrl     string
		PublicKey     string
		AppPublicKey  string
		AppPrivateKey string
		AppPrivatePem string
		Apis          map[string]string
		SingType      string
		Format        string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 验证码选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	ValidateCodeOption struct {
		Image                string //图片验证码名称
		Mobile               string //手机验证码名称
		VerifyWaitingMinutes int64  //动态验证码最大过期分钟数
		ErrorCount           int    //错误次数，需要输入验证码
		IsRegister           bool   //注册是否需要验证码
		IsLogon              bool   //登陆是否需要验证码
		IsComment            bool   //评论是否需要验证码
		IsEnabled            bool   //是否启用验证码
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 短信选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	SmsOption struct {
		AppKey         string
		AppSecret      string
		Host           SmsHostOption
		Method         SmsMethodOption
		Template       SmsTemplateOption
		MaxMobileCount int
		MaxCodeLength  int
		RegionId       string
		Format         string
		Type           string
		SignName       string
		SignMethod     string
		Version        string
		IsEnabled      bool
	}

	SmsHostOption struct {
		Http  string
		Https string
	}

	SmsMethodOption struct {
		Send      string
		SendQuery string
	}

	SmsTemplateOption struct {
		Code    SmsCodeOption
		Product string
	}

	SmsCodeOption struct {
		Register     string
		Logon        string
		Auth         string
		Info         string
		FindPassword string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 电子邮件选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	MailOption struct {
		Host     string //Smtp主机
		Port     int32  //Smtp端口
		Username string //发送邮件用户明
		Password string //发送邮件密码
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 即时通信选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	ImOption struct {
		OrgName      string //机构名
		AppName      string //应用名
		ClientId     string //client id
		ClientSecret string //client secret
		Host         string //主机地址
		IsEnabled    bool   //是否启用
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 文件存储选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	StorageOption struct {
		HtmlTemplate HtmlTemplateOption
		Static       StaticOption
		Upload       UploadOption
		Fdfs         FdfsOption
		IsFdfs       bool
	}

	HtmlTemplateOption struct {
		Path       string
		Extensions []string
	}

	StaticOption struct {
		File string
		Fdfs string
	}

	UploadOption struct {
		Root      string           //上传文件存储目录
		Temp      string           //上传文件存储临时目录
		Size      FileSizeOption   //文件大小
		Format    FileFormatOption //文件格式
		IsEnabled bool             //是否启用上传
	}

	FdfsOption struct {
		TrackerHosts []string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 文件大小
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FileSizeOption struct {
		Image int64 //图片文件（字节）
		Audio int64 //音频文件（字节）
		Video int64 //视频文件（字节）
		File  int64 //其它文件（字节）
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 文件格式
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	FileFormatOption struct {
		Images []string
		Audios []string
		Videos []string
		Files  []string
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 跨域选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	CorsOption struct {
		Domains          []string
		IsAllowAllDomain bool
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 测试选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	TestOption struct {
		Mobiles    []string //测试手机号码列表
		MobileCode string   //手机验证码
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 日志选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	LogOption struct {
		Level     string //级别（INFO | WARN | DEBUG | ERROR）
		IsEnabled bool   //是否开启日志
	}

	/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
	 * 版本选项
	 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
	VersionOption struct {
		Value       float64           //版本值
		Description string            //描述信息
		Urls        map[string]string //下载地址
		IsRequired  bool              //是否必须更新
	}
)
