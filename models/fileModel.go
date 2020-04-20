package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"

	"github.com/astaxie/beego"
)

// File struct is for handling the file attribute
type File struct {
	FileName         string
	FileAccess       os.FileMode
	FileContent      string
	FilePath         string
	FileSize         int64
	FileLastModified time.Time
}

// Directory struct is for handing the  Directory attribute
type Directory struct {
	DirName         string
	DirAccess       os.FileMode
	DirSize         int64
	DirPath         string
	DirLastModified time.Time
}

// DirectoryList struct is for handing the Directory and its SubName attribute
type DirectoryList struct {
	ChildrenDirs  []Directory
	ChildrenFiles []File
}

var (
	fileChan      chan File
	dirChan       chan Directory
	childrenFiles []File
	childrenDirs  []Directory
	file          File
	directory     Directory
)

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

// SFTPFileDirList function is for listing all files and directory in the remote host
func SFTPFileDirList(Path string, sftpConn *sftp.Client) (DirectoryList, error) {

	directoryList := DirectoryList{ChildrenDirs: nil, ChildrenFiles: nil}
	_, err := sftpConn.Stat(Path)
	if err != nil {
		beego.Error()
	}
	walkFiles, err := sftpConn.ReadDir(Path)
	if err != nil {
		beego.Error(err)
		return directoryList, err
	}
	for _, subFile := range walkFiles {
		if subFile.Name()[0] == '.' {
			continue
		}
		if subFile.IsDir() {
			directory.DirName = subFile.Name()
			directory.DirLastModified = subFile.ModTime()
			directory.DirPath = Path
			directory.DirSize = subFile.Size()
			directory.DirAccess = subFile.Mode()
			dirChan <- directory
		} else {
			file.FileName = subFile.Name()
			file.FileLastModified = subFile.ModTime()
			file.FileSize = subFile.Size()
			file.FileAccess = subFile.Mode()
			file.FilePath = Path
			fileChan <- file
		}
	}

	return directoryList, nil
}

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
