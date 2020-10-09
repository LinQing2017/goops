package error

type NodeShellError struct {
	ErrCode int
	ErrMsg  string
}

func (err *NodeShellError) Error() string {
	return err.ErrMsg
}
