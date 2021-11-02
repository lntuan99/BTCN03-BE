package routers

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/module/apmgin"
    "os"
    "web2/btcn/api/base"
    api_classroom "web2/btcn/api/routers/api-classroom"
    "web2/btcn/config"
)

func Initialize() *gin.Engine {
    r := gin.New()

    //Set GIN in RELEASE_MODE if BUSMAP_ENV is "production"
    if os.Getenv("ENV") == config.Production {
        gin.SetMode(gin.ReleaseMode)
    }

    corConfig := cors.DefaultConfig()
    corConfig.AllowAllOrigins = true
    corConfig.AllowHeaders = []string{
        "authorization", "Authorization",
        "content-type", "accept",
        "referer", "user-agent",
    }
    r.Use(cors.New(corConfig))

    r.Use(apmgin.Middleware(r))
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(base.MiddlewareClientVersion())

    routeVersion01 := r.Group("api/v1")

    // Multipart quota
    r.MaxMultipartMemory = 20971520 // Exactly 20MB
    r.Static("/media", "./public")

    classroomRoute := routeVersion01.Group("classroom")
    {
        classroomRoute.GET("/", api_classroom.HandlerGetClassroomList)
        classroomRoute.POST("/", api_classroom.HandlerCreateClassroom)
    }

    return r
}
