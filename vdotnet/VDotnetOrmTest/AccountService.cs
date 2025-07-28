using System;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;


public class AccountService
{
    private readonly HrDbContext _dbContext;
    private readonly IPasswordService _passwordService;

    public AccountService(HrDbContext dbContext, IPasswordService passwordService)
    {
        _dbContext = dbContext;
        _passwordService = passwordService;
    }

    public async Task CreateOrUpdateAsync(Account account, CancellationToken cancellationToken = default)
    {
        var existing = await _dbContext.Accounts
            .FirstOrDefaultAsync(a => a.Username == account.Username, cancellationToken);

        if (existing == null)
        {
            Console.WriteLine("Creating new account");
            _dbContext.Accounts.Add(account);
        }
        else
        {
            // KHÔNG cập nhật Username vì nó dùng để tìm
            existing.HashedPassword = account.HashedPassword;
            existing.IsLocked = account.IsLocked;
            existing.LockedUntil = account.LockedUntil;
            existing.FullName = account.FullName;
            existing.Email = account.Email;
            existing.Phone = account.Phone;
            existing.Role = account.Role;
            existing.UpdatedAt = DateTime.UtcNow;

            // KHÔNG cần set lại ID hoặc trạng thái Entity
        }

        await _dbContext.SaveChangesAsync(cancellationToken);
    }



}
