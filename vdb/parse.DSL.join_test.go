package vdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	/*
		builder := SqlBuilder.With(
				`(d:departments, u:users)->d.id=u.departmentId` //<-- FROM department d LEFT JOIN users u as d ON d.id = u.department_id
			)
	*/
	// info, err := ParseDslJoin(`(d:departments, u:users)->d.id=u.departmentId`)
	// assert.NoError(t, err)
	// assert.Equal(t, "FROM departments AS d LEFT JOIN users AS u ON d.id=u.departmentId", info.FullSQL())

	// info1, err := ParseDslJoin(`(d:departments, u:users)<-d.id=u.departmentId`)
	// assert.NoError(t, err)
	// assert.Equal(t, "FROM departments AS d RIGHT JOIN users AS u ON d.id=u.departmentId", info1.FullSQL())
	// `(d:departments, u:users)*-d.id=u.departmentId` left outer join
	// `(d:departments, u:users)-*d.id=u.departmentId` left outer join
	// `(d:departments, u:users)**d.id=u.departmentId` full outer join

	info2, err := ParseDslJoin(`(d:departments, u:users)<->d.id=u.departmentId`)
	assert.NoError(t, err)
	assert.Equal(t, "FROM departments AS d RIGHT JOIN users AS u ON d.id=u.departmentId", info2.FullSQL())

	//fmt.Printf("FROM %s AS %s\n", info.FromTable, info.FromAlias)
	//fmt.Println(info.JoinString())
	// fmt.Println(info.FullSQL())
}
func Test_ParseDslJoin_ThreeTables(t *testing.T) {
	dsl := `(d:departments, u:users, c:checks)->d.id=u.departmentId->u.checkId=c.id`
	info, err := ParseDslJoin(dsl)
	assert.NoError(t, err)

	expectedSQL := "FROM departments AS d LEFT JOIN users AS u ON d.id=u.departmentId LEFT JOIN checks AS c ON u.checkId=c.id"
	assert.Equal(t, expectedSQL, info.FullSQL())
}

/*
├── Documents
│   ├── Report.doc
│   └── Data.xls
└── Programs
    └── App.exe

├── Table1
│   ├── Table2
│   └── Table13
└── Programs
    └── App.exe
*/
