package services

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileService struct {
	DirectorySvc *DirectoryService
	UrlSvc       *UrlResolverService
}

func (fs *FileService) New() error {
	print("FileService.New")
	return nil
}

// Hàm mới để tìm số tiếp theo cho tên file
func (fs *FileService) findNextFileNumber(dirPath string) (int, error) {
	// Mở thư mục
	dir, err := os.Open(dirPath)
	if err != nil {
		// Trả về số 1 nếu thư mục rỗng hoặc không tồn tại (sẽ được tạo sau đó)
		if os.IsNotExist(err) {
			return 1, nil
		}
		return 0, err
	}
	defer dir.Close()

	// Đọc tất cả các tên file trong thư mục
	files, err := dir.Readdirnames(-1)
	if err != nil {
		return 0, err
	}

	maxNumber := 0
	for _, fileName := range files {
		// Lấy tên file không có phần mở rộng
		baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		// Chuyển đổi tên file thành số
		num, err := strconv.Atoi(baseName)
		if err == nil {
			// Nếu là một số hợp lệ và lớn hơn maxNumber, cập nhật maxNumber
			if num > maxNumber {
				maxNumber = num
			}
		}
	}

	// Số tiếp theo sẽ là số lớn nhất + 1
	return maxNumber + 1, nil
}
func (fs *FileService) GetFilePath(filePath string) (string, error) {
	ret := fs.DirectorySvc.DirUpload + "/" + filePath
	asbFilepath, err := filepath.Abs(ret)
	if err != nil {
		return "", err
	}

	return asbFilepath, nil
}
func (fs *FileService) SaveFile(file *multipart.FileHeader) (string, error) {
	// 1. Tạo hoặc lấy đường dẫn thư mục
	dirPath, err := fs.DirectorySvc.CreateDirectory()
	if err != nil {
		return "", err
	}

	// 2. Tìm số tiếp theo cho tên file từ các file hiện có
	newFileNumber, err := fs.findNextFileNumber(dirPath)
	if err != nil {
		return "", err
	}

	// 3. Định dạng tên file mới với số và phần mở rộng
	fileExt := filepath.Ext(file.Filename)
	if fileExt == "" {
		fileExt = ".dat"
	}
	newFileName := fmt.Sprintf("%04d%s", newFileNumber, fileExt)

	// 4. Tạo đường dẫn đầy đủ cho file mới
	destinationFilePath := filepath.Join(dirPath, newFileName)

	// 5. Mở file nhận được từ request và tạo file đích
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(destinationFilePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 6. Copy nội dung file
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return destinationFilePath, nil
}
func (fs *FileService) ListAllFiles() ([]string, error) {
	folderList, err := fs.DirectorySvc.ListAllDirectories()
	if err != nil {
		return nil, err
	}
	fileList := []string{}
	for _, folder := range folderList {
		//get all files in folder and append to fileList
		files, err := os.ReadDir(folder)
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			continue
		}

		for _, file := range files {
			fileList = append(fileList, filepath.Join(folder, file.Name()))
		}
	}

	// for i := 0; i < len(fileList); i++ {
	// 	fileList[i]=strings.TrimPrefix(fileList[i], directoryService.DirUpload[2:len(directoryService.DirUpload)-1)]))
	// }
	for i := 0; i < len(fileList); i++ {
		fileList[i] = strings.ReplaceAll(fileList[i], fs.DirectorySvc.DirUploadName+"\\", "")
		fileList[i] = fs.UrlSvc.MakeAbsUrl("api/media/files/" + fileList[i])
	}

	return fileList, nil

}

type DirectoryService struct {
	DirUpload     string
	DirUploadName string
	FullPathOfDir string
}

func (ds *DirectoryService) New() error {
	fmt.Println("DirectoryService.New")
	ds.DirUpload = "./uploads"

	fullPathOfDir, err := filepath.Abs(ds.DirUpload)
	if err != nil {
		return err
	}
	ds.FullPathOfDir = fullPathOfDir
	fmt.Println(fullPathOfDir)

	return nil
}
func (ds *DirectoryService) CreateDirectory() (string, error) {
	// get current time in UTC
	t := time.Now().UTC()
	// create directory YYYY
	datePath := t.Format("2006/01/02")

	// Kết hợp đường dẫn gốc với đường dẫn ngày tháng
	fullPath := filepath.Join(ds.DirUpload, datePath)

	// Tạo thư mục
	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		log.Fatalf("Không thể tạo thư mục: %v", err)
	}
	fullFilePath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	fmt.Printf("Đã tạo thư mục thành công: %s\n", fullFilePath)

	return fullPath, nil

}
func (ds *DirectoryService) ListAllDirectories() ([]string, error) {
	// Tạo slice để chứa đường dẫn của tất cả các thư mục
	var directories []string

	// Sử dụng filepath.Walk để duyệt qua toàn bộ cây thư mục
	// Tham số đầu tiên là thư mục gốc cần duyệt
	// Tham số thứ hai là một hàm callback sẽ được gọi cho mỗi file/thư mục
	err := filepath.Walk(ds.DirUpload, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Xử lý lỗi nếu có vấn đề khi truy cập thư mục/file
			return err
		}

		// Kiểm tra xem entry hiện tại có phải là thư mục không và không phải thư mục gốc
		if info.IsDir() && path != ds.DirUpload {
			// Nếu là thư mục, thêm đường dẫn vào slice
			directories = append(directories, path)
		}

		// Trả về nil để tiếp tục quá trình duyệt
		return nil
	})

	if err != nil {
		return nil, err
	}

	return directories, nil
}

type UrlResolverService struct {
	BaseUrl string
}

func (ur *UrlResolverService) New() error {
	fmt.Println("UrlResolver.New")

	return nil

}
func (ur *UrlResolverService) MakeAbsUrl(s string) string {
	fx := ur.BaseUrl + "/" + strings.ReplaceAll(s, "\\", "/")
	return fx

}
