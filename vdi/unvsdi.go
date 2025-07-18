package vdi

import (
	"fmt"
	"reflect"
	"sync"
)

type Singleton[TOwner any, T any] struct {
	Value T
	Owner interface{}
	Init  func(owner TOwner) T
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
		s.Value = s.Init(s.Owner.(TOwner))
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

	for fieldName, field := range meta.mappingFields {

		// Handle embedded (anonymous) fields
		if field.Anonymous {
			embeddedValue := retVal.Field(field.Index[0])
			if embeddedValue.CanAddr() {
				resolveByType(field.Type, embeddedValue, visited)
			}
			continue
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

			} else {

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

type initResolve struct {
	once sync.Once
	err  error
	val  *reflect.Value
}

var cacheResolve sync.Map

func resolve(typ reflect.Type) (*reflect.Value, error) {
	actual, _ := cacheResolve.LoadOrStore(typ, &initResolve{})
	resolve := actual.(*initResolve)
	resolve.once.Do(func() {
		retVal := reflect.New(typ).Elem()
		visited := make(map[reflect.Type]bool) // Track visited types to detect circular dependencies // 追踪访问的类型以检测循环依赖 // 循環依存を検出するために訪問した型を追跡 // Theo dõi kiểu đã visit để phát hiện vòng lặp
		retInterface, err := resolveByType(typ, retVal, visited)
		resolve.err = err
		resolve.val = retInterface

	})
	return resolve.val, resolve.err
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
	retVale, err := resolve(typ)
	if err != nil {
		return nil, err
	}
	return retVale.Interface().(*T), nil

}
