using Microsoft.OpenApi.Models;

var builder = WebApplication.CreateBuilder(args); //<--The name 'WebApplication' does not exist in the current context

// Thêm Swagger service
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1", new OpenApiInfo
    {
        Title = "My API",
        Version = "v1",
        Description = "Ví dụ Web API với Swagger trong .NET Core"
    });
});

var app = builder.Build();

// Bật Swagger UI trong môi trường Development
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI(c =>
    {
        c.SwaggerEndpoint("/swagger/v1/swagger.json", "My API V1");
        c.RoutePrefix = string.Empty; // Swagger sẽ hiển thị ở URL gốc "/"
    });
}

app.MapControllers();

app.Run();
