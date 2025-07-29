package vdi

import (
	"fmt"
	"reflect"
)

// resolveByType resolves dependencies of a given struct type by setting the Owner and Init fields of injector fields.
// resolveByType è§£æžç»™å®šç»“æž„ç±»åž‹çš„ä¾èµ–å…³ç³»ï¼Œé€šè¿‡è®¾ç½®æ³¨å…¥å™¨å­—æ®µçš„ Owner å’Œ Init å­—æ®µã€‚
// resolveByType é–¢æ•°ã¯ã€æŒ‡å®šã•ã‚ŒãŸæ§‹é€ ä½“åž‹ã®ä¾å­˜é–¢ä¿‚ã‚’è§£æ±ºã—ã€Injector ãƒ•ã‚£ãƒ¼ãƒ«ãƒ‰ã® Owner ãŠã‚ˆã³ Init ãƒ•ã‚£ãƒ¼ãƒ«ãƒ‰ã‚’è¨­å®šã—ã¾ã™ã€‚
// resolveByType giáº£i quyáº¿t phá»¥ thuá»™c cá»§a kiá»ƒu struct báº±ng cÃ¡ch gÃ¡n giÃ¡ trá»‹ cho Owner vÃ  Init cá»§a cÃ¡c injector field.
func resolveByType(typ reflect.Type, retVal reflect.Value, visited map[reflect.Type]bool) (*reflect.Value, error) {
	// Detect circular dependency
	// æ£€æµ‹å¾ªçŽ¯ä¾èµ–
	// å¾ªç’°ä¾å­˜ã‚’æ¤œå‡º
	// PhÃ¡t hiá»‡n vÃ²ng láº·p phá»¥ thuá»™c
	if visited[typ] {
		panic(fmt.Sprintf("circular dependency detected on type: %s", typ.String()))
	}
	visited[typ] = true
	defer delete(visited, typ) // Clean up when leaving this recursion // é€’å½’ç»“æŸæ—¶æ¸…ç† // å†å¸°çµ‚äº†æ™‚ã«ã‚¯ãƒªãƒ¼ãƒ³ã‚¢ãƒƒãƒ— // Dá»n dáº¹p khi Ä‘á»‡ quy káº¿t thÃºc

	meta := getMetaOfType(typ) // Get registered sample instance // èŽ·å–å·²æ³¨å†Œçš„æ ·æœ¬å®žä¾‹ // ç™»éŒ²æ¸ˆã¿ã‚µãƒ³ãƒ—ãƒ«ã‚¤ãƒ³ã‚¹ã‚¿ãƒ³ã‚¹ã‚’å–å¾— // Láº¥y sample instance Ä‘Ã£ Ä‘Äƒng kÃ½ (náº¿u cÃ³)
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

// Resolve resolves an instance of type T.
// Resolve è§£æžç±»åž‹ T çš„å®žä¾‹ã€‚
// Resolve é–¢æ•°ã¯åž‹ T ã®ã‚¤ãƒ³ã‚¹ã‚¿ãƒ³ã‚¹ã‚’è§£æ±ºã—ã¾ã™ã€‚
// Resolve táº¡o instance cá»§a kiá»ƒu T.
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
	visited := make(map[reflect.Type]bool) // Track visited types to detect circular dependencies // è¿½è¸ªè®¿é—®çš„ç±»åž‹ä»¥æ£€æµ‹å¾ªçŽ¯ä¾èµ– // å¾ªç’°ä¾å­˜ã‚’æ¤œå‡ºã™ã‚‹ãŸã‚ã«è¨ªå•ã—ãŸåž‹ã‚’è¿½è·¡ // Theo dÃµi kiá»ƒu Ä‘Ã£ visit Ä‘á»ƒ phÃ¡t hiá»‡n vÃ²ng láº·p
	retInterface, err := resolveByType(typ, retVal, visited)
	if err != nil {
		return nil, err
	}
	return retInterface.Interface().(*T), nil
}
