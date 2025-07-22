using BenchmarkDotNet.Attributes;
using BenchmarkDotNet.Running;
using System;
using System.Data;
using System.Data.Common;
using Microsoft.EntityFrameworkCore;                  // đã có
using Microsoft.EntityFrameworkCore.Infrastructure;   // cần thêm
using System.Data.Common;    
[MemoryDiagnoser]
public class EfCoreInsertBenchmark
{
    private HrDbContext db;

    [GlobalSetup]
    public void Setup()
    {
        db = new HrDbContext();
        db.Database.EnsureDeleted();
        db.Database.EnsureCreated();
    }

    [Benchmark]
    public void InsertSingleUser()
    {
        var name = "ef_user_" + Guid.NewGuid().ToString("N");
        var email = $"{name}@dotnet.com";

        var user = new User
        {
            UserId = Guid.NewGuid().ToString(),
            Email = email,
            Phone = "0987654321",
            Username = name,
            HashPassword = "ef_pass_123456",
            Description = "EFCore benchmark user",
            CreatedAt = DateTime.UtcNow
        };

        db.Users.Add(user);
        db.SaveChanges();

        // Optional check
        if (user.Username != name)
        {
            throw new Exception("Username mismatch");
        }
    }
}
[MemoryDiagnoser]
public class EfCoreInsertTransactionBenchmark
{
    private HrDbContext db;

    [GlobalSetup]
    public void Setup()
    {
        db = new HrDbContext();
        db.Database.EnsureDeleted();
        db.Database.EnsureCreated();
    }

    [Benchmark]
    public void InsertPositionDepartmentUserEmployee()
    {
        var strIndex = Guid.NewGuid().ToString("N");

        var position = new Position
        {
            Name = "CEO" + strIndex,
            Code = "CEO0" + strIndex,
            Title = "Chief Executive Officer " + strIndex,
            Level = 1,
            Description = "test",
            CreatedAt = DateTime.UtcNow
        };

        var department = new Department
        {
            Name = "CEO" + strIndex,
            Code = "CEO0" + strIndex,
            Description = "test",
            CreatedAt = DateTime.UtcNow
        };

        var user = new User
        {
            UserId = Guid.NewGuid().ToString(),
            Email = $"test{strIndex}@test.com",
            Phone = "0987654321",
            Username = "test001" + strIndex,
            HashPassword = "123456" + strIndex,
            Description = "test",
            CreatedAt = DateTime.UtcNow
        };

        using var tx = db.Database.BeginTransaction();
        try
        {
            db.Positions.Add(position);
            db.Departments.Add(department);
            db.Users.Add(user);
            db.SaveChanges();

            var emp = new Employee
            {
                PositionId = position.Id,
                DepartmentId = department.Id,
                UserId = user.Id,
                FirstName = "John",
                LastName = "Doe",
                Description = "test",
                CreatedAt = DateTime.UtcNow
            };

            db.Employees.Add(emp);
            db.SaveChanges();
            tx.Commit();
        }
        catch
        {
            tx.Rollback();
            throw;
        }
    }
}
public class QueryResult
{
    public string FullName { get; set; } = string.Empty;
    public int PositionId { get; set; }
    public int DepartmentId { get; set; }
    public string Email { get; set; } = string.Empty;
    public string Phone { get; set; } = string.Empty;
}
[MemoryDiagnoser]
public class EfCoreRawSqlBenchmark
{
    private HrDbContext db = null!;

    [GlobalSetup]
    public void Setup()
    {
        db = new HrDbContext();
        db.Database.OpenConnection(); // Quan trọng: mở kết nối trước để dùng low-level
    }

    [Benchmark]
    public void BenchmarkRawSelect()
    {
        //     var sql = @"
        //         SELECT 
        //             CONCAT(e.first_name, ' ', e.last_name) AS FullName,
        //             e.position_id AS PositionId,
        //             e.department_id AS DepartmentId,
        //             u.email AS Email,
        //             u.phone AS Phone
        //         FROM employees e
        //         LEFT JOIN users u ON e.user_id = u.id
        //         ORDER BY e.id
        //         LIMIT 1000";

        //     using var command = db.Database.GetDbConnection().CreateCommand();
        //     command.CommandText = sql;
        //     command.CommandType = CommandType.Text;

        //     using var reader = command.ExecuteReader();
        //     var results = new List<QueryResult>(1000);
        //     while (reader.Read())
        //     {
        //         results.Add(new QueryResult
        //         {
        //             FullName = reader.GetString(0),
        //             PositionId = reader.GetInt32(1),
        //             DepartmentId = reader.GetInt32(2),
        //             Email = reader.GetString(3),
        //             Phone = reader.GetString(4)
        //         });
        //     }
        // }
        var results = (from e in db.Employees
                       join u in db.Users on e.UserId equals u.Id into joined
                       from u in joined.DefaultIfEmpty()
                       orderby e.Id
                       select new QueryResult
                       {
                           FullName = e.FirstName + " " + e.LastName,
                           PositionId = e.PositionId,
                           DepartmentId = e.DepartmentId,
                           Email = u != null ? u.Email : null,
                           Phone = u != null ? u.Phone : null
                       })
                      .AsNoTracking()
                      .Skip(0)
                      .Take(5000)
                      .ToList();
                      
        
    }
}
