use serde::{Deserialize, Serialize};
use wasmcloud_component::{
    error,
    http::{self, Method},
};

use crate::{
    bindings::yasp::llm::ollama,
    utils::{make_empty_response, make_response, IncomingRequestExt},
};

#[derive(Debug, Deserialize)]
struct MessageRequest {
    pub role: String,
    pub content: String,
}

#[derive(Debug, Deserialize)]
struct Request {
    pub model: String,
    pub messages: Vec<MessageRequest>,
}

#[derive(Debug, Serialize)]
struct MessageResponse {
    pub role: String,
    pub content: String,
    pub thinking: String,
    pub images: Vec<Vec<u8>>,
}

#[derive(Debug, Serialize)]
struct Metrics {
    pub total: u64,
    pub load: u64,
    pub prompt_eval_count: u32,
    pub prompt_eval_duration: u64,
    pub eval_count: u32,
    pub eval_duration: u64,
}

#[derive(Debug, Serialize)]
struct Response {
    pub model: String,
    pub create_at: String,
    pub message: MessageResponse,
    pub done_reason: String,
    pub metrics: Metrics,
}

pub fn chat(mut request: http::IncomingRequest) -> http::Result<http::Response<String>> {
    if request.method() != Method::POST {
        return make_empty_response(http::StatusCode::METHOD_NOT_ALLOWED);
    }

    if !request.has_json_content_type() {
        return make_empty_response(http::StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    let body = match request.read_body_json::<Request>() {
        Ok(value) => value,
        Err(err) => {
            error!("Failed to deserialize request body: {err}");
            return make_response(
                http::StatusCode::BAD_REQUEST,
                &serde_json::json!({
                    "error": err.to_string()
                }),
            );
        }
    };

    let response = ollama::chat(&ollama::ChatRequest {
        model: body.model,
        messages: body
            .messages
            .into_iter()
            .map(|m| ollama::ChatMessageRequest {
                role: m.role,
                content: m.content,
            })
            .collect(),
    });

    make_response(
        http::StatusCode::OK,
        &Response {
            model: response.model,
            create_at: response.create_at,
            message: MessageResponse {
                role: response.message.role,
                content: response.message.content,
                thinking: response.message.thinking,
                images: response.message.images,
            },
            done_reason: response.done_reason,
            metrics: Metrics {
                total: response.metrics.total,
                load: response.metrics.load,
                prompt_eval_count: response.metrics.prompt_eval_count,
                prompt_eval_duration: response.metrics.prompt_eval_duration,
                eval_count: response.metrics.eval_count,
                eval_duration: response.metrics.eval_duration,
            },
        },
    )
}
