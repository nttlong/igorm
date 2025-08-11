# Hướng Dẫn Sử Dụng Thư Viện `vapi` Trong Go

Thư viện `vapi` là một công cụ mạnh mẽ giúp đơn giản hóa việc xây dựng API trong Go. Nó cung cấp các tính năng như tự động tạo URI, hỗ trợ xác thực, xử lý file, và sinh tài liệu Swagger mà không cần code-gen. Hướng dẫn này sẽ giải thích cách sử dụng `vapi` dựa trên các ví dụ thực tế.

## 1. Tổng Quan Về `vapi`

- **Mục đích**: `vapi` giúp lập trình viên định nghĩa các API handler một cách đơn giản, tự động ánh xạ các phương thức của struct thành các endpoint API.
- **Tính năng chính**:
  - Định nghĩa handler theo hai cách: trực tiếp và gián tiếp.
  - Tự động tạo URI dựa trên tên package, tên struct, và tên phương thức (theo chuẩn kebab-case).
  - Hỗ trợ xác thực thông qua `vapi.AuthClaims`.
  - Tích hợp Swagger để sinh tài liệu API.
  - Hỗ trợ upload và streaming file.
  - Cho phép sử dụng middleware để xử lý request.

## 2. Cài Đặt Và Thiết Lập Cơ Bản

1. **Cài đặt**:
   - Đảm bảo bạn đã thêm thư viện `vapi` vào dự án Go của mình:
     ```bash
     go get github.com/vapi
     ```
   - Import các gói cần thiết:
     ```go
     import (
         "github.com/google/uuid"
         "vapi"
         "os"
         "io"
         "mime/multipart"
         "path/filepath"
     )
     ```

2. **Khởi tạo server**:
   - Tạo một HTTP server với base URI và port:
     ```go
     server := vapi.NewHtttpServer("/api/v1", 8080, "localhost")
     ```
   - Base URI (`/api/v1`) sẽ được thêm vào trước các URI tương đối của handler.

3. **Đăng ký controller**:
   - Đăng ký các controller (struct chứa các handler) bằng `vapi.Controller`:
     ```go
     vapi.Controller(func() (*example.Media, error) {
         return &example.Media{}, nil
     })
     ```

4. **Bật Swagger (tùy chọn)**:
   - Gọi `server.Swagger()` để tự động sinh tài liệu Swagger:
     ```go
     server.Swagger()
     ```

5. **Thêm middleware**:
   - Sử dụng middleware để xử lý các yêu cầu HTTP, ví dụ: logging, CORS, hoặc nén dữ liệu:
     ```go
     server.Middleware(mw.Cors)
     server.Middleware(mw.Zip)
     ```

6. **Khởi động server**:
   - Bắt đầu server và xử lý lỗi:
     ```go
     err := server.Start()
     if err != nil {
         panic(err)
     }
     ```

## 3. Định Nghĩa Handler

`vapi` hỗ trợ hai cách định nghĩa handler: **trực tiếp** và **gián tiếp**.

### 3.1. Cách 1: Trực Tiếp

- **Cách khai báo**: Thêm tham số `*vapi.Handler` hoặc `vapi.Handler` vào phương thức của struct.
- **URI tự động**: URI được tạo theo cú pháp: `<package>/<struct>/<method>` (kebab-case).
- **Phương thức HTTP mặc định**: `POST`.
- **Ví dụ**: Liệt kê danh sách file trong thư mục:
  ```go
  type Media struct {
      vapi.Service
  }

  func (m *Media) ListOfFiles(ctx *vapi.Handler) ([]string, error) {
      folder := "./uploads"
      files, err := os.ReadDir(folder)
      if err != nil {
          return nil, err
      }
      ret := []string{}
      for _, file := range files {
          ret = append(ret, m.BaseUrl+"/"+ctx.BaseUrl+"/example/media/"+file.Name())
      }
      return ret, nil
  }
  ```
  - URI: `/api/v1/example/media/list-of-files` (nếu base URI là `/api/v1`).
  - Trả về danh sách các file trong thư mục `Uploads`.

### 3.2. Cách 2: Gián Tiếp

- **Cách khai báo**: Nhúng `vapi.Handler` vào một struct và sử dụng tag `route` để tùy chỉnh phương thức HTTP và URI.
- **Khi nào sử dụng**:
  - Cần thay đổi phương thức HTTP (ví dụ: `GET`, `PUT`).
  - Yêu cầu xác thực.
  - Cần URI động (ví dụ: chứa placeholder như `{FileName}`).
- **Ví dụ**: Xem nội dung file với phương thức `GET`:
  ```go
  func (m *Media) File(ctx struct {
      vapi.Handler `route:"method:get;uri:@/{FileName}"`
      FileName string
  }) error {
      fileName := "./Uploads/" + ctx.FileName
      return ctx.StreamingFile(fileName)
  }
  ```
  - URI: `/api/v1/example/media/<tên_file>` (ví dụ: `/api/v1/example/media/my-video-file.mp4`).
  - Placeholder `{FileName}` được ánh xạ vào field `FileName` của struct.

## 4. Xác Thực (Authentication)

- **Khai báo xác thực**: Nhúng `*vapi.AuthClaims` vào struct handler để yêu cầu xác thực.
- **Ví dụ**:
  ```go
  type AuthHandler struct {
      vapi.Handler
      Auth *vapi.AuthClaims
  }

  func (m *Media) DoShareFile(ctx AuthHandler, data struct {
      FilesShare []string `json:"files_share"`
      ShareTo    string   `json:"share_to"`
      ShareType  string   `json:"share_type"`
  }) (*UploadResult, error) {
      return &UploadResult{UploadId: uuid.New().String()}, nil
  }
  ```
  - API yêu cầu xác thực và hiển thị biểu tượng khóa trong Swagger.
  - URI: `/api/v1/example/media/do-share-file`.

## 5. Xử Lý File

### 5.1. Upload File

- **Khai báo**: Sử dụng `multipart.FileHeader` hoặc `multipart.File` trong struct body. Các trường file phải nằm ở cấp cao nhất (không lồng trong struct con).
- **Ví dụ**:
  ```go
  func (m *Media) Upload(ctx *AuthHandler, data struct {
      Files []*multipart.FileHeader `json:"file"`
      NoteFile multipart.File
      Info struct {
          FolderId string `json:"folder_id"`
      }
  }) ([]string, error) {
      if data.Files == nil {
          return nil, fmt.Errorf("file is required")
      }
      ret := []string{}
      for _, file := range data.Files {
          uploadDir := "./Uploads/"
          if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
              return nil, fmt.Errorf("không tạo được thư mục upload: %w", err)
          }
          f, err := file.Open()
          if err != nil {
              return nil, err
          }
          defer f.Close()
          out, err := os.Create(filepath.Join(uploadDir, file.Filename))
          if err != nil {
              return nil, err
          }
          defer out.Close()
          if _, err = io.Copy(out, f); err != nil {
              return nil, err
          }
      }
      return ret, nil
  }
  ```
  - URI: `/api/v1/example/media/upload`.
  - Hỗ trợ upload nhiều file và hiển thị trường upload trong Swagger.

### 5.2. Streaming File

- Sử dụng `ctx.StreamingFile` để trả về file dưới dạng stream, tiết kiệm bộ nhớ:
  ```go
  return ctx.StreamingFile(fileName)
  ```

## 6. URI Cố Định

- **Khai báo**: Sử dụng dấu `/` ở đầu URI trong tag `route` để định nghĩa URI cố định, không phụ thuộc vào package/struct.
- **Ví dụ**:
  ```go
  func (a *Auth) Oauth(ctx *struct {
      vapi.Handler `route:"uri:/api/@/token"`
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
      ret := struct {
          AccessToken  string `json:"access_token"`
          TokenType    string `json:"token_type"`
          ExpiresIn    int    `json:"expires_in"`
          RefreshToken string `json:"refresh_token"`
          Scope        string `json:"scope"`
      }{
          AccessToken: "12345556",
          TokenType: "Bearer",
          ExpiresIn: 3600,
      }
      return &ret, nil
  }
  ```
  - URI: `/api/oauth/token` (không bị ảnh hưởng bởi base URI).

## 7. Hỗ Trợ Swagger

- Khi bật Swagger (`server.Swagger()`), `vapi` tự động sinh tài liệu API, bao gồm:
  - Mô tả các tham số, body, và response của API.
  - Hiển thị biểu tượng khóa cho các API yêu cầu xác thực.
  - Hỗ trợ mô tả các trường upload file (`multipart.FileHeader`).

## 8. Lưu Ý Quan Trọng

- **URI động**: Sử dụng `@` trong tag `route` để đại diện cho đường dẫn mặc định của handler (`<package>/<struct>/<method>`).
- **File upload**: Các trường file (`multipart.File` hoặc `multipart.FileHeader`) phải nằm ở cấp cao nhất của struct body.
- **Middleware**: Sử dụng middleware để thêm các chức năng như CORS, logging, hoặc nén dữ liệu.
- **Xác thực**: Nhúng `*vapi.AuthClaims` để yêu cầu xác thực cho handler.
- **Khởi tạo server**: Đảm bảo base URI và port được cấu hình đúng.

## 9. Kết Luận

Thư viện `vapi` là một lựa chọn tuyệt vời để xây dựng API trong Go với mã nguồn tối giản và hiệu quả. Nó tự động hóa nhiều tác vụ như tạo URI, sinh Swagger, và xử lý file, đồng thời cung cấp sự linh hoạt trong việc định nghĩa handler và xác thực. Bằng cách sử dụng `vapi`, bạn có thể nhanh chóng tạo ra các API mạnh mẽ mà không cần các công cụ bổ sung.