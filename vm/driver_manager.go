package vm

//NewVMDriver providing default VM driver (currently otto)
//Another VM can be used by modifying this file to provide
//a desirable implmentation of VMDriver
func NewVMDriver() Driver {

	return &OttoDriver{}
}

// func GetVMDriver(funcs map[int]BuiltInFunctions) VMDriver {

// 	return &OttoVMDriver{funcs: funcs}
// }
