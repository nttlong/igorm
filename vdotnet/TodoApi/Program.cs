
var builder = WebApplication.CreateBuilder(args);
builder.Services.AddControllers();

builder.Services.AddControllers();
var app = builder.Build();

app.MapControllers();
app.Run();
