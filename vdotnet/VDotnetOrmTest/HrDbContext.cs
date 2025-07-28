using Microsoft.EntityFrameworkCore;

public class HrDbContext : DbContext
{
    public DbSet<User> Users => Set<User>();
    public DbSet<Account> Accounts => Set<Account>();
    public DbSet<Department> Departments => Set<Department>();
    public DbSet<Position> Positions => Set<Position>();
    public DbSet<Employee> Employees => Set<Employee>();
    // Add các DbSet còn lại

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        //dotnet tool install --global dotnet-ef
        //dotnet ef migrations add InitCreate
        //dotnet ef database update
        /*
        dotnet ef migrations remove
dotnet ef migrations add InitialCreate
dotnet ef database update

        */



        // optionsBuilder.UseMySql(
        //     "server=localhost;user=root;password=123456;database=dotenet_test_001;",
        //     new MySqlServerVersion(new Version(8, 0, 34))
        // );
        optionsBuilder.UseSqlServer(
        "Server=localhost;Database=dot_net_core_test;User Id=sa;Password=123456;Encrypt=False"
        );
        // optionsBuilder.UseNpgsql(
        //     "Host=localhost;Database=dotenet_test_001;Username=postgres;Password=123456"
        // );
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        // modelBuilder.Entity<Department>()
        //     .HasIndex(d => d.Name)
        //     .IsUnique();

        // modelBuilder.Entity<Department>()
        //     .HasIndex(d => d.Code)
        //     .IsUnique();

        // modelBuilder.Entity<Department>()
        //     .HasOne(d => d.Parent)
        //     .WithMany()
        //     .HasForeignKey(d => d.ParentId)
        //     .OnDelete(DeleteBehavior.NoAction);

        // modelBuilder.Entity<Employee>()
        //     .HasOne(e => e.Department)
        //     .WithMany()
        //     .HasForeignKey(e => e.DepartmentId);

        // modelBuilder.Entity<Employee>()
        //     .HasOne(e => e.Position)
        //     .WithMany()
        //     .HasForeignKey(e => e.PositionId);

        // modelBuilder.Entity<Employee>()
        //     .HasOne(e => e.User)
        //     .WithMany()
        //     .HasForeignKey(e => e.UserId);
        modelBuilder.Entity<Account>()
        .HasIndex(a => a.Username)
        .IsUnique();

        base.OnModelCreating(modelBuilder);

        // Tiếp tục cấu hình các quan hệ khác nếu có
    }
}
