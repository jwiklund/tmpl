package tmpl

func (fsRoot *fsRoot) Create(env Environment, target Target) error {
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
		err = WriteSingle(env, source, target)
		if err != nil {
			return err
		}
	}
	return nil
}
