package vapi

type inspectorType struct {
	helper *helperType
}

var inspector *inspectorType

func init() {
	inspector = &inspectorType{
		helper: &helperType{
			SpecialCharForRegex: "/\\?.$%^*-+",
		},
	}
}
