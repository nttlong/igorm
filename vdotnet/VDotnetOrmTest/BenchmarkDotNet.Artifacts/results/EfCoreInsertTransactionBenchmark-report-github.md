```

BenchmarkDotNet v0.15.2, Windows 11 (10.0.22631.5335/23H2/2023Update/SunValley3)
12th Gen Intel Core i7-12650H 2.30GHz, 1 CPU, 16 logical and 10 physical cores
.NET SDK 10.0.100-preview.6.25358.103
  [Host]     : .NET 10.0.0 (10.0.25.35903), X64 RyuJIT AVX2
  DefaultJob : .NET 10.0.0 (10.0.25.35903), X64 RyuJIT AVX2


```
| Method                               | Mean     | Error    | StdDev   | Gen0      | Gen1    | Gen2    | Allocated |
|------------------------------------- |---------:|---------:|---------:|----------:|--------:|--------:|----------:|
| InsertPositionDepartmentUserEmployee | 20.77 ms | 3.517 ms | 10.20 ms | 4867.1875 | 93.7500 | 23.4375 |  58.79 MB |
