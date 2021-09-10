package tibero

func cellValueUpdate(model IModel, fieldWhere string) (string, string, error) {
	var fields, err = GetValueFields(KEY_TAG, model)
	if err != nil {
		return "", "", err
	}
	var cells, where string
	for _, val := range fields {
		if val.Name == fieldWhere {
			where = val.Value
		}
		if val.Insert {
			if cells == "" {
				cells += val.Name + "=" + val.Value
			} else {
				cells += "," + val.Name + "=" + val.Value
			}
		}
	}
	return cells, where, err
}
