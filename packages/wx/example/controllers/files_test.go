package controllers

import (
	"fmt"
	"mime/multipart"
	"testing"
	"wx"

	"github.com/stretchr/testify/assert"
)

type S3Service struct {
}

type FileService struct {
}

func (fSvc *FileService) New() error {
	fmt.Println("FileService New")
	return nil
}
func (fSvc *FileService) SaveFile() error {
	fmt.Println("FileService SaveFile")
	return nil
}

type Files struct {
	fileSvc FileService
}

func (fs *Files) New(fileSvc wx.Global[FileService]) error {
	fService, err := fileSvc.Ins() //Load instance of FileService
	if err != nil {
		return err
	}
	fmt.Println(fService)
	fs.fileSvc = *fService
	return nil
}

type UploadResult struct {
	Url string
}

func (fs *Files) SaveFileContent() error {
	fmt.Println("SaveFileContent, SaveFileContent is non handler method")
	err := fs.fileSvc.SaveFile()
	return err
}

func (fs *Files) Upload(ctx *wx.Handler, data struct {
	File multipart.File
}) (*UploadResult, error) {
	return &UploadResult{
		Url: "Emulator url",
	}, nil
}

func TestFileController(t *testing.T) {

}
func TestFilesUpload(t *testing.T) {
	mt := wx.GetMethodByName[Files]("Upload")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "files/upload", info.Uri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, "files/upload", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, 0, len(info.UriParams))
	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 1, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, 2, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)

	mt2 := wx.GetMethodByName[Files]("SaveFileContent")
	assert.NotEmpty(t, *mt2)
	info2, err := wx.Helper.GetHandlerInfo(*mt2)
	assert.Nil(t, err)
	assert.Empty(t, info2)

}
func (fs *Files) UploadWithUser(ctx *wx.Handler, data struct {
	File multipart.File
}, user wx.UserClaims) (*UploadResult, error) {
	return &UploadResult{
		Url: "Emulator url",
	}, nil
}
func TestUploadWithUser(t *testing.T) {
	mt := wx.GetMethodByName[Files]("UploadWithUser")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "files/upload-with-user", info.Uri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, "files/upload-with-user", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, 0, len(info.UriParams))
	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 1, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, 2, info.IndexOfRequestBody)
	assert.Equal(t, 1, len(info.IndexOfAuthClaims))
	assert.Equal(t, []int{0}, info.IndexOfAuthClaims)
	assert.Equal(t, 3, info.IndexOfAuthClaimsArg)

}
func (fs *Files) UploadMutiTenant(ctx *struct {
	*wx.Handler `route:"{Tenant}/upload"`
	Tenant      string
}, data struct {
	File multipart.File
}, user wx.UserClaims) (*UploadResult, error) {
	return &UploadResult{
		Url: "Emulator url",
	}, nil
}
func TestUploadMutiTenant(t *testing.T) {
	mt := wx.GetMethodByName[Files]("UploadMutiTenant")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "files/{Tenant}/upload", info.Uri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, "^files/([^/]+)/upload$", info.RegexUri)
	assert.Equal(t, "files/", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "POST", info.HttpMethod)
	assert.Equal(t, 1, len(info.UriParams))
	assert.Equal(t, "Tenant", info.UriParams[0].Name)
	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 1, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, 2, info.IndexOfRequestBody)
	assert.Equal(t, 1, len(info.IndexOfAuthClaims))
	assert.Equal(t, []int{0}, info.IndexOfAuthClaims)
	assert.Equal(t, 3, info.IndexOfAuthClaimsArg)

}
func (fs *Files) Download(ctx *struct {
	wx.Handler `route:"method:get"`
}) error {
	return nil
}
func TestDownload(t *testing.T) {
	mt := wx.GetMethodByName[Files]("Download")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "files/download", info.Uri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, "files\\/download", info.RegexUri)
	assert.Equal(t, "files/download", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, 0, len(info.UriParams))

	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 0, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, -1, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Nil(t, info.IndexOfAuthClaims)
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)
}
func (fs *Files) DownloadInUri(ctx *struct {
	wx.Handler `route:"{FileName};method:get"`
	FileName   string
}) error {
	return nil
}
func TestDownloadInUri(t *testing.T) {
	mt := wx.GetMethodByName[Files]("DownloadInUri")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "files/{FileName}", info.Uri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, "^files/([^/]+)$", info.RegexUri)
	assert.Equal(t, "files/", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, 1, len(info.UriParams))

	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 0, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, -1, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Nil(t, info.IndexOfAuthClaims)
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)
}
func (fs *Files) DownloadFromAbsUri(ctx *struct {
	wx.Handler `route:"/{FileName};method:get"`
	FileName   string
}) error {
	return nil
}
func TestDownloadFromAbsUri(t *testing.T) {
	mt := wx.GetMethodByName[Files]("DownloadFromAbsUri")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "/files/{FileName}", info.Uri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, "^files/([^/]+)$", info.RegexUri)
	assert.Equal(t, "/files/", info.UriHandler)
	assert.Equal(t, true, info.IsAbsUri)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, 1, len(info.UriParams))

	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 0, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, -1, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Nil(t, info.IndexOfAuthClaims)
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)
}
func (fs *Files) DownloadWithQuery(ctx *struct {
	wx.Handler `route:"@/?file={FileName};method:get"`
	FileName   string
}) error {
	return nil
}
func TestDownloadWithQuery(t *testing.T) {
	mt := wx.GetMethodByName[Files]("DownloadWithQuery")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "files/download-with-query?file={FileName}", info.Uri)
	assert.Equal(t, true, info.IsQueryUri)
	assert.Equal(t, false, info.IsRegexHandler)
	assert.Equal(t, "files\\/download\\-with\\-query", info.RegexUri)
	assert.Equal(t, "files/download-with-query", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, 0, len(info.UriParams))

	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 0, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, -1, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Nil(t, info.IndexOfAuthClaims)
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)
}
func (fs *Files) DownloadWithParamAndQuery(ctx *struct {
	wx.Handler `route:"{TenantName}/@/?file={FileName};method:get"`
	FileName   string
	TenantName string
}) error {
	return nil
}
func TestDownloadWithParamAndQuery(t *testing.T) {
	mt := wx.GetMethodByName[Files]("DownloadWithParamAndQuery")
	assert.NotEmpty(t, *mt)
	info, err := wx.Helper.GetHandlerInfo(*mt)
	assert.Nil(t, err)
	assert.NotEmpty(t, info)
	assert.Equal(t, "{TenantName}/files/download-with-param-and-query?file={FileName}", info.Uri)
	assert.Equal(t, true, info.IsQueryUri)
	assert.Equal(t, true, info.IsRegexHandler)
	assert.Equal(t, "^([^/]+)/files/download\\-with\\-param\\-and\\-query$", info.RegexUri)
	assert.Equal(t, "", info.UriHandler)
	assert.Equal(t, false, info.IsAbsUri)
	assert.Equal(t, "GET", info.HttpMethod)
	assert.Equal(t, 1, len(info.UriParams))
	assert.Equal(t, "TenantName", info.UriParams[0].Name)
	assert.Equal(t, 0, len(info.IndexOfInjectors))
	assert.Equal(t, false, info.HasInjector)     // No inject
	assert.Equal(t, 0, len(info.FormUploadFile)) // fileupload is one
	assert.Equal(t, -1, info.IndexOfRequestBody)
	assert.Equal(t, 0, len(info.IndexOfAuthClaims))
	assert.Nil(t, info.IndexOfAuthClaims)
	assert.Equal(t, -1, info.IndexOfAuthClaimsArg)
}
