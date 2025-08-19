using Microsoft.AspNetCore.Mvc;
//dotnet-counters monitor --process-id 10476 --counters System.Runtime
//dotnet-counters collect --process-id 10476 --counters System.Runtime --output myapp_cpu_trace.nettrace
namespace MyApi.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class MediaController : ControllerBase
    {
        private readonly string rootPath = @"D:\code\go\news2\igorm\examples\media\cmd\uploads";
        [HttpGet("Hello")]
        public IActionResult Hello()
        {
            return Ok("Hello World");
        }
        [HttpGet("list-files")]
        public IActionResult ListFiles()
        {
            if (!Directory.Exists(rootPath))
            {
                return NotFound(new { message = $"Directory '{rootPath}' not found." });
            }

            var baseUrl = $"{Request.Scheme}://{Request.Host}/api/media/files";
            var result = new List<string>();

            // Thư mục con
            foreach (var dir in Directory.GetDirectories(rootPath, "*", SearchOption.AllDirectories))
            {
                var relativePath = Path.GetRelativePath(rootPath, dir).Replace("\\", "/");
                result.Add($"{baseUrl}/{relativePath}");
            }

            // File
            foreach (var file in Directory.GetFiles(rootPath, "*", SearchOption.AllDirectories))
            {
                var relativePath = Path.GetRelativePath(rootPath, file).Replace("\\", "/");
                result.Add($"{baseUrl}/{relativePath}");
            }

            return Ok(result);
        }

        // API phục vụ file trực tiếp (download/view)
        [HttpGet("files/{*filePath}")]
        public IActionResult GetFile(string filePath)
        {
            var fullPath = Path.Combine(rootPath, filePath);
            if (!System.IO.File.Exists(fullPath))
            {
                return NotFound();
            }

            var contentType = "application/octet-stream";
            return PhysicalFile(fullPath, contentType);
        }
    }
}
