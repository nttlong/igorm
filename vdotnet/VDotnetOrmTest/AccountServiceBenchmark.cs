using BenchmarkDotNet.Attributes;
using BenchmarkDotNet.Running;
using Microsoft.EntityFrameworkCore;
using System;
using System.Threading;
using System.Threading.Tasks;

[MemoryDiagnoser]
public class AccountServiceBenchmark
{
    private AccountService _accountService;
    private HrDbContext _dbContext;
    private BcryptPasswordService _passwordService;

    [GlobalSetup]
    public void Setup()
    {
        _dbContext = new HrDbContext();



        _passwordService = new BcryptPasswordService(); // Bạn đã viết interface này rồi
        _accountService = new AccountService(_dbContext, _passwordService);
    }

    // [Benchmark]
    // public async Task CreateOrUpdateAsync_NewAccount()
    // {

    //     var acc = new Account
    //     {
    //         Username = "existing_user",
    //         HashedPassword = _passwordService.Hash("password"), // ✅ Sử dụng đúng passwordService
    //         FullName = "Benchmark User",
    //         Role = "user",
    //         CreatedAt = DateTime.UtcNow
    //     };

    //     // ✅ Thêm vào DbContext
    //     _dbContext.Accounts.Add(acc);
    //     await _dbContext.SaveChangesAsync(); // Đảm bảo account đã tồn tại trước khi test update

    //     await _accountService.CreateOrUpdateAsync(acc);
    // }

    [Benchmark]
    public async Task CreateOrUpdateAsync_ExistingAccount()
    {
        var acc = new Account
        {
            Username = "existing_user",
            HashedPassword = _passwordService.Hash("password"), //<-- sua cho nay dung passwordService
            FullName = "Benchmark User",
            Role = "user",
            CreatedAt = DateTime.UtcNow
        };

        // // Insert trước để tồn tại
        // _dbContext.Accounts.Add(acc);
        // await _dbContext.SaveChangesAsync();

        acc.FullName = "Updated Name";
        await _accountService.CreateOrUpdateAsync(acc, CancellationToken.None);
    }
}
