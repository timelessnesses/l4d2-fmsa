package export

import "strings"

func Export(path string) error {
	if strings.Contains(path, ".fmsa") {
		return ExportFMSA(path)
	} else if strings.Contains(path, ".json") {
		return ExportJSON(path)
	} else {
		return ExportText(path)
	}
}
