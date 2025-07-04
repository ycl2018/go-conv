package b

type Struct struct {
	Student         Student
	Student2        Student
	Student3        Student
	Pojo            *Pojo
	Match_          string
	Caseinsensitive string
	IgnoreType      *Pojo
	IgnoreField     string
	Struct          struct {
		Foo struct {
			Bar string
		}
	}
}

type Pojo struct {
	Int    int
	String string
}

type Student struct {
	Name     string
	Class    Class
	Teachers []string
}

type Class struct {
	Name  string
	Grade string
}

type Foo struct {
	Bar Bar
}

type Bar struct {
	Field  string
	Field2 int
}
