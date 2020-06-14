package miuires

type ElementHelper interface {
	GetName() (name string)
	GetValue() (value string)
	Parse(base string) (ok bool)
	Write() []byte
}
