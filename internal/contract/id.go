package contract

const IDKey = "app:id"

type IDService interface {
	NewID() string
}
