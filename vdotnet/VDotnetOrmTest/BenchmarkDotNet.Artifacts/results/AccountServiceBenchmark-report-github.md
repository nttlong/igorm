```

BenchmarkDotNet v0.13.8, Windows 11 (10.0.22631.5335)
12th Gen Intel Core i7-12650H, 1 CPU, 16 logical and 10 physical cores
.NET SDK 10.0.100-preview.6.25358.103
  [Host]     : .NET 8.0.18 (8.0.1825.31117), X64 RyuJIT AVX2
  DefaultJob : .NET 8.0.18 (8.0.1825.31117), X64 RyuJIT AVX2


```
| Method                              | Mean     | Error   | StdDev  | Allocated |
|------------------------------------ |---------:|--------:|--------:|----------:|
| CreateOrUpdateAsync_ExistingAccount | 216.8 ms | 4.29 ms | 4.60 ms |  34.91 KB |
