using BenchmarkDotNet.Attributes;
using BenchmarkDotNet.Running;
using System;

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
