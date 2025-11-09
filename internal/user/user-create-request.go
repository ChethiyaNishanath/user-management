package user

type UserCreateRequest struct {
	FirstName string     `validate:"required,min=2,max=50"`
	LastName  string     `validate:"required,min=2,max=50"`
	Email     string     `validate:"required,email"`
	Phone     string     `validate:"omitempty,e164"`
	Age       int16      `validate:"omitempty,gt=0"`
	Status    UserStatus `validate:"omitempty,userStatus"`
}
