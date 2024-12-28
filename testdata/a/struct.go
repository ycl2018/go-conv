package a

type Struct struct {
	Student  Student
	Student2 Student
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
