using Microsoft.AspNetCore.Mvc;

[ApiController]
[Route("api/[controller]")]
public class TodoController : ControllerBase
{
    private static int idCounter = 0;
    private static readonly List<TodoItem> todos = new();

    [HttpPost]
    public IActionResult CreateTodo([FromBody] TodoItem item)
    {
        item.Id = Interlocked.Increment(ref idCounter);
        todos.Add(item);
        return Created("", item);
    }
}

public class TodoItem
{
    public int Id { get; set; }
    public string Name { get; set; }
    public bool IsComplete { get; set; }
}
