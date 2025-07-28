using BenchmarkDotNet.Running;

public class Program
{
    public static void Main(string[] args)
    {
        //BenchmarkRunner.Run<EfCoreInsertBenchmark>();
        //BenchmarkRunner.Run<EfCoreInsertTransactionBenchmark>();
        //BenchmarkRunner.Run<EfCoreRawSqlBenchmark>();
        BenchmarkRunner.Run<AccountServiceBenchmark>();
        //dotnet run -c Release

    }
}
