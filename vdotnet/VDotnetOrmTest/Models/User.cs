public class User : BaseModel
{
    public int Id { get; set; }
    public string? UserId { get; set; }
    public string Email { get; set; } = null!;
    public string Phone { get; set; } = null!;
    public string? Username { get; set; }
    public string? HashPassword { get; set; }
}
