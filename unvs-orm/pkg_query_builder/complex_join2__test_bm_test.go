package pkgquerybuilder

import (
	"testing"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func BenchmarkAliAsTableSelfJoin(t *testing.B) {
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < t.N; i++ {

		childDept := repo.Departments.Alias("childDept") //<-- self join need an alias to avoid table name conflict

		sql := childDept.ParentId.RightJoin(repo.Departments.DepartmentId).Select(
			childDept.ParentId,
			childDept.Name,
		)
		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(t, compilerResult.Err)
		sqlExpected := "SELECT [T1].[parent_id] AS [parent_id], [T1].[name] AS [name] FROM [departments] AS [T2] RIGHT JOIN [departments] AS [T1] ON [T1].[parent_id] = [T2].[department_id]"
		assert.Equal(t, sqlExpected, compilerResult.SqlText)
	}
}
func BenchmarkSelfJoin3Levels_RightJoin(t *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < t.N; i++ {
		child := repo.Departments.Alias("child")   // T1
		parent := repo.Departments.Alias("parent") // T2
		grand := repo.Departments.Alias("grand")   // T3
		join := child.RightJoin(parent,
			parent.DepartmentId.Eq(child.ParentId),
		)
		join = join.RightJoin(
			grand,
			grand.DepartmentId.Eq(parent.ParentId),
		)

		sql := join.Select(
			child.Name.As("ChildName"),
			parent.Name.As("ParentName"),
			grand.Name.As("GrandName"),
		)

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(t, compilerResult.Err)

		expected := "SELECT [T3].[name] AS [ChildName], [T2].[name] AS [ParentName], [T1].[name] AS [GrandName] FROM [departments] AS [T3] RIGHT JOIN [departments] AS [T2] ON [T2].[department_id] = [T3].[parent_id] RIGHT JOIN [departments] AS [T1] ON [T1].[department_id] = [T2].[parent_id]"

		assert.Equal(t, expected, compilerResult.SqlText)
	}
}
func BenchmarkJoinEmpDapartment(b *testing.B) {
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		dept := repo.Departments
		sql := emp.DepartmentId.Eq(dept.DepartmentId).Select(
			emp.Name.As("EmpName"),
			dept.Name.As("DeptName"),
		)
		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)
		expected := "SELECT [T1].[name] AS [EmpName], [T2].[name] AS [DeptName] FROM [employees] AS [T1] INNER JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id]"
		assert.Equal(b, expected, compilerResult.SqlText)
	}
}
func BenchmarkLeftJoinEmpDepartment(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		dept := repo.Departments

		sql := emp.DepartmentId.
			LeftJoin(dept.DepartmentId).
			Select(
				emp.Name.As("EmpName"),
				dept.Name.As("DeptName"), // Có thể null nếu không thuộc phòng
			)

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)

		expected := "SELECT [T1].[name] AS [EmpName], [T2].[name] AS [DeptName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id]"

		assert.Equal(b, expected, compilerResult.SqlText)
	}
}
func BenchmarkEmployeeWithoutDepartment(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		dept := repo.Departments

		sql := emp.DepartmentId.
			LeftJoin(dept.DepartmentId).
			// lọc nhân viên chưa có phòng
			Select(
				emp.Name.As("EmpName"),
			).Where(dept.DepartmentId.IsNull())

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)

		expected := "SELECT [T1].[name] AS [EmpName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id] " +
			"WHERE [T2].[department_id] IS NULL"

		assert.Equal(b, expected, compilerResult.SqlText)
	}
}
func BenchmarkDepartmentWithoutEmployee(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		dept := repo.Departments
		emp := repo.Employees

		sql := dept.DepartmentId.
			LeftJoin(emp.DepartmentId). // join theo FK
			Select(
				dept.Name.As("DeptName"),
			).Where(emp.EmployeeId.IsNull()) // lọc phòng không có nhân viên

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)

		expected := "SELECT [T1].[name] AS [DeptName] " +
			"FROM [departments] AS [T1] " +
			"LEFT JOIN [employees] AS [T2] ON [T1].[department_id] = [T2].[department_id] " +
			"WHERE [T2].[employee_id] IS NULL"

		assert.Equal(b, expected, compilerResult.SqlText)
	}
}
func BenchmarkJoinFullEmployeeDepartment_LeftJoin(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		dept := repo.Departments

		sql := emp.DepartmentId.
			LeftJoin(dept.DepartmentId).
			Select(
				emp.EmployeeId,
				emp.Name.As("EmpName"),
				emp.Email,
				emp.Phone,
				dept.Name.As("DeptName"),
				dept.Note,
			)

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)

		expected := "SELECT [T1].[employee_id] AS [EmployeeId], [T1].[name] AS [EmpName], [T1].[email] AS [Email], [T1].[phone] AS [Phone], [T2].[name] AS [DeptName], [T2].[note] AS [Note] FROM [employees] AS [T1] LEFT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id]"

		assert.Equal(b, expected, compilerResult.SqlText)
	}
}
func BenchmarkJoinEmployeeWithManager(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		manager := repo.Employees.Alias("manager") // Self-join alias tránh conflict tên bảng

		sql := emp.ManagerId.
			LeftJoin(manager.EmployeeId).
			Select(
				emp.Name.As("EmpName"),
				manager.Name.As("ManagerName"), // Có thể null nếu không có quản lý
			)

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err)

		expected := "SELECT [T1].[name] AS [EmpName], [T2].[name] AS [ManagerName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [employees] AS [T2] ON [T1].[manager_id] = [T2].[employee_id]"

		assert.Equal(b, expected, compilerResult.SqlText)
	}
}
