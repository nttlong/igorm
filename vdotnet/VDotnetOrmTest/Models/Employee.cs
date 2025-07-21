public class Employee : BaseModel
{
    public int Id { get; set; }
    public string FirstName { get; set; } = null!;
    public string LastName { get; set; } = null!;
    public int DepartmentId { get; set; }
    public Department Department { get; set; } = null!;
    public int PositionId { get; set; }
    public Position Position { get; set; } = null!;
    public int UserId { get; set; }
    public User User { get; set; } = null!;
}
