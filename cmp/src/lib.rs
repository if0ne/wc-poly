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
mod api;
mod utils;

use wasmcloud_component::http;

use crate::utils::make_empty_response;

struct Component;

impl http::Server for Component {
    fn handle(
        request: http::IncomingRequest,
    ) -> http::Result<http::Response<impl http::OutgoingBody>> {
        match request.uri().path() {
            "/api/chat" => api::chat(request),
            "/api/list" => api::list(request),
            "/api/pull" => api::pull(request),
            "/api/delete" => api::delete(request),
            _ => make_empty_response(http::StatusCode::NOT_FOUND),
        }
    }
}
