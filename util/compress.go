package util

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	workingDir, _ = os.Getwd()
	mediaDir      = fmt.Sprintf("%s/public/system", workingDir)
	TemporaryDir  = mediaDir + "/temporary"
)

func ZipFiles(filename string, files []string, isRemoveTimestamp bool) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file, isRemoveTimestamp); err != nil {
			return err
		}
	}
	return nil
}

func AddFileToZip(zipWriter *zip.Writer, filename string, isRemoveTimestamp bool) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	if isRemoveTimestamp {
		fileNameTemp := header.Name
		parts := strings.Split(fileNameTemp, ".")
		var i = len(parts) - 2
		parts = append(parts[:i], parts[i+1:]...)
		header.Name = strings.Join(parts, ".")
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	//header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, []string, error) {
	var dirNames []string
	var fileNames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return dirNames, fileNames, err
	}
	defer r.Close()
	suffix := time.Now().Unix()

	for _, f := range r.File {
		fileNamePart := strings.Split(f.Name, "/")
		fileNamePart[0] = fmt.Sprintf("%s_%d", fileNamePart[0], suffix)
		f.Name = strings.Join(fileNamePart, "/")
		fmt.Println(f.Name)

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return dirNames, fileNames, fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			dirNames = append(dirNames, fpath)
			continue
		}

		fileNames = append(fileNames, fpath)
		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return dirNames, fileNames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return dirNames, fileNames, err
		}

		rc, err := f.Open()
		if err != nil {
			return dirNames, fileNames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return dirNames, fileNames, err
		}
	}
	return dirNames, fileNames, nil
}

func UnzipV2(src string) ([]string, []string, error) {
	return Unzip(src, TemporaryDir)
}