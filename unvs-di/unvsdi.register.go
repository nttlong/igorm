package unvsdi

import (
	"fmt"
	"reflect"
	"sync"
)

type metaInfo struct {
	sampleInstanceReflectValue reflect.Value
	mappingFields              map[string]reflect.StructField // Cache injector fields
	sampleFieldValues          map[string]reflect.Value       // Cache sample field values
	initFuncs                  map[string]reflect.Value       // Cache Init functions
}

type containerInfo[T any] struct {
	sampleInstance T
	meta           metaInfo
}

func (r *containerInfo[T]) init(resolver func(obj *T) error) error {
	typ := reflect.TypeFor[T]()
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	_objVal := reflect.New(typ)
	obj := _objVal.Elem().Interface().(T)
	err := resolver(&obj)
	if err != nil {
		return err
	}
	objVal := reflect.ValueOf(obj)

	if objVal.Kind() == reflect.Ptr {
		objVal = objVal.Elem()
	}
	for i := 0; i < objVal.NumField(); i++ {
		field := objVal.Type().Field(i)
		if utils.isInjector(field) {
			r.meta.mappingFields[field.Name] = field
			r.meta.initFuncs[field.Name] = objVal.FieldByName("Init")
			r.meta.sampleFieldValues[field.Name] = objVal.Field(i)
		} else {
			r.meta.mappingFields[field.Name] = field
			r.meta.sampleFieldValues[field.Name] = objVal.Field(i)
		}
	}
	r.meta.sampleInstanceReflectValue = objVal
	r.sampleInstance = obj

	return nil

}
func (r *containerInfo[T]) GetInitFun(PropertyName string) (*reflect.Value, error) {
	val := reflect.ValueOf(r.sampleInstance)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	field, ok := val.Type().FieldByName(PropertyName)
	if !ok {
		return nil, fmt.Errorf("field %s not found in type %s", PropertyName, val.Type().String())
	}
	if !utils.isInjector(field) {
		return nil, fmt.Errorf("field %s is not an injector", PropertyName)
	}
	fn := val.FieldByName("Init")

	return &fn, nil

}

// var make sure the containerInfo[T] is only created once for each type T
// var containersCache = sync.Map{}

// var globalMeta : each containerInfo has meta map,
// this value to store reflect.ValueOf(sampleInstance),
var globalMeta = sync.Map{}

func getMetaOfType(typ reflect.Type) *metaInfo {
	key := typ.String()
	if v, ok := globalMeta.Load(key); ok {
		ret := v.(metaInfo)
		return &ret
	}
	return nil
}

// RegisterContainer registers a new container for a given type T.
//
// RegisterContainer 为指定类型 T 注册一个新的容器。
//
// RegisterContainer は指定された型 T のコンテナを登録します。
//
// RegisterContainer đăng ký container mới cho kiểu T.
func RegisterContainer[T any](resolver func(svc *T) error) (*containerInfo[T], error) {
	var t T
	key := reflect.TypeOf(t).String() // Get the type key as a string // 获取类型键字符串 // 型キーの文字列を取得 // Lấy key dạng chuỗi của kiểu

	// Check if a container for this type has already been registered
	// 检查该类型的容器是否已被注册
	// この型のコンテナが既に登録されているか確認
	// Kiểm tra container cho kiểu này đã được đăng ký chưa
	if _, ok := globalMeta.Load(key); ok {
		return nil, fmt.Errorf("container for type %s already exists", key)
	}

	ret := &containerInfo[T]{
		meta: metaInfo{
			mappingFields:              make(map[string]reflect.StructField),
			initFuncs:                  make(map[string]reflect.Value),
			sampleFieldValues:          make(map[string]reflect.Value),
			sampleInstanceReflectValue: reflect.ValueOf(t),
		},
	} // Create new containerInfo instance // 创建新的 containerInfo 实例 // 新しい containerInfo インスタンスを作成 // Tạo containerInfo mới

	// Initialize the container using the provided resolver function
	// 使用提供的 resolver 函数初始化容器
	// 渡された resolver 関数を使ってコンテナを初期化
	// Khởi tạo container bằng hàm resolver được truyền vào
	err := ret.init(resolver)
	if err != nil {
		return nil, err // Initialization failed // 初始化失败 // 初期化失敗 // Khởi tạo thất bại
	}

	// Store metadata for later reference (for injection purposes)
	// 存储元数据以供后续引用（用于注入）
	// 後で参照するためのメタデータを保存（注入用）
	// Lưu metadata cho lần resolve sau (dùng cho injection)
	globalMeta.Store(key, ret.meta)

	// Store the container instance in the cache
	// 将容器实例存储到缓存中
	// コンテナインスタンスをキャッシュに保存
	// Lưu container instance vào cache

	return ret, nil // Return the created container // 返回创建的容器 // 作成されたコンテナを返す // Trả về container vừa tạo
}
