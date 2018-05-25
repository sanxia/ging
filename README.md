# ging
--------------
ging for gin micro web framework extending

auth: 美丽的地球啊 - mliu

date: 20161001


----- Example -----

+++++ main.go +++++

func main() {
    //解析参数
    appSettings, appRouter := parseApplication()

    if appSettings == nil || appRouter == nil {

        return

    }

    //引导器初始化

    Bootstrap(appSettings)

    //启动服务器

    serverOption := &ging.ServerOption{

        Host:  appSettings.Server.Host,

        Ports: appSettings.Server.Ports,

    }

    ging.Start(serverOption, appRouter)
}


func parseApplication() (*ging.Settings, ging.IHttpRouter) {

    //解析参数

    appNameFlag := flag.String("app", "", "请输入App名称")

    appHostFlag := flag.String("host", "", "请输入绑定Ip")

    flag.Parse()

    //App名称

    appName := *appNameFlag

    //App设置
    appSettingsList := map[string]*ging.Settings{

        "myapp": myapp.AppSettings,

    }

    appSettings, isOk := appSettingsList[appName]

    if !isOk {

        return nil, nil

    }

    //当前应用名称

    appSettings.AppName = appName


    //自定义绑定Ip优先
    appHost := *appHostFlag

    if len(appHost) > 0 {

        appSettings.Server.Host = appHost

    } else {

        return nil, nil

    }

    //App路由
    routers := map[string]ging.IHttpRouter{

        "myapp": myapp.NewHttpRouter(),

    }

    appRouter, isRouterOk := routers[appName]

    if !isRouterOk {

        log.Printf("parse application appName: %s isRouter error", appName)

        return nil, nil

    }

    return appSettings, appRouter

}


----- App -----

+++++ app.go +++++

package myapp

import (
    "github.com/gin-gonic/gin"

    "github.com/sanxia/ging"

    "github.com/sanxia/ging/middleware/authentication/cookie"

    "github.com/sanxia/ging/middleware/session"

)

// 创建路由

func NewHttpRouter() ging.IHttpRouter {

    return &HttpRouter{}

}

// Route 路由注册

func (r *HttpRouter) Route() *gin.Engine {

    httpEngine := ging.NewHttpEngine(AppSettings.Storage.HtmlTemplate.Path, ging.RELEASE)

    //静态路由

    httpEngine.Static("/static", "./assets")

    //认证中间件

    httpEngine.Middleware(cookie.CookieAuthenticationMiddleware(cookie.CookieExtend{
        Option: &cookie.CookieOption{
            Name:     AppSettings.Forms.Authentication.Cookie.Name,
            Path:     AppSettings.Forms.Authentication.Cookie.Path,
            Domain:   AppSettings.Forms.Authentication.Cookie.Domain,
            MaxAge:   AppSettings.Forms.Authentication.Cookie.MaxAge,
            HttpOnly: AppSettings.Forms.Authentication.Cookie.HttpOnly,
            Secure:   AppSettings.Forms.Authentication.Cookie.Secure,
        },
        LogonUrl:   AppSettings.Forms.Authentication.LogonUrl,
        DefaultUrl: AppSettings.Forms.Authentication.DefaultUrl,
        PassUrls:   AppSettings.Forms.Authentication.PassUrls,
        EncryptKey: AppSettings.Security.EncryptKey,
        IsJson:     true,
        IsEnabled:  true,
    }))

    //会话中间件
    httpEngine.Middleware(session.SessionMiddleware(&session.SessionOption{
        Cookie: &session.SessionCookieOption{
            Name:     AppSettings.Session.Cookie.Name,
            Path:     AppSettings.Session.Cookie.Path,
            Domain:   AppSettings.Session.Cookie.Domain,
            MaxAge:   AppSettings.Session.Cookie.MaxAge,
            HttpOnly: AppSettings.Session.Cookie.HttpOnly,
            Secure:   AppSettings.Session.Cookie.Secure,
        },
        Redis: &session.SessionRedisOption{
            Host:     AppSettings.Session.RedisStore.Host,
            Port:     AppSettings.Session.RedisStore.Port,
            Password: AppSettings.Session.RedisStore.Password,
            Prefix:   AppSettings.Redis.Prefix,
        },
        EncryptKey: AppSettings.Security.EncryptKey,
        StoreType:  AppSettings.Session.StoreType,
    }))

    //passport router
    passportGroup := httpEngine.Group("/passport")
    passportController := passport.NewController()
    passportController.Filter(filter.PassportFilter())
    {
        passportGroup.POST("/validatecode", passportController.Action(passportController.ValidateCode))
    }

    //public router
    publicGroup := httpEngine.Group("/public")
    publicController := public.NewController()
    {
        //test
        publicGroup.GET("/user/test", publicController.Action(publicController.TestUser))
    }

    return httpEngine.Engine()
}
