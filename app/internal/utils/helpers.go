package utils

// Helpers simples usados em múltiplos arquivos.
func AwsString(s string) *string { return &s }
func AwsStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
