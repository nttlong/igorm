package orm

func (e *expression) Prepare(input string) (string, error) {
	if e.keywords == nil {
		e.keywords = []string{
			"select",
			"from",
			"where",
			"group",
			"order",
			"limit",
			"offset",
		}
	}
	for _, keyword := range e.keywords {
		markList, err := e.getMarkList(input, keyword)
		if err != nil {
			return "", err
		}
		input = e.InsertMarks(input, markList)
	}
	return input, nil

}
