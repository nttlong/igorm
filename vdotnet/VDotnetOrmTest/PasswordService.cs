using BCrypt.Net;
public interface IPasswordService
{
    string Hash(string password);
    bool Verify(string password, string hashed);
}


public class BcryptPasswordService : IPasswordService
{
    private readonly int _workFactor;

    public BcryptPasswordService(int workFactor = 12)
    {
        _workFactor = workFactor > 0 ? workFactor : 12;
    }

    public string Hash(string password)
    {
        return BCrypt.Net.BCrypt.HashPassword(password, _workFactor);
    }

    public bool Verify(string password, string hashed)
    {
        return BCrypt.Net.BCrypt.Verify(password, hashed);
    }
}
