package vgrpc

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	// Thay thế bằng đường dẫn module của bạn
	"vgrpc/caller"
)

// packValue là hàm hỗ trợ để đóng gói một giá trị bất kỳ vào *anypb.Any
func packValue(input interface{}) (*anypb.Any, error) {
	// Kiểm tra nếu input là một Protobuf message (ưu tiên)
	if p, ok := input.(proto.Message); ok {
		return anypb.New(p)
	}

	// Xử lý các kiểu dữ liệu generic với reflection
	v := reflect.ValueOf(input)
	switch v.Kind() {
	case reflect.Slice:
		// Xử lý slice byte đặc biệt
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return anypb.New(wrapperspb.Bytes(v.Bytes()))
		}
		// Xử lý slice thông thường bằng cách đóng gói vào AnyList
		anyList := &caller.AnyList{}
		for j := 0; j < v.Len(); j++ {
			packedElem, err := packValue(v.Index(j).Interface())
			if err != nil {
				return nil, fmt.Errorf("không thể đóng gói phần tử slice thứ %d: %w", j, err)
			}
			anyList.Values = append(anyList.Values, packedElem)
		}
		return anypb.New(anyList)

	case reflect.String:
		return anypb.New(wrapperspb.String(v.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return anypb.New(wrapperspb.Int64(v.Int()))
	case reflect.Bool:
		return anypb.New(wrapperspb.Bool(v.Bool()))
	case reflect.Float32, reflect.Float64:
		return anypb.New(wrapperspb.Double(v.Float()))
	default:
		return nil, fmt.Errorf("kiểu dữ liệu không được hỗ trợ: %T", input)
	}
}

// packInputs converts a slice of Go interfaces into a slice of Any messages.
func packInputs(inputs ...interface{}) ([]*anypb.Any, error) {
	anyInputs := make([]*anypb.Any, len(inputs))
	for i, input := range inputs {
		anyMsg, err := packValue(input)
		if err != nil {
			return nil, fmt.Errorf("không thể đóng gói input thứ %d: %w", i, err)
		}
		anyInputs[i] = anyMsg
	}
	return anyInputs, nil
}

// unpackValue là hàm hỗ trợ để giải nén một giá trị từ *anypb.Any
func unpackValue(anyMsg *anypb.Any) (interface{}, error) {
	protoMessage, err := anyMsg.UnmarshalNew()
	if err != nil {
		return nil, fmt.Errorf("lỗi khi giải nén Any message: %w", err)
	}

	// Kiểm tra nếu là AnyList
	if list, ok := protoMessage.(*caller.AnyList); ok {
		return unpackAnySlice(list.GetValues())
	}

	// Xử lý các kiểu Wrapper types
	if wrapper, ok := protoMessage.(*wrapperspb.StringValue); ok {
		return wrapper.GetValue(), nil
	} else if wrapper, ok := protoMessage.(*wrapperspb.Int64Value); ok {
		return wrapper.GetValue(), nil
	} else if wrapper, ok := protoMessage.(*wrapperspb.BoolValue); ok {
		return wrapper.GetValue(), nil
	} else if wrapper, ok := protoMessage.(*wrapperspb.DoubleValue); ok {
		return wrapper.GetValue(), nil
	} else if wrapper, ok := protoMessage.(*wrapperspb.BytesValue); ok {
		return wrapper.GetValue(), nil
	}

	// Trả về chính Protobuf message nếu nó không phải là kiểu Wrapper
	return protoMessage, nil
}

// unpackAnySlice converts a slice of Any messages back to a slice of Go interfaces.
func unpackAnySlice(anyResults []*anypb.Any) ([]interface{}, error) {
	results := make([]interface{}, len(anyResults))
	for i, anyResult := range anyResults {
		unpackedValue, err := unpackValue(anyResult)
		if err != nil {
			return nil, fmt.Errorf("lỗi khi giải nén kết quả thứ %d: %w", i, err)
		}
		results[i] = unpackedValue
	}
	return results, nil
}
