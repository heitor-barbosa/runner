package br.gov.go.ses.assinador.http;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.model.SignRequest;
import br.gov.go.ses.assinador.model.ValidateRequest;
import br.gov.go.ses.assinador.service.SignatureService;
import br.gov.go.ses.assinador.validation.SignRequestValidator;
import br.gov.go.ses.assinador.validation.ValidateRequestValidator;
import br.gov.go.ses.assinador.validation.ValidationError;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpExchange;

import java.io.IOException;
import java.io.OutputStream;
import java.nio.charset.StandardCharsets;
import java.util.List;

public class SignatureController {

    private static final String CONTENT_TYPE_JSON = "application/json; charset=utf-8";

    private final ObjectMapper mapper;
    private final SignatureService service;
    private final SignRequestValidator signValidator;
    private final ValidateRequestValidator validateValidator;

    public SignatureController(
            SignatureService service,
            SignRequestValidator signValidator,
            ValidateRequestValidator validateValidator
    ) {
        this.mapper = new ObjectMapper();
        this.service = service;
        this.signValidator = signValidator;
        this.validateValidator = validateValidator;
    }

    public void handleSign(HttpExchange exchange) throws IOException {
        if (!requirePost(exchange)) {
            return;
        }

        SignRequest request;
        try {
            request = mapper.readValue(readBody(exchange), SignRequest.class);
        } catch (Exception error) {
            writeResponse(exchange, 400, AssinadorResponse.error(
                    "FORMAT.BUNDLE-MALFORMED",
                    "JSON de entrada invalido: " + error.getMessage()
            ));
            return;
        }

        List<ValidationError> errors = signValidator.validate(request);
        if (!errors.isEmpty()) {
            writeResponse(exchange, 400, validationErrorResponse(errors));
            return;
        }

        AssinadorResponse response = service.sign(request);
        writeResponse(exchange, response.isSuccess() ? 200 : 400, response);
    }

    public void handleValidate(HttpExchange exchange) throws IOException {
        if (!requirePost(exchange)) {
            return;
        }

        ValidateRequest request;
        try {
            request = mapper.readValue(readBody(exchange), ValidateRequest.class);
        } catch (Exception error) {
            writeResponse(exchange, 400, AssinadorResponse.error(
                    "FORMAT.JWS-MALFORMED",
                    "JSON de entrada invalido: " + error.getMessage()
            ));
            return;
        }

        List<ValidationError> errors = validateValidator.validate(request);
        if (!errors.isEmpty()) {
            writeResponse(exchange, 400, validationErrorResponse(errors));
            return;
        }

        AssinadorResponse response = service.validate(request);
        writeResponse(exchange, response.isSuccess() ? 200 : 400, response);
    }

    private boolean requirePost(HttpExchange exchange) throws IOException {
        if ("POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            return true;
        }

        exchange.getResponseHeaders().add("Allow", "POST");
        writeResponse(exchange, 405, AssinadorResponse.error(
                "HTTP.METHOD-NOT-ALLOWED",
                "Metodo nao permitido. Use POST."
        ));
        return false;
    }

    private AssinadorResponse validationErrorResponse(List<ValidationError> errors) {
        ValidationError first = errors.get(0);
        String details = errors.size() > 1
                ? first.getMessage() + " (e mais " + (errors.size() - 1) + " problema(s))"
                : first.getMessage();
        return AssinadorResponse.error(first.getCode(), details);
    }

    private String readBody(HttpExchange exchange) throws IOException {
        return new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
    }

    private void writeResponse(HttpExchange exchange, int statusCode, AssinadorResponse response) throws IOException {
        byte[] payload = mapper.writeValueAsBytes(response);
        exchange.getResponseHeaders().set("Content-Type", CONTENT_TYPE_JSON);
        exchange.sendResponseHeaders(statusCode, payload.length);
        try (OutputStream output = exchange.getResponseBody()) {
            output.write(payload);
        }
    }
}
