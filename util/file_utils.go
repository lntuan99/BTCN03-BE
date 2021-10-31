package util

import (
    "busmap.vn/librarycore/config"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"
)

func GetPermissionFile() os.FileMode {
    return 0755
}

func GetWorkingDir() string {
    // Ex: C:\Users\kienquoc\go\src\busmap.vn\librarycore
    workingDir, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    return workingDir
}

func CreateFolder(folderPath string) {
    // Example:
    // + WorkingDir = C:\Users\kienquoc\go\src\busmap.vn\librarycore
    // + ParentFolderName = \public
    // + FolderPath = \system\upload\app\1\ (folder name is user-id)
    parentFolderName := config.Config.MediaDir
    workingDir := GetWorkingDir()
    fullPath := fmt.Sprintf("%s%s%s", workingDir, parentFolderName, folderPath)
    if _, err := os.Stat(fullPath); err != nil && os.IsNotExist(err) {
        err := os.MkdirAll(fullPath, GetPermissionFile())
        if err != nil {
            panic(err)
        }
    }
}

func CreateFolderV2(folderPath string) {
    //Create Folder if not exists
    if _, err := os.Stat(folderPath); err != nil && os.IsNotExist(err) {
        err := os.MkdirAll(folderPath, GetPermissionFile())
        if err != nil {
            panic(err)
        }
    }
}

func RemoveFolder(folderPath string) {
    // Example:
    // + WorkingDir = C:\Users\kienquoc\go\src\busmap.vn\librarycore
    // + ParentFolderName = \public
    // + FolderPath = \system\upload\app\1\ (folder name is user-id)
    parentFolderName := config.Config.MediaDir
    workingDir := GetWorkingDir()
    fullPath := fmt.Sprintf("%s%s%s", workingDir, parentFolderName, folderPath)
    _ = os.RemoveAll(fullPath)
}

func RemoveFile(folderPath, fileName string) {
    // Example:
    // + ParentFolderName = /public
    // + FolderPath = /system/upload/app/1/ (folder name is user-id)
    // + FileName = hello_world.png
    parentFolderName := config.Config.MediaDir
    dst := fmt.Sprintf(".%s%s/%s", parentFolderName, folderPath, fileName)
    if _, err := os.Stat(dst); err != nil && os.IsNotExist(err) {
        return
    } else {
        _ = os.Remove(dst)
    }
}

func RemoveFileV2(filePath string) {
    // Example:
    // + ParentFolderName = /public
    // + FilePath = /system/upload/app/1/hello_world.png
    parentFolderName := config.Config.MediaDir
    dst := fmt.Sprintf(".%s%s", parentFolderName, filePath)
    if _, err := os.Stat(dst); err != nil && os.IsNotExist(err) {
        return
    } else {
        _ = os.Remove(dst)
    }
}

func SaveUploadedFile(file *multipart.FileHeader, folderPath, fileName string) error {
    // Example:
    // + ParentFolderName = /public
    // + FolderPath = /system/upload/app/1/ (folder name is user-id)
    // + FileName = hello_world.png
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    parentFolderName := config.Config.MediaDir
    dst := fmt.Sprintf(".%s%s/%s", parentFolderName, folderPath, fileName)
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}

func CopyFile(src string, dest string) error {
    fileInput, err := ioutil.ReadFile(src)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(dest, fileInput, 0644)
    if err != nil {
      return err
    }

    return nil
}

func CopyFileV2(src string, dest string) string {
    //This function make sure file name copy is unique
    unix := time.Now().Unix()
    fileFullName := filepath.Base(src)
    fileNamePart := strings.Split(fileFullName, ".")
    fileNameResult := fmt.Sprintf("%s_%d.%s", fileNamePart[0], unix, fileNamePart[1])

    dest = fmt.Sprintf("%s/%s", dest, fileNameResult) // for name is always unique

    if err := CopyFile(src, dest); err != nil {
        return ""
    }

    return fileNameResult
}

func DownloadFile(URL, fileName string) error {
    //Get the response bytes from the url
    response, err := http.Get(URL)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        return errors.New("Received non 200 response code")
    }
    //Create a empty file
    file, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer file.Close()

    //Write the bytes to the file
    _, err = io.Copy(file, response.Body)
    if err != nil {
        return err
    }

    return nil
}
