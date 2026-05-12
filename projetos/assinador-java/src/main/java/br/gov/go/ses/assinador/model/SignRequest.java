package br.gov.go.ses.assinador.model;

public class SignRequest {

    private String bundle;
    private String provenance;
    private String credentialType;
    private String credentialContent;
    private String credentialPassword;
    private String credentialAlias;
    private String certificateChain;
    private Long referenceTimestamp;
    private String strategy;
    private String policyUri;

    public String getBundle() {
        return bundle;
    }

    public void setBundle(String bundle) {
        this.bundle = bundle;
    }

    public String getProvenance() {
        return provenance;
    }

    public void setProvenance(String provenance) {
        this.provenance = provenance;
    }

    public String getCredentialType() {
        return credentialType;
    }

    public void setCredentialType(String credentialType) {
        this.credentialType = credentialType;
    }

    public String getCredentialContent() {
        return credentialContent;
    }

    public void setCredentialContent(String credentialContent) {
        this.credentialContent = credentialContent;
    }

    public String getCredentialPassword() {
        return credentialPassword;
    }

    public void setCredentialPassword(String credentialPassword) {
        this.credentialPassword = credentialPassword;
    }

    public String getCredentialAlias() {
        return credentialAlias;
    }

    public void setCredentialAlias(String credentialAlias) {
        this.credentialAlias = credentialAlias;
    }

    public String getCertificateChain() {
        return certificateChain;
    }

    public void setCertificateChain(String certificateChain) {
        this.certificateChain = certificateChain;
    }

    public Long getReferenceTimestamp() {
        return referenceTimestamp;
    }

    public void setReferenceTimestamp(Long referenceTimestamp) {
        this.referenceTimestamp = referenceTimestamp;
    }

    public String getStrategy() {
        return strategy;
    }

    public void setStrategy(String strategy) {
        this.strategy = strategy;
    }

    public String getPolicyUri() {
        return policyUri;
    }

    public void setPolicyUri(String policyUri) {
        this.policyUri = policyUri;
    }
}
