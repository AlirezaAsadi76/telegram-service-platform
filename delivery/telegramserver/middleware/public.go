package middleware

func Public() []Middleware {
	return []Middleware{
		Recover(),
		Logger(),
	}
}
