package types


type Student struct{
	ID int `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Age int `json:"age" validate:"required"`
}