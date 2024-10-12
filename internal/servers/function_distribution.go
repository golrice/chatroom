package servers

type FunctionTable map[string]interface{}

type FunctionDistribution interface {
	Register(name string, f interface{})
	Get(name string) interface{}
}

// 全局下的函数表
var globalFunctionTable = FunctionTable{}

type functionDistribution struct {
}

func (fd *functionDistribution) Register(name string, f interface{}) {
	globalFunctionTable[name] = f
}

func (fd *functionDistribution) Get(name string) interface{} {
	return globalFunctionTable[name]
}
