package idgen

type IdGen128 interface {
	Generate() [2]int64
	GenerateString() string
	GenerateWithError() ([2]int64, error)
	GenerateStringWithError() (string, error)
}

func NewIdGen(nodeId int16, elementId int32) IdGen128 {
	return newIdGen(nodeId, elementId)
}

func Generate128() [2]int64 {
	return defaultGen.Generate()
}

func Generate128String() string {
	return defaultGen.GenerateString()
}

func Generate128WithError() ([2]int64, error) {
	return defaultGen.GenerateWithError()
}

func Generate128StringWithError() (string, error) {
	return defaultGen.GenerateStringWithError()
}
