package shelldiff

func must(_ any, err error) {
	if err != nil {
		panic(err)
	}
}
