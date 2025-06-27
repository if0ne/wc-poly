mod bindings {
    use super::Component;

    wit_bindgen::generate!({ generate_all });

    export!(Component);
}

use crate::bindings::{
    exports::wasmcloud::messaging::handler::Guest,
    wasmcloud::messaging::{
        consumer,
        types::{self, BrokerMessage},
    },
    yasp::llm::ollama,
};

struct Component;

impl Guest for Component {
    fn handle_message(msg: BrokerMessage) -> Result<(), String> {
        let result = match msg.subject.split('.').skip(3).next().unwrap() {
            "chat" => {
                let body = String::from_utf8_lossy(&msg.body);
                ollama::chat(&body)
            }
            "list" => ollama::model_list(),
            "pull" => {
                let body = String::from_utf8_lossy(&msg.body);
                ollama::pull(&body)
            }
            "delete" => {
                let body = String::from_utf8_lossy(&msg.body);
                ollama::delete(&body)
            }
            _ => Err("unknown api".to_string()),
        };

        if let Some(reply) = msg.reply_to {
            let body = match result {
                Ok(body) => body.into_bytes(),
                Err(err) => serde_json::to_vec(&serde_json::json!({
                    "error": err
                }))
                .unwrap(),
            };

            let _ = consumer::publish(&types::BrokerMessage {
                subject: reply,
                reply_to: None,
                body,
            });
        }

        Ok(())
    }
}
