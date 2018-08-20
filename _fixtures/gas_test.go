package a

func badFunction() string {
	u, _ := ErrorHandle()
	return u
}

func ErrorHandle() (u string, err error) {
	return u
}
