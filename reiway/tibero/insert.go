package tibero

func cellValueInserts(model IModel) (string, string, error) {
	var fields, err = GetValueFields(KEY_TAG, model)
	if err != nil {
		return "", "", err
	}
	var cells, values string
	for _, val := range fields {
		if val.Insert {
			if cells == "" {
				cells += val.Name
				values += val.Value
			} else {
				cells += "," + val.Name
				values += "," + val.Value
			}
		}
	}
	return cells, values, err
}
