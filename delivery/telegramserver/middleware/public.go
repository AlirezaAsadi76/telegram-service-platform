package middleware

func Public() []Middleware {
	return []Middleware{
		Logger(),
		Recover(),
	}
}
