package models

import (
	"fmt"
	"github.com/pkg/sftp"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/astaxie/beego"
)

// File struct is for handling the file attribute
type File struct {
	FileName         string
	FileContent      string
	FilePath         string
	FileLastModified string
}

// Directory struct is for handing the  Directory attribute
type Directory struct {
	DirName         string
	DirPath         string
	DirLastModified string
}

// DirectoryList struct is for handing the Directory and its SubName attribute
type DirectoryList struct {
	ChildrenDirs  []Directory
	ChildrenFiles []File
}

type IndexValue struct {
	mux sync.Mutex
	wg sync.WaitGroup
	childrenDirTmp Directory
}

var (
	fileChan      chan File
	dirChan       chan Directory
	childrenFiles []File
	childrenDirs  []Directory
	file          File
	directory     Directory
	fileBuilder   strings.Builder
)

//init function to init the file read channel
func init(){
	fileChan = make(chan File)
	dirChan = make(chan Directory)
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

// SFTPFileRead function is for reading the file in the remote host
func SFTPFileRead(readFile File, sftpConn *sftp.Client) (file File, err error) {
	file = readFile
	sftpFile, err := sftpConn.Open(readFile.FilePath + "/" +readFile.FileName)
	if err != nil {
		beego.Error(err)
		return file, nil
	}
	defer sftpFile.Close()
	sftpFileByte, err := ioutil.ReadAll(sftpFile)
	if err != nil {
		beego.Error(err)
		return file, err
	}
	file.FileContent = string(sftpFileByte)
	return file, nil
}

// SFTPFileDirList function is for listing all files and directory in the remote host
func SFTPFileDirList(Path string, sftpConn *sftp.Client) (DirectoryList, error) {

	directoryList := DirectoryList{ChildrenDirs: nil, ChildrenFiles: nil}
	var filePath *string
	fileInfo, err := sftpConn.Stat(Path)
	if err != nil {
		beego.Error(err)
	} else if fileInfo.IsDir() {
		filePath = &Path
	} else {
		dirSlice := strings.Split(Path, "/")
		dirSlice = dirSlice[:len(dirSlice) - 1]
		Path = strings.Join(dirSlice, "/")
		filePath = &Path
	}
	walkFiles, err := sftpConn.ReadDir(*filePath)
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
			directory.DirLastModified = subFile.ModTime().Format("2006/Jan/02 15:04")
			directory.DirPath = *filePath
			childrenDirs = append(childrenDirs, directory)
		} else {
			file.FileName = subFile.Name()
			file.FileLastModified = subFile.ModTime().Format("2006/Jan/02 15:04")
			file.FilePath = *filePath
			childrenFiles = append(childrenFiles, file)
		}
	}
	//for aFile := range fileChan {
	//	childrenFiles = append(childrenFiles, aFile)
	//}
	//for bFile := range dirChan {
	//	childrenDirs = append(childrenDirs, bFile)
	//}
	directoryList = DirectoryList{ChildrenDirs: childrenDirs, ChildrenFiles: childrenFiles}
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

// FileWrite function is for writing the file in localhost
func FileWrite(filePath string, content string) (err error) {
	readFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 766)
	defer readFile.Close()
	if err != nil {
		beego.Error(err)
		return err
	} else if os.IsNotExist(err) {
		beego.Info("This File Not Existed")
		return err
	} else {
		if _, err := readFile.Write([]byte(content)); err != nil {
			beego.Error(err)
			return err
		}
		return nil
	}
}

// SFTPFileWrite function is for writing the file in the remote host
func SFTPFileWrite(file File, sftpConn *sftp.Client) (err error) {
	sftpFileWriter, err := sftpConn.OpenFile(file.FilePath+file.FileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE)
	defer sftpFileWriter.Close()
	if err != nil {
		beego.Error(err)
		return err
	} else if os.IsNotExist(err) {
		beego.Info("This File Not Existed")
		return err
	} else {
		if _, err := sftpFileWriter.Write([]byte(file.FileContent)); err != nil {
			beego.Error(err)
			return err
		}
		return nil
	}
}

// // FileCompare function is for comparing two different files and output difference
// func FileCompare(OldFile File, NewFile File) (diff string, err error) {

// }
