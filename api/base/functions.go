package base

import (
    "busmap.vn/librarycore/util"
    "github.com/gin-gonic/gin"
    "strconv"
    "strings"
)

func ResponseErrorWithCode(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "status": 0,
        "code":   message,
    })
}

func ResponseError(c *gin.Context, messageCode string) {
    c.JSON(200, gin.H{
        "status": 0,
        "code":   messageCode,
    })
}

func ResponseErrorWithData(c *gin.Context, code string, payload interface{}) {
    c.JSON(200, gin.H{
        "status": 0,
        "code":   code,
        "data":   payload,
    })
}

func ResponseResult(c *gin.Context, payload interface{}) {
    if payload == nil {
        c.JSON(200, gin.H{
            "status": 1,
            "code":   CodeSuccess,
        })
    } else {
        c.JSON(200, gin.H{
            "status": 1,
            "code":   CodeSuccess,
            "data":   payload,
        })

    }
}

// -----------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------

func GetIntQuery(c *gin.Context, key string) int {
    s := c.Query(key)
    v, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return v
}

func GetInt64Query(c *gin.Context, key string) int64 {
    s := c.Query(key)
    v, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return v
}

func GetUIntQuery(c *gin.Context, key string) uint {
    s := c.Query(key)
    v, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return uint(v)
}

func GetUInt64Query(c *gin.Context, key string) uint64 {
    s := c.Query(key)
    v, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return uint64(v)
}

func GetIntParam(c *gin.Context, key string) int {
    s := c.Param(key)
    v, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return v
}

func GetInt64Param(c *gin.Context, key string) int64 {
    s := c.Param(key)
    v, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return v
}

func GetUIntParam(c *gin.Context, key string) uint {
    s := c.Param(key)
    v, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return uint(v)
}

func GetUInt64Param(c *gin.Context, key string) uint64 {
    s := c.Param(key)
    v, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0
    }
    return uint64(v)
}

func GetBoolQuery(c *gin.Context, key string) bool {
    s := c.Query(key)
    if s == "0" || s == "false" {
        return false
    } else if s == "1" || s == "true" {
        return true
    } else {
        return false
    }
}

// -----------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------
// -----------------------------------------------------------------------------------------

func MiddlewareClientVersion() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientVersion := c.GetHeader("client-version")
        arr := strings.Split(clientVersion, ":")
        os := "unknown"
        versionCode := 0

        if len(arr) >= 2 {
            os = arr[0]
            versionCode, _ = strconv.Atoi(arr[1])
        }
        c.Set("os", os)
        c.Set("versionCode", versionCode)

        language := c.GetHeader("language")
        if util.EmptyOrBlankString(language) {
            language = "vi"
        }
        c.Set("language", language)
    }
}
