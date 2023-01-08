package utils

func GetSubfix(s string) string {
	for i := len(s)-1;i>-1;i--{
		if s[i] == '.' {
            return s[i+1:]
        }
	}
	return ""
}