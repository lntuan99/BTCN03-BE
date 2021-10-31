package routers

import (
    "busmap.vn/librarycore/api/base"
    api_classroom "busmap.vn/librarycore/api/routers/api-classroom"
    "busmap.vn/librarycore/config"
    "fmt"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "go.elastic.co/apm/module/apmgin"
    "os"
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
    corConfig.AllowCredentials = true
    r.Use(cors.Default())

    fmt.Println(corConfig)

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
