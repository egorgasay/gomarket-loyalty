package service

import "gomarket-loyalty/model"

type ProductService interface {
	Create(request model.CreateProductRequest) (response model.CreateProductResponse)
	List() (responses []model.GetProductResponse)
}
