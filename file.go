package ging

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/sanxia/gfs"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * 文件数据
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IFileStorage interface {
		IFilePath
		IFileOperation
	}

	IFilePath interface {
		UrlToPath(url string, args ...string) string
		PathToUrl(path string, args ...string) string
	}

	IFileOperation interface {
		Upload(data []byte, fileExtName string, args ...bool) (*File, error)
		Delete(filename string, args ...bool) error
	}

	IDiskFileOperation interface {
		UploadFileToDisk(data []byte, fileExtName string) (*File, error)
		DeleteDiskFile(filename string) error
	}

	IFdfsFileOperation interface {
		GetFileByFileId(fileId string) (*File, error)
		UploadFileToFdfs(data []byte, fileExtName string) (*File, error)
		DeleteFileByFileId(fileId string) error
	}
)

type (
	fileStorage struct {
		client *gfs.FdfsClient
	}

	File struct {
		Id   string `form:"id" json:"id"`
		Path string `form:"path" json:"path"`
		Url  string `form:"url" json:"-"`
		Data []byte `form:"data" json:"-"`
		Size int64  `form:"size" json:"-"`
	}
)

/* ================================================================================
 * 初始化FileStorage
 * author: smalltour - 老牛|共生美好
 * ================================================================================ */
func NewFileStorage() IFileStorage {
	return &fileStorage{}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 上传文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) Upload(data []byte, fileExtName string, args ...bool) (*File, error) {
	isToDisk := true
	if len(args) > 0 {
		isToDisk = args[0]
	}

	if isToDisk {
		return s.UploadFileToDisk(data, fileExtName)
	}

	return s.UploadFileToFdfs(data, fileExtName)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) Delete(filename string, args ...bool) error {
	isToDisk := true
	if len(args) > 0 {
		isToDisk = args[0]
	}

	if isToDisk {
		return s.DeleteDiskFile(filename)
	}

	return s.DeleteFileByFileId(filename)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 上传文件到磁盘
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) UploadFileToDisk(data []byte, fileExtName string) (*File, error) {
	filename := glib.Guid() + "." + fileExtName

	//保存文件
	fullFilename, err := glib.SaveFile(data, filename, GetApp().GetSetting().Storage.Upload.Root)
	if err != nil {
		return nil, err
	}

	//判断文件资源码
	resourceCode := s.FileExtNameToResourceCode(fileExtName)

	relativePath := glib.GetRelativePath(fullFilename)

	file := new(File)
	file.Path = relativePath
	file.Url = s.PathToUrl(relativePath, resourceCode)

	return file, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除磁盘文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) DeleteDiskFile(filename string) error {
	if filename == "" {
		return errors.New("argument error")
	}

	if strings.HasPrefix(filename, "http://") {
		filename = s.UrlToPath(filename)
	}

	fileRoot := ""
	if !strings.HasPrefix(filename, GetApp().GetSetting().Storage.Upload.Root) {
		fileRoot = GetApp().GetSetting().Storage.Upload.Root
	}

	err := glib.DeleteFile(filename, fileRoot)

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 根据文件id获取文件数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) GetFileByFileId(fileId string) (*File, error) {
	s.initFdfsClient()

	if s.client == nil {
		return nil, errors.New("fdfs error")
	}

	downloadResponse, err := s.client.DownloadToBuffer(fileId, 0, 0)
	if err != nil {
		return nil, errors.New("fdfs download file error")
	}

	file := new(File)
	file.Path = downloadResponse.FileId
	if data, isOk := downloadResponse.Content.([]byte); isOk {
		file.Data = data
	}

	file.Size = downloadResponse.Size

	return file, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 上传文件到Fdfs
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) UploadFileToFdfs(data []byte, fileExtName string) (*File, error) {
	s.initFdfsClient()

	if s.client == nil {
		return nil, errors.New("fdfs error")
	}

	uploadResponse, err := s.client.UploadByBuffer(data, fileExtName)
	if err != nil {
		return nil, errors.New("fdfs upload file error")
	}

	//判断文件资源码
	resourceCode := s.FileExtNameToResourceCode(fileExtName)

	file := new(File)
	file.Path = uploadResponse.FileId
	file.Url = s.PathToUrl(uploadResponse.FileId, resourceCode)

	return file, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除Fdfs文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) DeleteFileByFileId(fileId string) error {
	s.initFdfsClient()

	if s.client == nil {
		return errors.New("fdfs error")
	}

	if strings.HasPrefix(fileId, "http://") {
		fileId = s.UrlToPath(fileId)
	}

	err := s.client.DeleteFile(fileId)
	if err != nil {
		return errors.New("fdfs delete file error")
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 根据文件物理路径和资源编码获取文件Url地址
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) PathToUrl(path string, args ...string) string {
	domain, url := "", GetApp().GetSetting().Image.Default
	resourceCode := "image"
	if len(args) == 1 {
		resourceCode = strings.ToLower(args[0])
	}

	hosts := make(map[string]string, 0)
	hosts["image"] = GetApp().GetSetting().Domain.ImageHost
	hosts["audio"] = GetApp().GetSetting().Domain.AudioHost
	hosts["video"] = GetApp().GetSetting().Domain.VideoHost
	hosts["file"] = GetApp().GetSetting().Domain.FileHost

	if host, isOk := hosts[resourceCode]; isOk {
		domain = host
	}

	if GetApp().GetSetting().Domain.IsSsl {
		if !strings.HasPrefix(domain, "https") {
			domain = fmt.Sprintf("%s://%s", "https", domain)
		}
	} else {
		if !strings.HasPrefix(domain, "http") {
			domain = fmt.Sprintf("%s://%s", "http", domain)
		}
	}

	if len(path) > 0 {
		if strings.HasPrefix(path, "group") {
			path = s.FileIdToFilePath(path)
		} else {
			staticPaths := strings.Split(GetApp().GetSetting().Storage.Upload.Root, string(os.PathSeparator))
			if !strings.HasPrefix(path, GetApp().GetSetting().Storage.Upload.Root) {
				path = filepath.Join(GetApp().GetSetting().Storage.Upload.Root, path)
			}

			//文件路径替换成url路由
			path = strings.Replace(path, staticPaths[0], GetApp().GetSetting().Storage.Static.File, -1)
		}
		url = fmt.Sprintf("%s/%s", domain, path)
	}

	return url
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 根据Url和资源编码获取文件相对物理路径
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) UrlToPath(url string, args ...string) string {
	domain, path := "", ""
	resourceCode := "image"
	if len(args) == 1 {
		resourceCode = strings.ToLower(args[0])
	}

	hosts := make(map[string]string, 0)
	hosts["image"] = GetApp().GetSetting().Domain.ImageHost
	hosts["audio"] = GetApp().GetSetting().Domain.AudioHost
	hosts["video"] = GetApp().GetSetting().Domain.VideoHost
	hosts["file"] = GetApp().GetSetting().Domain.FileHost

	if host, isOk := hosts[resourceCode]; isOk {
		domain = host
	}

	if len(url) > 0 {
		url = glib.FilterHostProtocol(url)
		url = strings.Replace(url, domain, "", -1)

		if strings.HasPrefix(url, string(os.PathSeparator)) {
			url = url[1:]
		}

		fdfsRouter := GetApp().GetSetting().Storage.Static.Fdfs + string(os.PathSeparator)
		if strings.HasPrefix(url, fdfsRouter) {
			path = s.FilePathToFileId(url)
		} else {
			//去除url路由和上传根目录名: static/upload
			//只存储目录和文件名：2017/09/21/c2bac40d2759fa536f2debdd264d656a.png
			paths := strings.Split(url, string(os.PathSeparator))[2:]
			path = strings.Join(paths, string(os.PathSeparator))
		}
	}

	return path
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 文件路径转换成Fdfs文件id
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) FilePathToFileId(filePath string) string {
	filePaths := strings.Split(filePath, string(os.PathSeparator))
	dir := filePaths[2]

	group := "group" + filePaths[1]
	store := "M" + dir[:2]
	dir1 := dir[2:4]
	dir2 := dir[4:]
	filename := filePaths[3]

	fileIds := []string{group, store, dir1, dir2, filename}

	return strings.Join(fileIds, string(os.PathSeparator))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Fdfs文件id转换成文件路径
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) FileIdToFilePath(fileId string, args ...string) string {
	prefix := "s"
	if len(args) > 0 {
		prefix = args[0]
	}

	fileIds := strings.Split(fileId, string(os.PathSeparator))
	group := strings.Replace(fileIds[0], "group", "", -1)
	store := strings.Replace(fileIds[1], "M", "", -1)
	dir1 := fileIds[2]
	dir2 := fileIds[3]
	filename := fileIds[4]

	dir := store + dir1 + dir2
	filePaths := []string{prefix, group, dir, filename}

	return strings.Join(filePaths, string(os.PathSeparator))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 文件扩展名转换成资源码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) FileExtNameToResourceCode(fileExtName string) string {
	resourceCode := ""
	resourceTypeMap := map[string]string{
		"jpg":  "image",
		"jpeg": "image",
		"png":  "image",
		"bmp":  "image",
		"gif":  "image",
		"mp3":  "audio",
		"mp4":  "video",
	}

	if resourceType, isOk := resourceTypeMap[fileExtName]; isOk {
		resourceCode = resourceType
	}

	return resourceCode
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 连接fdfs服务
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *fileStorage) initFdfsClient() {
	if s.client == nil {
		if client, err := gfs.NewFdfsClient(GetApp().GetSetting().Storage.Server.Trackers); err != nil {
			fmt.Printf("fdfs client error: %v\n", err)
		} else {
			s.client = client
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取int64值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *File) IdInt64() int64 {
	return glib.StringToInt64(s.Id)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置int64值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *File) SetIdInt64(id int64) {
	s.Id = glib.Int64ToString(id)
}
