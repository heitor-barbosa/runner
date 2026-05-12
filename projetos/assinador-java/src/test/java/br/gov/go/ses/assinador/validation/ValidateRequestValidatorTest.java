package br.gov.go.ses.assinador.validation;

import br.gov.go.ses.assinador.model.ValidateRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertTrue;

class ValidateRequestValidatorTest {

    private static final String VALID_SIGNATURE_DATA = Base64.getEncoder()
            .encodeToString("{\"payload\":\"ABC\",\"signatures\":[{\"protected\":\"HDR\",\"signature\":\"SIG\"}]}"
                    .getBytes(StandardCharsets.UTF_8));

    private ValidateRequestValidator validator;

    @BeforeEach
    void setUp() {
        validator = new ValidateRequestValidator();
    }

    @Test
    void validRequestShouldNotGenerateErrors() {
        assertTrue(validator.validate(buildValid()).isEmpty());
    }

    @Test
    void invalidBase64ShouldGenerateError() {
        ValidateRequest request = buildValid();
        request.setSignatureData("!!!INVALIDO!!!");
        assertHasCode(validator.validate(request), "FORMAT.BASE64-INVALID");
    }

    @Test
    void invalidTimestampShouldGenerateError() {
        ValidateRequest request = buildValid();
        request.setReferenceTimestamp(1_000_000L);
        assertHasCode(validator.validate(request), "CONFIG.TIMESTAMP-OUT-OF-RANGE");
    }

    @Test
    void missingPolicyShouldGenerateError() {
        ValidateRequest request = buildValid();
        request.setPolicyUri(null);
        assertHasCode(validator.validate(request), "POLICY.MISSING");
    }

    @Test
    void malformedJwsShouldGenerateError() {
        String payloadOnly = Base64.getEncoder()
                .encodeToString("{\"payload\":\"ABC\"}".getBytes(StandardCharsets.UTF_8));
        ValidateRequest request = buildValid();
        request.setSignatureData(payloadOnly);
        assertHasCode(validator.validate(request), "FORMAT.JWS-MALFORMED");
    }

    private ValidateRequest buildValid() {
        ValidateRequest request = new ValidateRequest();
        request.setSignatureData(VALID_SIGNATURE_DATA);
        request.setReferenceTimestamp(1_751_328_001L);
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
