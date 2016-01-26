package tmpl

func (fsRoot *fsRoot) Create(target Target) error {
	list, err := fsRoot.List(FilterFileAllowAll)
	if err != nil {
		return err
	}
	for _, tmpl := range list {
		source, err := fsRoot.Reader(tmpl)
		if err != nil {
			return err
		}
		target, err := target.Writer(tmpl)
		if err != nil {
			return err
		}
		var props map[string]string
		err = WriteSingle(source, target, props)
		if err != nil {
			return err
		}
	}
	return nil
}
