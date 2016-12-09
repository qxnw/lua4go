package lua4go

//Context 脚本执行上下文
type Context struct {
	Session     string
	LoggerName  string
	Input       string
	Body        string
	HTTPContext interface{}
	DataMap     map[string]interface{}
}
