package example

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"wx"

	"github.com/google/uuid"
)

/*
Duoi day la mau khai bao he thong API, trong package exmaple1.
Mot vi du: ve quan ly file.
Su dung thu vien vapi, mot thu vien ho tro 1 cach don gian nhat de tao cac API handler
va tham chi tu dong sinh ra Swagger (neu bat Swagger) ma kg can phai Code Gen
*/
type Media struct {
	wx.ControllerContext
	User wx.UserClaims
}

/*
Co 2 cach de tao
*/

/*
#Cach 1: truc tiep, hay noi cach khac, la khai bao truc tiep 1 tham so co kieu
vapi.Handler hoac *vapi.Handler de bao cho vapi biet day la 1 ham
se duoc Handle boi web server
Luu y: phuong thuc mac dinh se la POST va Uri cua api la
<package name cua Media dang ToKebab Case>/<ten cua struct dong vai tro la Controller, cu the o day la Media, va cung o dang ToKebabCase>/
<Ten cua method>
Nhu bay trong truong hop cu the nay uri cua api la example/media/list-of-file
Luu y:"example/media/list-of-file" chua phai la Url cuoi cung de handler tai http server
*/
func (m *Media) ListOfFiles(ctx *wx.Handler, // <-- nhung method nao cua
	// Media co mot tham so kieu vapi.Handler hoac *vapi.Handler, thi duoc xem nhu la 1 handler
	// Nhu vay la o day ta da co 1 api liet ke danh sach cac file
	data struct {
		Page int `json:"page"`
		Size int `json:"size"`
	}) ([]string, error) {
	folder := "./uploads"
	streamFileUri, err := wx.GetUriOfHandler[Media]("ListOfFiles")
	if err != nil {
		return nil, err
	}
	//create folder if not exists
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.Mkdir(folder, 0750)
		if err != nil {
			return nil, err
		}
	}
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}
	ret := []string{}

	for _, file := range files {
		ret = append(ret, ctx.GetAbsRootUri()+streamFileUri+"/example/media/file/"+file.Name())
	}
	return ret, nil
}

/*
#Cach 2. Gian tiep. Khai bao 1 struct va embed vapi.Handler vao
Cach nay thuong duoc su dung khi: muon tao 1 handler bat buoc pahi chung thuc
Hoac thay doi method POST thanh GET,PUT,...
Trong vi du duoi day la tao 1 handler bat buoc phai chung thu moi duoc su dung
Luu y: moi Method cua struct co su dung AuthHandler (nhu trong truong hop nay) duoc xem nhu la co chung thuc
*/
type AuthHandler struct { // <-- co hai cach de khai bao 1 handler
	wx.Handler                // khai bao 1 struc va embed vapi.handler
	Auth       *wx.AuthClaims //<-- day la tuy chon, neu co 1 field vapi.AuthClaim, dieu nay co nghia la
	// Tat cac handler su dung AuthHandler bat buoc phai chung thuc, truong hop kg can chung thuc hay bo qua
	// va dung khai bao bat ky 1 field nao trong cau truc co kieu vapi.AuthClaims haoc * vapi.AuthClaim
}

/*
Api nay cho phep xem noi dung file.  De xem duoc file phai su dung method Http Get.
Nhu vay minh buoc phai su dung cach 2 "gian tiep"
Tuy nhien trong truong hop nay phuc tap hon 1 chut la ta phai khai bao 1 uri trong tag route
De co the truy cap dang example/media/my-video-file.mp4,..
De co the truy cap duoc url dang example/media/my-video-file.mp4 trong Uri se co 1 placeHolder map voi 1 Field trong Handler
Luu y: phai su dung ky tu @ de dai dien cho <package>/<struct name>/<method>
Vi du: tag uri @/{FileName} thi co nghia la truy cap example/media/my-video-file.mp4 de co duoc file
Neu dat {FileName}/@ ->my-video-file.mp4/example/media
*/
func (m *Media) File(ctx struct {
	wx.Handler `route:"method:get;uri:@/{FileName}"` //<-- tai sao la @/{FileName}?
	//
	//											^
	//				|--------------------------	|
	//										[Lay o cho nay]
	//			[Dat vao cho nay]
	FileName string //<-- vapi se tu dong map tu Url request vao bien nay
}) error {
	fileName := "./uploads/" + ctx.FileName
	return ctx.StreamingFile(fileName) //<-- StreamingFile la 1 ham da co san trong vapi.Handler)

}

type UploadResult struct {
	UploadId string
}

/*
Vid du share file
uri cho API la example/media/do-share-file
Http method Post
Body theo tam so data
Luu y: Neu bat Swagger len cac tham so se hien thi day du
*/
func (m *Media) DoShareFile( //<-- day la 1 vi du Dang ky 1 media
	ctx AuthHandler, //<-- cach 2 (gian tiep),  neu bat swagger len cho nay se co hinh O khoa
	data struct { //<-- body
		FilesShare []string `json:"files_share"`
		ShareTo    string   `json:"share_to"`
		ShareType  string   `json:"share_type"`
	},
) (*UploadResult, error) {
	return &UploadResult{
		UploadId: uuid.New().String(),
	}, nil
}
func (m *Media) Upload(ctx *AuthHandler, //<-- Van su dung cach 2 la gian tipe do Upload file doi hoi phai co  chung thuc
	data struct { // <-- day phan Body handler
		//
		Files []*multipart.FileHeader `json:"file"` //<-- Khi co bat cu field nao co kieu la
		//multipart.FileHeader hoac multipart.File, hoac la slice of multipart.File,hoac la slice of multipart.FileHeader
		// thi Swagger se tu dong xuat hien file Upload (neu bat Swagger)
		NoteFile multipart.File //<-- Muon upload file phai de ngoai cung
		Info     struct {       // Ngoai file ra cung co the yeu cau them 1 so thong tin khac, vi d nhu struc duoi day
			FolderId string `json:"folder_id"`
			// Luu y: Khi muon Upload file cac file pahi de ngoai cung, kg duoc long trong struct
			// NoteFile multipart.File //<-- khai bao nhu vay la kg hop le
		}
	}) ([]string, error) {
	if data.Files == nil {
		return nil, fmt.Errorf("file is required")
	}
	ret := []string{}
	for _, file := range data.Files {
		uploadDir := "./uploads/"

		// Tạo thư mục nếu chưa tồn tại
		if err := os.MkdirAll(filepath.Clean(uploadDir), 0750); err != nil {
			return nil, fmt.Errorf("không tạo được thư mục upload: %w", err)
		}

		f, err := file.Open() // file là *multipart.FileHeader
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// Tạo file đích
		out, err := os.Create(filepath.Join(uploadDir, file.Filename))
		if err != nil {
			return nil, err
		}
		defer out.Close()

		// Copy dữ liệu
		if _, err = io.Copy(out, f); err != nil {
			return nil, err
		}

	}
	return ret, nil
}

// phan tiep theo se noi ve cach dat uri co dinh cho API
// Uri co dinh la Uri kg co chua package path, bang cach them 1 dau '/' vao dau khai bao
// Vi du duoi day
type Auth struct {
}

func (a *Auth) Oauth(ctx *struct {
	wx.Handler `route:"uri:/api/@/token"` //--> se phat sinh la uri la /api/oauth/token
}, data struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}) (*struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}, error) {
	fmt.Println(data.Password)
	fmt.Println(data.UserName)
	ret := struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}{}
	ret.AccessToken = "12345556"
	ret.TokenType = "Bearer"
	ret.ExpiresIn = 3600
	return &ret, nil
}

/*
func main() {
	vapi.Controller(func() (*example.Media, error) { // Dang ky Controller Media
		return &example.Media{}, nil
	})
	vapi.Controller(func() (*example.Auth, error) { // Dang ky Controller Auth
		return &example.Auth{}, nil
	})
	server := vapi.NewHtttpServer("/api/v1", 8080, "localhost") // tao Http Server voi baseUri la /api/v1
	// Luu y baseUri chi se duoc dat vao truoc cac Uri cua Handler nhu da noi phan tren de tao ra
	// uri thuc su cho API, va dac biet chi ap dung cho cac Uri tuong doi tuc la cac Uri kg bat dau ban dau '/'
	// Vi du: controller Auth tren method Oauth se kg bi anh huong boi baseUri trong ham  vapi.NewHtttpServer
	server.MiddleWare(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		....
		Noi dung cua MiddleWare
	})
	server.Swagger() // Bat Swagger <-- la tuy chon
	server.Middleware(mw.LogAccessTokenClaims) //<-- khai bao cac Middleware

	server.Middleware(mw.Cors)
	server.Middleware(mw.Zip)
	err := server.Start()
	if err != nil {
		panic(err)
	}

}
*/
