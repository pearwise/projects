package utils

func StatusText(code int) string {
	switch code {
	case 200:
		return " OK"
	case 400:
        return " Bad Request"
	case 401:
        return " Unauthorized"
	case 403:
        return " Forbidden"
	case 404:
		return " Not Found"
	case 500:
		return " Internal Server Error"
	case 101:
		return " Switching Protocols"
	default:
		return " "
	}
}
