package kernel

type Email string

func NewEmail(email string) Email { return Email(email) }
func (e Email) String() string    { return string(e) }
func (e Email) IsEmpty() bool     { return string(e) == "" }

type Phone string

type FirstName string

type LastName string
