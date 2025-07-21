public class Department : BaseModel
{
    public int Id { get; set; }
    public string Name { get; set; } = null!;
    public string Code { get; set; } = null!;
    public int? ParentId { get; set; }
    public Department? Parent { get; set; }
}
