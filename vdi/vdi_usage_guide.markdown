# Hướng Dẫn Sử Dụng Thư Viện VDI

Thư viện `vdi` là một thư viện Dependency Injection (DI) cho Go, cho phép quản lý các dependency theo các lifecycle khác nhau: **Singleton**, **Scoped**, và **Transient**. Tài liệu này hướng dẫn cách sử dụng thư viện dựa trên các ví dụ kiểm thử, bao gồm cách đăng ký và truy xuất các dependency.

## 1. Giới thiệu về VDI

Thư viện `vdi` hỗ trợ:
- **Singleton**: Một instance duy nhất được chia sẻ trên toàn bộ ứng dụng.
- **Scoped**: Một instance duy nhất trong phạm vi (scope) cụ thể, các scope khác nhau sẽ có instance khác nhau.
- **Transient**: Mỗi lần truy xuất sẽ tạo một instance mới.

Các tính năng chính:
- Tạo container gốc (`RootContainer`) để quản lý các dependency.
- Tạo scope để giới hạn vòng đời của các dependency.
- Đăng ký và truy xuất dependency theo type.

## 2. Cài đặt

Giả định bạn đã cài đặt thư viện `vdi`. Để thêm thư viện vào dự án Go, sử dụng:

```bash
go get github.com/<path-to-vdi>
```

Thêm các import cần thiết trong mã:

```go
import (
    vdi "vdi"
)
```

## 3. Cách sử dụng

### 3.1. Khởi tạo Container

Sử dụng `vdi.NewRootContainer()` để tạo một container gốc. Container này là điểm bắt đầu để đăng ký và quản lý các dependency.

```go
root := vdi.NewRootContainer()
```

### 3.2. Đăng ký Dependency

Thư viện hỗ trợ ba loại lifecycle:

#### 3.2.1. Singleton
- **Đặc điểm**: Một instance duy nhất được chia sẻ trên toàn bộ ứng dụng, bất kể scope.
- **Cách đăng ký**: Sử dụng `RegisterSingleton` với một hàm factory trả về instance của dependency.

```go
root.RegisterSingleton(func() *Logger {
    return &Logger{ID: "singleton"}
})
```

- **Ứng dụng**: Phù hợp với các dependency cần duy trì trạng thái toàn cục, ví dụ: kết nối cơ sở dữ liệu, cấu hình ứng dụng.

#### 3.2.2. Scoped
- **Đặc điểm**: Một instance duy nhất trong một scope cụ thể. Mỗi scope mới sẽ tạo một instance mới.
- **Cách đăng ký**: Sử dụng `RegisterScoped` với một hàm factory.

```go
root.RegisterScoped(func() *Logger {
    return &Logger{ID: "scoped"}
})
```

- **Ứng dụng**: Phù hợp với các dependency cần giới hạn trong một phiên làm việc, ví dụ: thông tin người dùng trong một request HTTP.

#### 3.2.3. Transient
- **Đặc điểm**: Mỗi lần truy xuất sẽ tạo một instance mới.
- **Cách đăng ký**: Sử dụng `RegisterTransient` với một hàm factory.

```go
root.RegisterTransient(func() *Logger {
    return &Logger{ID: "transient"}
})
```

- **Ứng dụng**: Phù hợp với các dependency không cần duy trì trạng thái, ví dụ: các đối tượng tạm thời hoặc logger cho mỗi thao tác.

### 3.3. Tạo Scope

Tạo một scope từ container gốc bằng `CreateScope()`. Scope cho phép giới hạn vòng đời của các dependency `Scoped`.

```go
scope := root.CreateScope()
```

### 3.4. Truy xuất Dependency

Sử dụng `ResolveByType` để truy xuất một dependency dựa trên type của nó. Type được lấy bằng `vdi.TypeOf[T]()`, với `T` là kiểu của dependency.

```go
t := vdi.TypeOf[*Logger]()
logger, err := scope.ResolveByType(t)
if err != nil {
    // Xử lý lỗi
}
```

- **Lưu ý**: 
  - Với **Singleton**, cùng một instance được trả về ở mọi scope và container gốc.
  - Với **Scoped**, cùng một instance được trả về trong cùng một scope, nhưng khác scope sẽ trả về instance khác.
  - Với **Transient**, mỗi lần gọi `ResolveByType` sẽ trả về một instance mới.

### 3.5. Ví dụ minh họa

#### 3.5.1. Singleton
```go
root := vdi.NewRootContainer()
root.RegisterSingleton(func() *Logger {
    return &Logger{ID: "singleton"}
})

scope1 := root.CreateScope()
scope2 := root.CreateScope()

l1, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]())
l2, _ := scope2.ResolveByType(vdi.TypeOf[*Logger]())
l3, _ := root.ResolveByType(vdi.TypeOf[*Logger]())

// l1, l2, l3 là cùng một instance
```

#### 3.5.2. Scoped
```go
root := vdi.NewRootContainer()
root.RegisterScoped(func() *Logger {
    return &Logger{ID: "scoped"}
})

scope1 := root.CreateScope()
scope2 := root.CreateScope()

l1, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]())
l2, _ := scope1.ResolveByType(vdi.TypeOf[*Logger]()) // cùng instance với l1
l3, _ := scope2.ResolveByType(vdi.TypeOf[*Logger]()) // instance khác

// l1 == l2, nhưng l1 != l3
```

#### 3.5.3. Transient
```go
root := vdi.NewRootContainer()
root.RegisterTransient(func() *Logger {
    return &Logger{ID: "transient"}
})

scope := root.CreateScope()

l1, _ := scope.ResolveByType(vdi.TypeOf[*Logger]())
l2, _ := scope.ResolveByType(vdi.TypeOf[*Logger]())

// l1 != l2
```

## 4. Kiểm thử hiệu suất

Thư viện hỗ trợ kiểm thử hiệu suất thông qua benchmark. Dưới đây là các bài kiểm tra hiệu suất từ code mẫu:

- **BenchmarkTestScoped_ReuseWithinScope**: Đo hiệu suất của việc truy xuất dependency `Scoped` trong cùng một scope (instance được tái sử dụng).
- **BenchmarkTestTransient_AlwaysNew**: Đo hiệu suất của việc truy xuất dependency `Transient` (instance mới mỗi lần gọi).
- **BenchmarkTestSingleton_SameEverywhere**: Đo hiệu suất của việc truy xuất dependency `Singleton` trên các scope khác nhau (luôn trả về cùng instance).

Để chạy benchmark:

```bash
go test -bench=.
```

## 5. Lưu ý khi sử dụng

- **Kiểm tra lỗi**: Luôn kiểm tra lỗi trả về từ `ResolveByType` để xử lý trường hợp dependency không được đăng ký.
- **Type an toàn**: Sử dụng `vdi.TypeOf[T]()` để đảm bảo truy xuất đúng kiểu dependency.
- **Quản lý scope**: Đảm bảo scope được sử dụng hợp lý để tránh rò rỉ bộ nhớ, đặc biệt với các dependency `Scoped`.

## 6. Tích hợp với ứng dụng

- **Web server**: Sử dụng scope cho mỗi request HTTP để quản lý dependency `Scoped` (ví dụ: thông tin người dùng, logger).
- **Ứng dụng lớn**: Sử dụng `Singleton` cho các tài nguyên toàn cục như kết nối cơ sở dữ liệu hoặc cấu hình.
- **Kiểm thử**: Sử dụng `Transient` để tạo các đối tượng độc lập trong các bài kiểm thử.

## 7. Kết luận

Thư viện `vdi` cung cấp một cách đơn giản và hiệu quả để quản lý dependency trong Go với ba lifecycle: **Singleton**, **Scoped**, và **Transient**. Với API dễ sử dụng và khả năng kiểm thử hiệu suất, `vdi` phù hợp cho cả ứng dụng nhỏ và lớn.

Nếu bạn cần hỗ trợ thêm hoặc ví dụ chi tiết hơn, hãy liên hệ hoặc tham khảo mã nguồn của thư viện.