package br.gov.go.ses.assinador.cli;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.model.SignRequest;
import br.gov.go.ses.assinador.model.ValidateRequest;
import br.gov.go.ses.assinador.http.AssinadorHttpServer;
import br.gov.go.ses.assinador.http.HttpServerMain;
import br.gov.go.ses.assinador.service.FakeSignatureService;
import br.gov.go.ses.assinador.service.SignatureService;
import br.gov.go.ses.assinador.validation.SignRequestValidator;
import br.gov.go.ses.assinador.validation.ValidateRequestValidator;
import br.gov.go.ses.assinador.validation.ValidationError;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.util.List;

public class Main {

    private static final ObjectMapper MAPPER = new ObjectMapper();
    private static final SignatureService SERVICE = new FakeSignatureService();
    private static final SignRequestValidator SIGN_VALIDATOR = new SignRequestValidator();
    private static final ValidateRequestValidator VALIDATE_VALIDATOR = new ValidateRequestValidator();

    public static void main(String[] args) {
        if (args.length < 2) {
            printUsage();
            System.exit(1);
        }

        String command = args[0].toLowerCase();

        if ("server".equals(command)) {
            new HttpServerMain().startAndBlock(extractPort(args));
            return;
        }

        String json = extractJson(args);

        if (json == null || json.isBlank()) {
            printError("CONFIG.MISSING-PARAMETER", "Parametro --json e obrigatorio.");
            System.exit(1);
        }

        switch (command) {
            case "sign" -> handleSign(json);
            case "validate" -> handleValidate(json);
            default -> {
                printError(
                        "CONFIG.INVALID-PARAMETER",
                "Comando desconhecido: '" + command + "'. Use 'sign', 'validate' ou 'server'."
                );
                System.exit(1);
            }
        }
    }

    private static void handleSign(String json) {
        SignRequest request;
        try {
            request = MAPPER.readValue(json, SignRequest.class);
        } catch (Exception error) {
            printError("FORMAT.BUNDLE-MALFORMED", "JSON de entrada invalido: " + error.getMessage());
            System.exit(1);
            return;
        }

        List<ValidationError> errors = SIGN_VALIDATOR.validate(request);
        if (!errors.isEmpty()) {
            printValidationErrors(errors);
            System.exit(1);
            return;
        }

        AssinadorResponse response = SERVICE.sign(request);
        printResponse(response);
        System.exit(response.isSuccess() ? 0 : 1);
    }

    private static void handleValidate(String json) {
        ValidateRequest request;
        try {
            request = MAPPER.readValue(json, ValidateRequest.class);
        } catch (Exception error) {
            printError("FORMAT.JWS-MALFORMED", "JSON de entrada invalido: " + error.getMessage());
            System.exit(1);
            return;
        }

        List<ValidationError> errors = VALIDATE_VALIDATOR.validate(request);
        if (!errors.isEmpty()) {
            printValidationErrors(errors);
            System.exit(1);
            return;
        }

        AssinadorResponse response = SERVICE.validate(request);
        printResponse(response);
        System.exit(response.isSuccess() ? 0 : 1);
    }

    private static void printResponse(AssinadorResponse response) {
        try {
            System.out.println(MAPPER.writeValueAsString(response));
        } catch (Exception error) {
            System.err.println("Erro ao serializar resposta: " + error.getMessage());
        }
    }

    private static void printError(String code, String message) {
        printResponse(AssinadorResponse.error(code, message));
    }

    private static void printValidationErrors(List<ValidationError> errors) {
        ValidationError first = errors.get(0);
        String details = errors.size() > 1
                ? first.getMessage() + " (e mais " + (errors.size() - 1) + " problema(s))"
                : first.getMessage();
        printResponse(AssinadorResponse.error(first.getCode(), details));
    }

    private static void printUsage() {
        System.err.println("""
                Assinador HubSaude v0.1.0 - Sistema Runner SES-GO/UFG

                Uso:
                  java -jar assinador.jar sign     --json '<SignRequest JSON>'
                  java -jar assinador.jar validate --json '<ValidateRequest JSON>'
                  java -jar assinador.jar server   --port 8080
                """);
    }

    private static int extractPort(String[] args) {
        for (int index = 1; index < args.length - 1; index++) {
            if ("--port".equals(args[index])) {
                try {
                    int port = Integer.parseInt(args[index + 1]);
                    if (port < 0 || port > 65535) {
                        throw new NumberFormatException("porta fora do intervalo");
                    }
                    return port;
                } catch (NumberFormatException error) {
                    printError("CONFIG.INVALID-PARAMETER", "Porta invalida: " + args[index + 1]);
                    System.exit(1);
                }
            }
        }
        return AssinadorHttpServer.DEFAULT_PORT;
    }

    private static String extractJson(String[] args) {
        for (int index = 1; index < args.length - 1; index++) {
            if ("--json".equals(args[index])) {
                return args[index + 1];
            }
        }
        if (args.length >= 2 && !args[1].startsWith("--")) {
            return args[1];
        }
        return null;
    }
}
