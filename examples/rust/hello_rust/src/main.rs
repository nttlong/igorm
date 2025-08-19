use std::pin::Pin;
use std::future::Future;
use actix_web::{get, App, HttpServer, Responder, HttpResponse, HttpRequest};
use serde_json::json;
use std::path::{PathBuf, Path};
use tokio::fs;
use tokio::io;
//use tokio::task;
use tokio::sync::Mutex;

static ROOT_PATH: &str = "D:\\code\\go\\news2\\igorm\\examples\\media\\cmd\\uploads";

#[get("/api/media/hello")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello World")
}

fn list_dir_recursive(
    dir: PathBuf,
    base_url: String,
    results: &Mutex<Vec<String>>,
) -> Pin<Box<dyn Future<Output = io::Result<()>> + Send + '_>> {
    Box::pin(async move {
        let mut read_dir = fs::read_dir(&dir).await?;
        while let Some(entry) = read_dir.next_entry().await? {
            let path = entry.path();
            if path.is_dir() {
                // gọi đệ quy boxed
                list_dir_recursive(path, base_url.clone(), results).await?;
            } else {
                if let Ok(rel) = path.strip_prefix(ROOT_PATH) {
                    let url_path = rel.to_string_lossy().replace("\\", "/");
                    results.lock().await.push(format!("{}/{}", base_url, url_path));
                }
            }
        }
        Ok(())
    })
}

#[get("/api/media/list-files")]
async fn list_files(req: HttpRequest) -> impl Responder {
    if !Path::new(ROOT_PATH).exists() {
        return HttpResponse::NotFound().body("Directory not found.");
    }

    let base_url = format!("http://{}{}", req.connection_info().host(), req.path());
    let results = tokio::sync::Mutex::new(Vec::with_capacity(1024));

    /*
    spawn task để tránh blocking
    */
    if let Err(e) = list_dir_recursive(PathBuf::from(ROOT_PATH), base_url, &results).await {
        eprintln!("Error reading directory: {}", e);
        
        return HttpResponse::InternalServerError().body("Internal server error.");
    }

    let locked = results.lock().await;
    HttpResponse::Ok()
        .content_type("application/json")
        .body(json!(&*locked).to_string())
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Server running at http://127.0.0.1:8080");

    HttpServer::new(|| {
        App::new()
            
            .service(list_files)
    })
    .bind("0.0.0.0:8082")?
    .run()
    .await
}
