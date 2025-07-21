public class Position : BaseModel
{
    public int Id { get; set; }
    public string Code { get; set; } = null!;
    public string Name { get; set; } = null!;
    public string Title { get; set; } = null!;
    public int Level { get; set; }
}
