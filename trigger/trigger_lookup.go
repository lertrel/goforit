package trigger

//Lookup Trigger repository
type Lookup interface {

	//GetTrigger getting Trigger by name
	GetTrigger(triggerName string) (Trigger, error)

	//GetAllTriggers getting all Trigger(s)
	GetAllTriggers() (Iterator, error)

	//GetTriggers search all Trigger(s) that matches the given filter (script)
	GetTriggers(filter string) (Iterator, error)
}

//Iterator a interator of Trigger
type Iterator interface {
	HasNext() bool
	Next() Trigger
}
