using Microsoft.EntityFrameworkCore;

public class HrDbContext : DbContext
{
    public DbSet<User> Users => Set<User>();
    public DbSet<Department> Departments => Set<Department>();
    public DbSet<Position> Positions => Set<Position>();
    public DbSet<Employee> Employees => Set<Employee>();
    // Add các DbSet còn lại

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        optionsBuilder.UseMySql(
            "server=localhost;user=root;password=123456;database=dotenet_test_001;",
            new MySqlServerVersion(new Version(8, 0, 34))
        );
    }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<Department>()
            .HasIndex(d => d.Name)
            .IsUnique();

        modelBuilder.Entity<Department>()
            .HasIndex(d => d.Code)
            .IsUnique();

        modelBuilder.Entity<Department>()
            .HasOne(d => d.Parent)
            .WithMany()
            .HasForeignKey(d => d.ParentId)
            .OnDelete(DeleteBehavior.NoAction);

        modelBuilder.Entity<Employee>()
            .HasOne(e => e.Department)
            .WithMany()
            .HasForeignKey(e => e.DepartmentId);

        modelBuilder.Entity<Employee>()
            .HasOne(e => e.Position)
            .WithMany()
            .HasForeignKey(e => e.PositionId);

        modelBuilder.Entity<Employee>()
            .HasOne(e => e.User)
            .WithMany()
            .HasForeignKey(e => e.UserId);

        // Tiếp tục cấu hình các quan hệ khác nếu có
    }
}
