package actions

type Action interface {
	Execute()
	Result() interface{}
}
