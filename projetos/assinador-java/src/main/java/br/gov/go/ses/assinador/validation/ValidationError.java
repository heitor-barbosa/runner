package br.gov.go.ses.assinador.validation;

public class ValidationError {

    private final String code;
    private final String message;

    public ValidationError(String code, String message) {
        this.code = code;
        this.message = message;
    }

    public String getCode() {
        return code;
    }

    public String getMessage() {
        return message;
    }

    @Override
    public String toString() {
        return "[" + code + "] " + message;
    }
}
