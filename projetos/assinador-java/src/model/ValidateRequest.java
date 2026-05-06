package br.gov.go.ses.assinador.model;

public class ValidateRequest {

    private String signatureData;
    private Long referenceTimestamp;
    private String policyUri;

    public ValidateRequest() {
    }

    public String getSignatureData() {
        return signatureData;
    }

    public void setSignatureData(String signatureData) {
        this.signatureData = signatureData;
    }

    public Long getReferenceTimestamp() {
        return referenceTimestamp;
    }

    public void setReferenceTimestamp(Long referenceTimestamp) {
        this.referenceTimestamp = referenceTimestamp;
    }

    public String getPolicyUri() {
        return policyUri;
    }

    public void setPolicyUri(String policyUri) {
        this.policyUri = policyUri;
    }
}
