using BenchmarkDotNet.Attributes;
using BenchmarkDotNet.Running;
using System;
using System.Data;
using System.Data.Common;
using Microsoft.EntityFrameworkCore;                  // đã có
using Microsoft.EntityFrameworkCore.Infrastructure;   // cần thêm
using System.Data.Common;
using BenchmarkDotNet.Order;
using BenchmarkDotNet.Engines;

[MemoryDiagnoser]
[GcForce] // ép GC chạy giữa các iteration
[MinColumn, MaxColumn, MeanColumn, MedianColumn]
[RankColumn]
[WarmupCount(5)]
[IterationCount(10)]
[InvocationCount(1)]
[ThreadingDiagnoser]
[Orderer(SummaryOrderPolicy.FastestToSlowest)]
[BenchmarkCategory("Insert")]
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
        //db.Database.EnsureDeleted();
        //db.Database.EnsureCreated();
    }

    [Benchmark]
    public void InsertPositionDepartmentUserEmployee()
    {
        for (int i = 0; i < 50000; i++)
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
        var rows = 20000;
        var q = (from e in db.Employees
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
                      .Take(rows);


        var start = DateTime.Now;
        var item = q.ToList();
        var end = DateTime.Now;
        Console.WriteLine($"Elapsed: {(end - start).TotalMilliseconds} ms");
        Console.WriteLine($"Count: {item.Count()}");
        if (item.Count() != rows)
        {
            throw new Exception("Count mismatch");
        }


    }
}
