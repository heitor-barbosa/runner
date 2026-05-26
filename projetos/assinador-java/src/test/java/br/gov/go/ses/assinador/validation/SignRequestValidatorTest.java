package br.gov.go.ses.assinador.validation;

import br.gov.go.ses.assinador.model.SignRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class SignRequestValidatorTest {

    private static final long NOW = System.currentTimeMillis() / 1000L;
    private SignRequestValidator validator;

    @BeforeEach
    void setUp() {
        validator = new SignRequestValidator();
    }

    @Test
    void validRequestShouldNotGenerateErrors() {
        assertTrue(validator.validate(buildValid()).isEmpty());
    }

    @Test
    void missingPolicyShouldGenerateError() {
        SignRequest request = buildValid();
        request.setPolicyUri(null);
        assertHasCode(validator.validate(request), "POLICY.MISSING");
    }

    @Test
    void invalidTimestampShouldGenerateError() {
        SignRequest request = buildValid();
        request.setReferenceTimestamp(1_000_000_000L);
        assertHasCode(validator.validate(request), "CONFIG.TIMESTAMP-OUT-OF-RANGE");
    }

    @Test
    void invalidStrategyShouldGenerateError() {
        SignRequest request = buildValid();
        request.setStrategy("rsa");
        assertHasCode(validator.validate(request), "CONFIG.INVALID-STRATEGY");
    }

    @Test
    void malformedBundleShouldGenerateError() {
        SignRequest request = buildValid();
        request.setBundle("{\"entry\":[{}]}");
        assertHasCode(validator.validate(request), "FORMAT.BUNDLE-MALFORMED");
    }

    @Test
    void malformedProvenanceShouldGenerateError() {
        SignRequest request = buildValid();
        request.setProvenance("{\"resourceType\":\"Provenance\"}");
        assertHasCode(validator.validate(request), "FORMAT.PROVENANCE-INVALID");
    }

    @Test
    void invalidCredentialTypeShouldGenerateError() {
        SignRequest request = buildValid();
        request.setCredentialType("SSH");
        assertHasCode(validator.validate(request), "CONFIG.INVALID-PARAMETER");
    }

    @Test
    void incompleteCertificateChainShouldGenerateError() {
        SignRequest request = buildValid();
        request.setCertificateChain("[\"CERT1\"]");
        assertHasCode(validator.validate(request), "CERT.CHAIN-INCOMPLETE");
    }

    @Test
    void pkcs12ShouldRequirePasswordAndAlias() {
        SignRequest request = buildValid();
        request.setCredentialType("PKCS12");
        request.setCredentialPassword(null);
        request.setCredentialAlias(null);

        List<ValidationError> errors = validator.validate(request);
        assertFalse(errors.isEmpty());
        assertHasCode(errors, "CONFIG.MISSING-PARAMETER");
    }

    @Test
    void tokenShouldRequirePkcs11ConfigAndAlias() {
        SignRequest request = buildValid();
        request.setCredentialType("TOKEN");
        request.setPkcs11ConfigPath(null);
        request.setCredentialAlias(null);

        List<ValidationError> errors = validator.validate(request);

        assertFalse(errors.isEmpty());
        assertHasCode(errors, "PKCS11.CONFIG-MISSING");
        assertHasCode(errors, "PKCS11.ALIAS-MISSING");
    }

    @Test
    void smartCardShouldAcceptPkcs11Metadata() {
        SignRequest request = buildValid();
        request.setCredentialType("SMARTCARD");
        request.setPkcs11ConfigPath("pkcs11.cfg");
        request.setCredentialAlias("assinatura");
        request.setTokenLabel("token-a");

        assertTrue(validator.validate(request).isEmpty());
    }

    private SignRequest buildValid() {
        SignRequest request = new SignRequest();
        request.setBundle("{\"resourceType\":\"Bundle\",\"entry\":[{}]}");
        request.setProvenance("{\"resourceType\":\"Provenance\",\"target\":[{\"reference\":\"urn:uuid:abc\"}]}");
        request.setCredentialType("PEM");
        request.setCredentialContent("-----BEGIN PRIVATE KEY-----\nFAKE\n-----END PRIVATE KEY-----");
        request.setCertificateChain("[\"CERT1\",\"CERT2\"]");
        request.setReferenceTimestamp(NOW);
        request.setStrategy("iat");
        request.setPolicyUri(validPolicy());
        return request;
    }

    private String validPolicy() {
        return "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2";
    }

    private void assertHasCode(List<ValidationError> errors, String expectedCode) {
        assertTrue(errors.stream().anyMatch(error -> error.getCode().equals(expectedCode)));
    }
}
