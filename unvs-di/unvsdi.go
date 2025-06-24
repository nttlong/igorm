package unvsdi

import (
	"fmt"
	"reflect"
	"sync"
)

type Singleton[TOwner any, T any] struct {
	Value T
	Owner interface{}
	Init  func(owner *TOwner) T
	once  sync.Once
}

// The life cycle of a service is controlled by the owner.
type Scoped[TOwner any, T any] struct {
	Value T
	Owner interface{}
	Init  func(owner *TOwner) T
}

// The life cycle of a service is controlled by the container.
type Transient[TOwner any, T any] struct {
	Value T
	Owner interface{}
	Init  func(owner *TOwner) T
}

func (s *Transient[TOwner, T]) Get() T {
	if s.Owner == nil {
		panic("Transient[TOwner, T] requires an owner")
	}
	return s.Init(s.Owner.(*TOwner))
}

func (s *Singleton[TOwner, T]) Get() T {
	if s.Owner == nil {
		panic("Singleton[TOwner, T] requires an owner")
	}
	s.once.Do(func() {
		s.Value = s.Init(s.Owner.(*TOwner))
	})
	return s.Value
}

func (s *Scoped[TOwner, T]) Get() T {
	if s.Owner == nil {
		panic("Scoped[TOwner, T] requires an owner")
	}
	if s.Init == nil {
		return s.Value
	}
	s.Value = s.Init(s.Owner.(*TOwner))
	s.Init = nil
	return s.Value
}

// resolveByType resolves dependencies of a given struct type by setting the Owner and Init fields of injector fields.
// resolveByType 解析给定结构类型的依赖关系，通过设置注入器字段的 Owner 和 Init 字段。
// resolveByType 関数は、指定された構造体型の依存関係を解決し、Injector フィールドの Owner および Init フィールドを設定します。
// resolveByType giải quyết phụ thuộc của kiểu struct bằng cách gán giá trị cho Owner và Init của các injector field.
func resolveByType(typ reflect.Type, retVal reflect.Value, visited map[reflect.Type]bool) (*reflect.Value, error) {
	// Detect circular dependency
	// 检测循环依赖
	// 循環依存を検出
	// Phát hiện vòng lặp phụ thuộc
	if visited[typ] {
		panic(fmt.Sprintf("circular dependency detected on type: %s", typ.String()))
	}
	visited[typ] = true
	defer delete(visited, typ) // Clean up when leaving this recursion // 递归结束时清理 // 再帰終了時にクリーンアップ // Dọn dẹp khi đệ quy kết thúc

	meta := getMetaOfType(typ) // Get registered sample instance // 获取已注册的样本实例 // 登録済みサンプルインスタンスを取得 // Lấy sample instance đã đăng ký (nếu có)
	if meta == nil {
		return nil, fmt.Errorf("type %s is not registered", typ.String())
	}
	// init owner

	for fieldName, field := range meta.injectorFields {

		// Handle embedded (anonymous) fields

		if field.Anonymous {
			fieldValue := meta.sampleFieldValues[fieldName]

			ft := field.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				fieldValue = fieldValue.Elem()
			}
			_, err := resolveByType(ft, fieldValue, visited)
			if err != nil {
				return nil, err
			}
		}

		if field.Type.Kind() == reflect.Struct {
			fieldType := field.Type

			if utils.isInjector(field) {
				sampleFieldValue := meta.sampleFieldValues[fieldName]

				// Set owner field for the injector

				ownerField := retVal.Field(field.Index[0]).FieldByName("Owner")
				if ownerField.IsValid() && ownerField.CanSet() {
					ownerField.Set(retVal.Addr())
				}
				sampleValueFieldInitFn := sampleFieldValue.FieldByName("Init")
				initField := retVal.Field(field.Index[0]).FieldByName("Init")
				initField.Set(sampleValueFieldInitFn)
				// // Copy Init function from registered sample if available
				// // 如果有已注册的样本，则复制 Init 函数
				// // 登録されたサンプルから Init 関数をコピー（可能な場合）
				// // Nếu có sample đã đăng ký, copy hàm Init từ sample
				// if meta.sampleValue != nil {
				// 	fieldInSample := (*sampleValue).FieldByName(field.Name)
				// 	if fieldInSample.IsValid() {
				// 		fnVal := fieldInSample.FieldByName("Init")
				// 		if fnVal.IsValid() && fnVal.Type().Kind() == reflect.Func {
				// 			initField := val.FieldByName("Init")
				// 			if initField.IsValid() && initField.CanSet() {
				// 				initField.Set(fnVal)
				// 			}
				// 		}
				// 	}
				// }

			} else {
				// Continue resolving normal (non-injector) struct fields
				// 继续解析普通（非注入器）结构字段
				// 通常の（Injectorではない）構造体フィールドの解決を継続
				// Tiếp tục resolve field kiểu struct bình thường (không phải injector)
				fieldValue := retVal.Field(field.Index[0])
				if fieldValue.CanAddr() {
					resolveByType(fieldType, fieldValue, visited)
				}
			}
		}
	}

	retInterface := retVal.Addr()
	return &retInterface, nil
}

// Resolve resolves an instance of type T.
// Resolve 解析类型 T 的实例。
// Resolve 関数は型 T のインスタンスを解決します。
// Resolve tạo instance của kiểu T.
func Resolve[T any]() (*T, error) {
	var t T
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		panic("Resolve[T] requires a struct type")
	}
	retVal := reflect.New(typ).Elem()
	visited := make(map[reflect.Type]bool) // Track visited types to detect circular dependencies // 追踪访问的类型以检测循环依赖 // 循環依存を検出するために訪問した型を追跡 // Theo dõi kiểu đã visit để phát hiện vòng lặp
	retInterface, err := resolveByType(typ, retVal, visited)
	if err != nil {
		return nil, err
	}
	return retInterface.Interface().(*T), nil
}
