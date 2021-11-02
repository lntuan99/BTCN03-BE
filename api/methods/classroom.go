package methods

import (
    "web2/btcn/api/base"
    req_res "web2/btcn/api/req_res_struct"
    "web2/btcn/model"
    "web2/btcn/util"
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
    _ = c.Request.ParseMultipartForm(20971520)

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

    _, header, errFile := c.Request.FormFile("coverImage")
    if errFile == nil {
        newFileName := fmt.Sprintf("%v%v", time.Now().Unix(), filepath.Ext(header.Filename))
        folderDst := fmt.Sprintf("/system/classrooms/%v", newClassroom.ID)

        util.CreateFolder(folderDst)

        fileDst := fmt.Sprintf("%v/%v", folderDst, newFileName)
        if err := util.SaveUploadedFile(header, folderDst, newFileName); err == nil {
            model.DBInstance.
                Model(&newClassroom).
                Update("cover_image_url", fileDst)
        }
    }

    return true, base.CodeSuccess, nil
}