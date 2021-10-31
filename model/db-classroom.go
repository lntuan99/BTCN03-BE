package model

import "github.com/jinzhu/gorm"

type Classroom struct {
	gorm.Model
	Name          string
	CoverImageURL string
	Code          string `gorm:"uniqueIndex"`
	Description   string
	//Teacher []Teacher
	//Student []Student
}

type ClassroomRes struct {
	ID            uint
	Name          string `json:"name"`
	CoverImageURL string `json:"coverImageUrl"`
	Code          string `json:"code"`
	Description   string `json:"description"`
}

func (classroom Classroom) ToRes() ClassroomRes {
	return ClassroomRes{
		ID:            classroom.ID,
		Name:          classroom.Name,
		CoverImageURL: classroom.CoverImageURL,
		Code:          classroom.Code,
		Description:   classroom.Description,
	}
}

//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
//============================================================
func (classroom Classroom) FindClassroomByCode(code string) Classroom {
	var res Classroom
	DBInstance.First(&res, "code = ?", code)

	return res
}