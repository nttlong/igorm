using System;
using System.Security.Cryptography;
using System.Text;
using Isopoh.Cryptography.Argon2;

public class PasswordHasher
{
    private const int SaltLen = 16;
    private const int KeyLen = 32;
    private const int ArgonTime = 3;       // t
    private const int ArgonMemory = 64 * 1024; // KiB (64 MB)
    private const int ArgonThreads = 2;    // p

    public static string HashPassword(string password)
    {
        if (string.IsNullOrEmpty(password))
            throw new ArgumentException("empty password");

        // Tạo salt ngẫu nhiên
        byte[] salt = new byte[SaltLen];
        using (var rng = RandomNumberGenerator.Create())
        {
            rng.GetBytes(salt);
        }

        var config = new Argon2Config
        {
            Type = Argon2Type.HybridAddressing, //<-- 'Argon2Type' does not contain a definition for 'ID' [D:\code\go\news2\igorm\examples\dotnet\TestApi\TestApi.csproj]
            Version = Argon2Version.Nineteen,
            TimeCost = ArgonTime,
            MemoryCost = ArgonMemory,
            Lanes = ArgonThreads,
            Threads = ArgonThreads,
            Salt = salt,
            Password = Encoding.UTF8.GetBytes(password),
            HashLength = KeyLen
        };

        using (var argon2 = new Argon2(config))
        {
            var hash = argon2.Hash();
            byte[] hashBytes = hash.Buffer;
            // Format giống Go: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
            string encoded = $"$argon2id$v=19$m={ArgonMemory},t={ArgonTime},p={ArgonThreads}$" +
                             $"{Convert.ToBase64String(salt)}$" +
                             $"{Convert.ToBase64String(hashBytes)}"; //<--'SecureArray<byte
                                                                     //>' does not contain a definition for 'HashBytes' and no accessible extension method 'HashBytes' accepting a first ar
                                                                     //gument of type 'SecureArray<byte>' could be found (are you missing a using directive or an assembly reference?)
            return encoded;
        }
    }
}
