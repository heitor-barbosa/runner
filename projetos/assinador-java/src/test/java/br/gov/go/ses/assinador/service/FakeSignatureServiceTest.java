package br.gov.go.ses.assinador.service;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.model.SignRequest;
import br.gov.go.ses.assinador.model.ValidateRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.util.Base64;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FakeSignatureServiceTest {

    private FakeSignatureService service;

    @BeforeEach
    void setUp() {
        service = new FakeSignatureService();
    }

    @Test
    void signShouldReturnBase64SignatureData() {
        AssinadorResponse response = service.sign(buildValidSignRequest());

        assertTrue(response.isSuccess());
        assertNotNull(response.getData());
        assertFalse(response.getData().isBlank());
    }

    @Test
    void signShouldReturnJwsLikePayload() {
        AssinadorResponse response = service.sign(buildValidSignRequest());

        byte[] decoded = Base64.getDecoder().decode(response.getData());
        String json = new String(decoded, StandardCharsets.UTF_8);
        assertTrue(json.contains("\"payload\""));
        assertTrue(json.contains("\"signatures\""));
        assertTrue(json.contains("\"protected\""));
        assertTrue(json.contains("\"signature\""));
    }

    @Test
    void validateShouldAcceptGeneratedSignature() {
        AssinadorResponse signResponse = service.sign(buildValidSignRequest());

        ValidateRequest request = new ValidateRequest();
        request.setSignatureData(signResponse.getData());
        request.setReferenceTimestamp(1_751_328_001L);
        request.setPolicyUri(validPolicy());

        AssinadorResponse response = service.validate(request);
        assertTrue(response.isSuccess());
        assertEquals("VALIDATION.SUCCESS", response.getData());
    }

    @Test
    void validateShouldRejectMalformedBase64() {
        ValidateRequest request = new ValidateRequest();
        request.setSignatureData("!!!INVALIDO!!!");
        request.setReferenceTimestamp(1_751_328_001L);
        request.setPolicyUri(validPolicy());

        AssinadorResponse response = service.validate(request);
        assertFalse(response.isSuccess());
        assertEquals("FORMAT.BASE64-INVALID", response.getErrorCode());
    }

    @Test
    void signShouldReturnPkcs11ErrorWhenDeviceIsUnavailable() {
        FakeSignatureService pkcs11Service = new FakeSignatureService(new Pkcs11ProviderLoader() {
            @Override
            public java.security.Provider load(String configPath) {
                throw new IllegalStateException("simulated unavailable device");
            }
        });
        SignRequest request = buildValidSignRequest();
        request.setCredentialType("TOKEN");
        request.setPkcs11ConfigPath("pkcs11.cfg");
        request.setCredentialAlias("assinatura");

        AssinadorResponse response = pkcs11Service.sign(request);

        assertFalse(response.isSuccess());
        assertEquals("PKCS11.DEVICE-UNAVAILABLE", response.getErrorCode());
    }

    private SignRequest buildValidSignRequest() {
        SignRequest request = new SignRequest();
        request.setBundle("{\"resourceType\":\"Bundle\",\"entry\":[{}]}");
        request.setProvenance("{\"resourceType\":\"Provenance\",\"target\":[{\"reference\":\"urn:uuid:abc\"}]}");
        request.setCredentialType("PEM");
        request.setCredentialContent("-----BEGIN PRIVATE KEY-----\nFAKE\n-----END PRIVATE KEY-----");
        request.setCertificateChain("[\"CERT1_BASE64\",\"CERT2_BASE64\"]");
        request.setReferenceTimestamp(System.currentTimeMillis() / 1000L);
        request.setStrategy("iat");
        request.setPolicyUri(validPolicy());
        return request;
    }

    private String validPolicy() {
        return "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2";
    }
}
