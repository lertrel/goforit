package goforit

//TriggerLookup Trigger repository
type TriggerLookup interface {

	//GetTrigger getting Trigger by name
	GetTrigger(triggerName string) (Trigger, error)

	//GetAllTriggers getting all Trigger(s)
	GetAllTriggers() (TriggerIterator, error)

	//GetTriggers search all Trigger(s) that matches the given filter (script)
	GetTriggers(filter string) (TriggerIterator, error)
}

//TriggerIterator a interator of Trigger
type TriggerIterator interface {
	HasNext() bool
	Next() Trigger
}
