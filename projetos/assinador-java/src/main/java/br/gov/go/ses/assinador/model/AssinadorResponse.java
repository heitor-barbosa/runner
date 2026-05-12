package br.gov.go.ses.assinador.model;

public class AssinadorResponse {

    private boolean success;
    private String data;
    private String errorCode;
    private String errorMessage;

    public static AssinadorResponse ok(String data) {
        AssinadorResponse response = new AssinadorResponse();
        response.success = true;
        response.data = data;
        return response;
    }

    public static AssinadorResponse error(String code, String message) {
        AssinadorResponse response = new AssinadorResponse();
        response.success = false;
        response.errorCode = code;
        response.errorMessage = message;
        return response;
    }

    public boolean isSuccess() {
        return success;
    }

    public String getData() {
        return data;
    }

    public String getErrorCode() {
        return errorCode;
    }

    public String getErrorMessage() {
        return errorMessage;
    }
}
