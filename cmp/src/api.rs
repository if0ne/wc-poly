use wasmcloud_component::{
    error,
    http::{self, Method},
};

use crate::{
    bindings::yasp::llm::ollama,
    utils::{make_empty_response, make_response, make_response_from_string, IncomingRequestExt},
};

pub fn chat(mut request: http::IncomingRequest) -> http::Result<http::Response<String>> {
    if request.method() != Method::POST {
        return make_empty_response(http::StatusCode::METHOD_NOT_ALLOWED);
    }

    if !request.has_json_content_type() {
        return make_empty_response(http::StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    let body = match request.read_body_string() {
        Ok(value) => value,
        Err(err) => {
            error!("Failed to read request body: {err}");
            return make_response(
                http::StatusCode::BAD_REQUEST,
                &serde_json::json!({
                    "error": err.to_string()
                }),
            );
        }
    };

    let response = ollama::chat(&body);

    match response {
        Ok(response) => make_response_from_string(http::StatusCode::OK, response),
        Err(error) => make_response_from_string(http::StatusCode::BAD_REQUEST, error),
    }
}

pub fn pull(mut request: http::IncomingRequest) -> http::Result<http::Response<String>> {
    if request.method() != Method::POST {
        return make_empty_response(http::StatusCode::METHOD_NOT_ALLOWED);
    }

    if !request.has_json_content_type() {
        return make_empty_response(http::StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    let body = match request.read_body_string() {
        Ok(value) => value,
        Err(err) => {
            error!("Failed to read request body: {err}");
            return make_response(
                http::StatusCode::BAD_REQUEST,
                &serde_json::json!({
                    "error": err.to_string()
                }),
            );
        }
    };

    let response = ollama::pull(&body);

    match response {
        Ok(response) => make_response_from_string(http::StatusCode::OK, response),
        Err(error) => make_response_from_string(http::StatusCode::BAD_REQUEST, error),
    }
}

pub fn list(request: http::IncomingRequest) -> http::Result<http::Response<String>> {
    if request.method() != Method::GET {
        return make_empty_response(http::StatusCode::METHOD_NOT_ALLOWED);
    }

    let response = ollama::model_list();

    match response {
        Ok(response) => make_response_from_string(http::StatusCode::OK, response),
        Err(error) => make_response_from_string(http::StatusCode::BAD_REQUEST, error),
    }
}

pub fn delete(mut request: http::IncomingRequest) -> http::Result<http::Response<String>> {
    if request.method() != Method::DELETE {
        return make_empty_response(http::StatusCode::METHOD_NOT_ALLOWED);
    }

    if !request.has_json_content_type() {
        return make_empty_response(http::StatusCode::UNSUPPORTED_MEDIA_TYPE);
    }

    let body = match request.read_body_string() {
        Ok(value) => value,
        Err(err) => {
            error!("Failed to read request body: {err}");
            return make_response(
                http::StatusCode::BAD_REQUEST,
                &serde_json::json!({
                    "error": err.to_string()
                }),
            );
        }
    };

    let response = ollama::delete(&body);

    match response {
        Ok(response) => make_response_from_string(http::StatusCode::OK, response),
        Err(error) => make_response_from_string(http::StatusCode::BAD_REQUEST, error),
    }
}
