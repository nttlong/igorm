using System;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

[Table("accounts")]
public class Account
{
    [Key]
    [Column("id")]
    [DatabaseGenerated(DatabaseGeneratedOption.Identity)] // <- Tự động tăng

    public int ID { get; set; }

    [Column("user_id")]
    [MaxLength(36)]
    public string UserID { get; set; } = string.Empty;

    [Column("username")]
    [MaxLength(50)]
    public string Username { get; set; } = string.Empty;

    [Column("hashed_password")]
    [MaxLength(255)]
    public string HashedPassword { get; set; } = string.Empty;

    [Column("is_locked")]
    public bool IsLocked { get; set; }

    [Column("locked_until")]
    public DateTime? LockedUntil { get; set; }

    [Column("full_name")]
    [MaxLength(255)]
    public string FullName { get; set; } = string.Empty;

    [Column("email")]
    [MaxLength(255)]
    public string? Email { get; set; }

    [Column("phone")]
    [MaxLength(20)]
    public string? Phone { get; set; }

    [Column("role")]
    [MaxLength(50)]
    public string Role { get; set; } = "user";

    [Column("created_at", TypeName = "datetime2(3)")]
    public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

    [Column("updated_at", TypeName = "datetime2(3)")]
    public DateTime? UpdatedAt { get; set; }
}
