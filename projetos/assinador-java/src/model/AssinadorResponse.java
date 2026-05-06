package br.gov.go.ses.assinador.model;

public class AssinadorResponse {

    private boolean success;
    private String data;
    private String errorCode;
    private String errorMessage;

    public AssinadorResponse() {
    }

    public AssinadorResponse(boolean success, String data, String errorCode, String errorMessage) {
        this.success = success;
        this.data = data;
        this.errorCode = errorCode;
        this.errorMessage = errorMessage;
    }

    public static AssinadorResponse ok(String data) {
        return new AssinadorResponse(true, data, null, null);
    }

    public static AssinadorResponse error(String errorCode, String errorMessage) {
        return new AssinadorResponse(false, null, errorCode, errorMessage);
    }

    public boolean isSuccess() {
        return success;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }

    public String getData() {
        return data;
    }

    public void setData(String data) {
        this.data = data;
    }

    public String getErrorCode() {
        return errorCode;
    }

    public void setErrorCode(String errorCode) {
        this.errorCode = errorCode;
    }

    public String getErrorMessage() {
        return errorMessage;
    }

    public void setErrorMessage(String errorMessage) {
        this.errorMessage = errorMessage;
    }
}
