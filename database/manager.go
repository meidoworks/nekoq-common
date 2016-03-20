package database

import "errors"

var engineMap = make(map[string]DatabaseEngine)
var managerConfig = make(map[string]interface{})

type DatabaseEngine interface {
	GetEngineName() string
	NewDatabase() (Database, error)
}

func InitManagerConfig(config map[string]interface{}) {
	managerConfig = config
}

func GetManagerConfig() map[string]interface{} {
	return managerConfig
}

func Register(engine DatabaseEngine) error {
	_, ok := engineMap[engine.GetEngineName()]
	if !ok {
		engineMap[engine.GetEngineName()] = engine
	} else {
		return errors.New("engine exists. name: " + engine.GetEngineName())
	}
	return nil
}

func GetEngine(name string) (DatabaseEngine, error) {
	engine, ok := engineMap[name]
	if !ok {
		return nil, errors.New("engine not exists. name: " + name)
	}
	return engine, nil
}
