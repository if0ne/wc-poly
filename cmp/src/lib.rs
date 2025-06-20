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

use std::io::Read;

use bindings::yasp::giga::wc;
use wasmcloud_component::http;

struct Component;

impl http::Server for Component {
    fn handle(
        mut request: http::IncomingRequest,
    ) -> http::Result<http::Response<impl http::OutgoingBody>> {
        let mut arg = String::new();
        let _ = request.body_mut().read_to_string(&mut arg);
        let giga_wc = wc::whoami(&arg);

        Ok(http::Response::new(giga_wc))
    }
}
