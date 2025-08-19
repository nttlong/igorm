using Microsoft.AspNetCore.Mvc;

namespace MyApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class HelloController : ControllerBase
    {
        [HttpGet("say")]
        public IActionResult SayHello()
        {
            return Ok(new { message = "Hello from Swagger!" });
        }
    }
}
