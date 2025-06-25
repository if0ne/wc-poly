mod bindings {
    use super::Component;

    wit_bindgen::generate!({
        with: {
            "wasi:http/incoming-handler@0.2.2": wasmcloud_component::wasi::exports::http::incoming_handler,
        },
        generate_all
    });

    wasmcloud_component::http::export!(Component);
}
mod utils;

use bindings::yasp::giga::wc;
use serde::{Deserialize, Serialize};
use wasmcloud_component::{error, http};

use crate::{
    bindings::yasp,
    utils::{make_response, IncomingRequestExt},
};

struct Component;

#[derive(Debug, Deserialize)]
struct ChatMessageRequest {
    pub role: String,
    pub content: String,
}

#[derive(Debug, Deserialize)]
struct ChatRequest {
    pub model: String,
    pub messages: Vec<ChatMessageRequest>,
}

#[derive(Debug, Serialize)]
struct ChatMessageResponse {
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
struct ChatResponse {
    pub model: String,
    pub create_at: String,
    pub message: ChatMessageResponse,
    pub done_reason: String,
    pub metrics: Metrics,
}

impl http::Server for Component {
    fn handle(
        mut request: http::IncomingRequest,
    ) -> http::Result<http::Response<impl http::OutgoingBody>> {
        let body = match request.read_body_json::<ChatRequest>() {
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
        let response = wc::chat(&yasp::giga::wc::ChatRequest {
            model: body.model,
            messages: body
                .messages
                .into_iter()
                .map(|m| yasp::giga::wc::ChatMessageRequest {
                    role: m.role,
                    content: m.content,
                })
                .collect(),
        });

        make_response(
            http::StatusCode::OK,
            &ChatResponse {
                model: response.model,
                create_at: response.create_at,
                message: ChatMessageResponse {
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
}
