package methods

import (
    "busmap.vn/librarycore/api/base"
    req_res "busmap.vn/librarycore/api/req_res_struct"
    "busmap.vn/librarycore/config"
    "busmap.vn/librarycore/model"
    "busmap.vn/librarycore/util"
    "fmt"
    "github.com/gin-gonic/gin"
    "path/filepath"
    "time"
)

func MethodGetClassroomList(c *gin.Context) (bool, string, interface{}) {
    var classroomArray = make([]model.Classroom, 0)
    model.DBInstance.Find(&classroomArray)

    var classroomResArray = make([]model.ClassroomRes, 0)
    for _, classroom := range classroomArray {
        classroomResArray = append(classroomResArray, classroom.ToRes())
    }

    return true, base.CodeSuccess, classroomResArray
}

func MethodCreateClassroom(c *gin.Context) (bool, string, interface{}) {
    var classroomInfo req_res.PostCreateClassroom
    if err := c.ShouldBind(&classroomInfo); err != nil {
        return false, base.CodeBadRequest, nil
    }

    var newClassroom = model.Classroom{
        Name:          classroomInfo.Name,
        CoverImageURL: "",
        Code:          classroomInfo.Code,
        Description:   classroomInfo.Description,
    }

    existedClassroomCode := model.Classroom{}.FindClassroomByCode(newClassroom.Code)

    if existedClassroomCode.ID > 0 {
        return false, base.CodeExistedClassroomCode, nil
    }

    err := model.DBInstance.Create(&newClassroom).Error

    if err != nil {
        return false, base.CodeCreateClassroomFail, nil
    }

    coverImage, err := c.FormFile("coverImage")
    if err != nil {
        newFileName := fmt.Sprintf("%v_%v.%v",coverImage.Filename, time.Now().Unix(), filepath.Ext(coverImage.Filename))
        folderDst := fmt.Sprintf("%v/system/classrooms/%v", config.Config.MediaDir, newClassroom.ID)
        util.CreateFolderV2(folderDst)
        fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
        _ = c.SaveUploadedFile(coverImage, fileDst)
    }

    return true, base.CodeSuccess, nil
}