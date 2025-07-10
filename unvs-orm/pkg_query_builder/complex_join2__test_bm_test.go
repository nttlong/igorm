package pkgquerybuilder

import (
	"testing"
	"time"
	orm "unvs-orm"

	"github.com/stretchr/testify/assert"
)

func BenchmarkComplexWhereQuery(t *testing.B) {
	fx := orm.Utils.ToSnakeCase("EmployeeCount")
	assert.Equal(t, "employee_count", fx)
	repo := orm.Repository[OrderRepository]() // Giả sử vẫn dùng Departments
	for i := 0; i < t.N; i++ {
		// Tạo một ngày cụ thể, ví dụ: 1/1/2023
		specificDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

		sql := repo.Departments.
			Filter(
				repo.Departments.Name.Like("Sales%").
					Or(
						repo.Departments.EmployeeCount.Gt(50).
							And(repo.Departments.CreatedAt.Gt(specificDate)),
					),
			).
			Select(repo.Departments.EmployeeCount) // Chọn tất cả các cột

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(t, compilerResult.Err())

		// SQL dự kiến (có thể cần điều chỉnh tùy theo cách ORM sinh ra)
		sqlExpected := "SELECT [employeecount] AS [EmployeeCount], [department_id] AS [DepartmentId], [name] AS [Name], [parent_id] AS [ParentId], [updatedat] AS [UpdatedAt], [order_no] AS [OrderNo], [note] AS [Note], [createdat] AS [CreatedAt], [updated_at] AS [UpdatedAt], [departmentid] AS [DepartmentId], [code] AS [Code], [level] AS [Level], [created_by] AS [CreatedBy], [createdby] AS [CreatedBy], [updated_by] AS [UpdatedBy], [employee_count] AS [EmployeeCount], [parentid] AS [ParentId], [orderno] AS [OrderNo], [created_at] AS [CreatedAt], [updatedby] AS [UpdatedBy] FROM [departments] WHERE [departments].[name] LIKE ? OR [departments].[employee_count] > ? AND [departments].[created_at] > ?"
		assert.Equal(t, sqlExpected, compilerResult.String())
		// Ngoài ra, kiểm tra các tham số được truyền vào (parameters) cũng rất quan trọng
		// assert.Equal(t, []interface{}{"Sales%", 50, specificDate}, compilerResult.Parameters)
	}
}
func BenchmarkAliAsTableSelfJoin(t *testing.B) {
	repo := orm.Repository[OrderRepository]()
	for i := 0; i < t.N; i++ {

		childDept := repo.Departments.Alias("childDept") //<-- self join need an alias to avoid table name conflict

		sql := childDept.RightJoin(
			repo.Departments,
			childDept.DepartmentId.Eq(repo.Departments.ParentId),
		).Select(
			childDept.Name.As("ChildName"),
			repo.Departments.Name.As("ParentName"),
		)
		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(t, compilerResult.Err())
		sqlExpected := "SELECT [T1].[name] AS [ChildName], [T2].[name] AS [ParentName] FROM [departments] AS [T1] RIGHT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[parent_id]"
		assert.Equal(t, sqlExpected, compilerResult.String())
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
		assert.NoError(t, compilerResult.Err())

		expected := "SELECT [T3].[name] AS [ChildName], [T2].[name] AS [ParentName], [T1].[name] AS [GrandName] FROM [departments] AS [T3] RIGHT JOIN [departments] AS [T2] ON [T2].[department_id] = [T3].[parent_id] RIGHT JOIN [departments] AS [T1] ON [T1].[department_id] = [T2].[parent_id]"

		assert.Equal(t, expected, compilerResult.String())
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
		assert.NoError(b, compilerResult.Err())
		expected := "SELECT [T1].[name] AS [EmpName], [T2].[name] AS [DeptName] FROM [employees] AS [T1] INNER JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id]"
		assert.Equal(b, expected, compilerResult.String())
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
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T1].[name] AS [EmpName], [T2].[name] AS [DeptName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id]"

		assert.Equal(b, expected, compilerResult.String())
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
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T1].[name] AS [EmpName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id] " +
			"WHERE [T2].[department_id] IS NULL"

		assert.Equal(b, expected, compilerResult.String())
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
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T1].[name] AS [DeptName] " +
			"FROM [departments] AS [T1] " +
			"LEFT JOIN [employees] AS [T2] ON [T1].[department_id] = [T2].[department_id] " +
			"WHERE [T2].[employee_id] IS NULL"

		assert.Equal(b, expected, compilerResult.String())
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
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T1].[employee_id] AS [EmployeeId], [T1].[name] AS [EmpName], [T1].[email] AS [Email], [T1].[phone] AS [Phone], [T2].[name] AS [DeptName], [T2].[note] AS [Note] FROM [employees] AS [T1] LEFT JOIN [departments] AS [T2] ON [T1].[department_id] = [T2].[department_id]"

		assert.Equal(b, expected, compilerResult.String())
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
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T1].[name] AS [EmpName], [T2].[name] AS [ManagerName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [employees] AS [T2] ON [T1].[manager_id] = [T2].[employee_id]"

		assert.Equal(b, expected, compilerResult.String())
	}
}
func BenchmarkEmployeeWithoutManager(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		manager := repo.Employees.Alias("manager")

		sql := emp.ManagerId.
			LeftJoin(manager.EmployeeId).
			Select(emp.Name.As("EmpName")).
			Where(manager.EmployeeId.IsNull())

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T1].[name] AS [EmpName] " +
			"FROM [employees] AS [T1] " +
			"LEFT JOIN [employees] AS [T2] ON [T1].[manager_id] = [T2].[employee_id] " +
			"WHERE [T2].[employee_id] IS NULL"

		assert.Equal(b, expected, compilerResult.String())
	}
}
func BenchmarkEmployeeWithGrandManager(b *testing.B) {
	repo := orm.Repository[OrderRepository]()

	for i := 0; i < b.N; i++ {
		emp := repo.Employees
		mgr := repo.Employees.Alias("mgr")
		grand := repo.Employees.Alias("grand")

		join := emp.Join(mgr, emp.ManagerId.Eq(mgr.EmployeeId)) //<--emp.ManagerId.LeftJoin(mgr.EmployeeId)
		join = join.LeftJoin(grand, grand.EmployeeId.Eq(mgr.ManagerId))

		sql := join.Select(
			emp.Name.As("EmpName"),
			mgr.Name.As("ManagerName"),
			grand.Name.As("GrandManagerName"),
		)

		dialect := mssql()
		compilerResult := sql.Compile(dialect)
		assert.NoError(b, compilerResult.Err())

		expected := "SELECT [T3].[name] AS [EmpName], [T2].[name] AS [ManagerName], [T1].[name] AS [GrandManagerName] FROM [employees] AS [T3] INNER JOIN [employees] AS [T2] ON [T3].[manager_id] = [T2].[employee_id] LEFT JOIN [employees] AS [T1] ON [T1].[employee_id] = [T2].[manager_id]"

		assert.Equal(b, expected, compilerResult.String())
		//"SELECT [T3].[name] AS [EmpName], [T2].[name] AS [ManagerName], [T1].[name] AS [GrandManagerName] FROM [employees] AS [T3] INNER JOIN [employees] AS [T2] ON [T3].[manager_id] = [T2].[employee_id] LEFT JOIN [employees] AS [T1] ON [T1].[employee_id] = [T2].[manager_id]"
	}
}
