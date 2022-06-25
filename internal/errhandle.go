package internal

// MustBeNil panics if err isn't nil
func MustBeNil(err error) {
	if err != nil {
		panic(err)
	}
}
