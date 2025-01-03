// Code generated by github.com/ycl2018/go-conv DO NOT EDIT.

package cases

import (
	"github.com/ycl2018/go-conv/testdata/a"
	"github.com/ycl2018/go-conv/testdata/b"
)

func PtrAStructToPtrBStruct(src *a.Struct) (dst *b.Struct) {
	if src != nil {
		dst = new(b.Struct)
		dst.Student.Name = src.Student.Name
		dst.Student.Class.Name = src.Student.Class.Name
		dst.Student.Class.Grade = string(src.Student.Class.Grade)
		dst.Student.Teachers = src.Student.Teachers
		dst.Student2.Name = src.Student2.Name
		dst.Student2.Class.Name = src.Student2.Class.Name
		dst.Student2.Class.Grade = transfer(src.Student2.Class.Grade)
		filteredSrcStudent2Teachers := filter(src.Student2.Teachers)
		dst.Student2.Teachers = filteredSrcStudent2Teachers
		dst.Match_ = src.Match
	}
	return
}
func init() {
	Struct2Struct = PtrAStructToPtrBStruct
}
