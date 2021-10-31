package util

import (
    "busmap.vn/librarycore/config"
    "busmap.vn/librarycore/config/constants"
    "errors"
    "fmt"
    "mime/multipart"
    "path/filepath"
    "time"
)

func GetFolderPathForUserFile(userID uint) string {
    if userID == 0 {
        return ""
    }
    folderPath := constants.UPLOAD_FILE_FROM_APP
    folderPath = fmt.Sprintf("%v/%v", folderPath, userID)
    return folderPath
}

func UploadUserFile(fileHeader *multipart.FileHeader, userID uint) (string, string, string, error) {
    // Evaluate parameters
    if userID == 0 {
        messageError := "[FILE_UPDATE_REMOVE][LIBRARY] Invalid user-id\n"
        fmt.Print(messageError)
        return "", "", "", errors.New(messageError)
    }

    // Create folder if not exists
    folderPath := GetFolderPathForUserFile(userID)

    // Upload file
    return uploadFile(fileHeader, folderPath)
}

func RemoveUserFile(fileName string, userID uint) (bool, error) {
    if userID == 0 {
        messageError := "[FILE_UPDATE_REMOVE][LIBRARY] Invalid user-id\n"
        fmt.Print(messageError)
        return false, errors.New(messageError)
    }

    folderPath := GetFolderPathForUserFile(userID)
    RemoveFile(folderPath, fileName)
    return true, nil
}

func uploadFile(
    fileHeader *multipart.FileHeader,
    folderPath string,
) (string, string, string, error) {
    // Evaluate file-header
    if fileHeader == nil {
        messageError := "[FILE_UPDATE_REMOVE][LIBRARY] Multipart-file-header mustn't be empty\n"
        fmt.Print(messageError)
        return "", "", "", errors.New(messageError)
    }

    // Create folder
    CreateFolder(folderPath)

    // Process hash-file-name
    fileName := fileHeader.Filename
    fileFullName := filepath.Base(fileName)
    fileExtension := filepath.Ext(fileName)
    fileNameForHash := fmt.Sprintf("%s_%d", fileFullName, time.Now().Unix())
    hashFileName := fmt.Sprintf("%s%s", HexSha256String([]byte(fileNameForHash)), fileExtension)

    // Save image file
    if err := SaveUploadedFile(fileHeader, folderPath, hashFileName); err != nil {
        errorMessage := fmt.Sprintf("[FILE_UPDATE_REMOVE][LIBRARY] Error has occurred: %v\n", err.Error())
        fmt.Print(errorMessage)
        return "", "", "", errors.New(errorMessage)
    }

    // ShortFileUrl for storing in tables of database.
    // FullFireUrl for viewing on app or web.
    apiDomain := config.Config.Domain
    shortFileUrl := fmt.Sprintf("%s/%s", folderPath, hashFileName)
    fullFileUrl := fmt.Sprintf("%s/media%s", apiDomain, shortFileUrl)
    return hashFileName, shortFileUrl, fullFileUrl, nil
}
