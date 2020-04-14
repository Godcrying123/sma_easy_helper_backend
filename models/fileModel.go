package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
)

// File struct is for handling the file attribute
type File struct {
	FileName         string
	FileType         string
	FileAccess       os.FileMode
	FileContent      string
	FilePath         string
	FileSize         int
	FileLastModified time.Time
}

// Directory struct is for handing the  Directory attribute
type Directory struct {
	DirName         string
	DirAccess       os.FileMode
	DirSize         int
	DirPath         string
	DirLastModified time.Time
}

// DirectoryList struct is for handing the Directory and its SubName attribute
type DirectoryList struct {
	ChildrenDirs  []Directory
	ChildrenFiles []File
}

// FileRead function is for reading the file in localhost
func FileRead(filePath string) (fileContent string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		beego.Error(err)
		return "", err
	}
	defer file.Close()
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		beego.Error(err)
		return "", err
	}
	// beego.Info(string(fileByte))
	return string(fileByte), nil
}

// // // SFTPFileRead function is for reading the file in the remote host
// func SFTPFileRead(filePath string, sftpConn *sftp.Client) (file File, err error) {
// }

// FileListDir function is for listing all file in a specific dirpath
func FileListDir(DirPath string, FilesChan chan<- string) chan<- string {
	for _, entry := range dirents(DirPath) {
		if entry.IsDir() {
			subdir := filepath.Join(DirPath, entry.Name())
			FileListDir(subdir, FilesChan)
		} else {
			FilesChan <- filepath.Join(DirPath, entry.Name())
		}
	}
	return FilesChan
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
	}
	return entries
}

// // SFTPFileListDir function is for listing all file in a specific dirpath from a remote host
// func SFTPFileListDir(DirPath string, sftpConn *sftp.Client) (File []File, wrongRead string, err error) {

// }

// // FileWrite function is for writing the file in localhost
// func FileWrite(filePath string, file File) (err error) {

// }

// // SFTPFileWrite function is for writing the file in the remote host
// func SFTPFileWrite(filePathn string, file File, sftpConn *sftp.Client) (err error) {

// }

// // FileCompare function is for comparing two different files and output difference
// func FileCompare(OldFile File, NewFile File) (diff string, err error) {

// }
