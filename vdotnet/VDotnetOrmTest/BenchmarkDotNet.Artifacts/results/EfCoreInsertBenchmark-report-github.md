```

BenchmarkDotNet v0.15.2, Windows 11 (10.0.22631.5335/23H2/2023Update/SunValley3)
12th Gen Intel Core i7-12650H 2.30GHz, 1 CPU, 16 logical and 10 physical cores
.NET SDK 10.0.100-preview.6.25358.103
  [Host]     : .NET 10.0.0 (10.0.25.35903), X64 RyuJIT AVX2
  DefaultJob : .NET 10.0.0 (10.0.25.35903), X64 RyuJIT AVX2


```
| Method           | Mean     | Error     | StdDev   | Gen0     | Gen1    | Gen2    | Allocated |
|----------------- |---------:|----------:|---------:|---------:|--------:|--------:|----------:|
| InsertSingleUser | 5.414 ms | 0.7935 ms | 2.340 ms | 933.5938 | 35.1563 | 11.7188 |  11.27 MB |
