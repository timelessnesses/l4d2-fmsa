package export

func ExportFMSA(path string) error {
	err := ExportJSON(path)
	if err != nil {
		return err
	}
	return nil
}
