package user

type UserUpdateRequest struct {
	FirstName string `validate:"omitempty,min=2,max=50"`
	LastName  string `validate:"omitempty,min=2,max=50"`
	Email     string `validate:"omitempty,email"`
	Phone     string `validate:"omitempty,e164"`
	Age       int16  `validate:"omitempty,gt=0"`
	Status    string `validate:"omitempty,userStatus"`
}
