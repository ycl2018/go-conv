package a

type Struct struct {
	Student  Student
	Student2 Student
	Student3 Student
	Pojo     *Pojo
	Match    string
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
	Grade int
}
