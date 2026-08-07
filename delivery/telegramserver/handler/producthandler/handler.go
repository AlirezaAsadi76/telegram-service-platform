package producthandler

type Handler struct {
	productService ProductService
}

func New(productService ProductService) Handler {
	return Handler{
		productService: productService,
	}
}
