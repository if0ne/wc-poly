use std::io::Read;

use serde::{de::DeserializeOwned, Serialize};
use wasmcloud_component::http;

pub trait IncomingRequestExt {
    fn has_json_content_type(&self) -> bool;

    fn read_body_json<T: DeserializeOwned>(&mut self) -> anyhow::Result<T>;
}

impl IncomingRequestExt for http::IncomingRequest {
    fn has_json_content_type(&self) -> bool {
        self.headers()
            .get(http::header::CONTENT_TYPE)
            .is_some_and(|h| h.to_str().is_ok_and(|h| h.starts_with("application/json")))
    }

    fn read_body_json<T: DeserializeOwned>(&mut self) -> anyhow::Result<T> {
        let content_length = self
            .headers()
            .get(http::header::CONTENT_LENGTH)
            .and_then(|h| h.to_str().ok())
            .and_then(|s| s.parse::<usize>().ok())
            .unwrap_or(64);

        let body = self.body_mut();

        let mut buf = Vec::with_capacity(content_length);
        body.read_to_end(&mut buf)
            .map_err(|e| anyhow::anyhow!("failed to read request body: {e}"))?;

        serde_json::from_slice::<T>(&buf)
            .map_err(|e| anyhow::anyhow!("failed to decode request body: {e}"))
    }
}

pub fn make_response<T: Serialize + ?Sized>(
    code: http::StatusCode,
    body: &T,
) -> http::Result<http::Response<String>> {
    let body = serde_json::to_string(body)
        .map_err(|e| http::ErrorCode::InternalError(Some(format!("failed to deserialize: {e}"))))?;

    http::Response::builder()
        .status(code)
        .header(
            http::header::CONTENT_TYPE,
            "application/json; charset=utf-8",
        )
        .body(body)
        .map_err(|e| http::ErrorCode::InternalError(Some(format!("failed to build response: {e}"))))
}

pub fn make_empty_response(code: http::StatusCode) -> http::Result<http::Response<String>> {
    http::Response::builder()
        .status(code)
        .body(String::new())
        .map_err(|e| http::ErrorCode::InternalError(Some(format!("failed to build response: {e}"))))
}
