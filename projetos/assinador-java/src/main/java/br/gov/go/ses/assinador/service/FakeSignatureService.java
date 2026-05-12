package br.gov.go.ses.assinador.service;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.model.SignRequest;
import br.gov.go.ses.assinador.model.ValidateRequest;

import java.nio.charset.StandardCharsets;
import java.util.Base64;

public class FakeSignatureService implements SignatureService {

    private static final String FAKE_PAYLOAD =
            "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";

    private static final String FAKE_SIGNATURE =
            "FAKE_SIGNATURE_VALUE_BASE64URL_PLACEHOLDER_SPRINT2";

    @Override
    public AssinadorResponse sign(SignRequest request) {
        String jws = buildFakeJws(request.getReferenceTimestamp());
        String signatureData = Base64.getEncoder()
                .encodeToString(jws.getBytes(StandardCharsets.UTF_8));
        return AssinadorResponse.ok(signatureData);
    }

    @Override
    public AssinadorResponse validate(ValidateRequest request) {
        String data = request.getSignatureData();

        if (data == null || data.isBlank()) {
            return AssinadorResponse.error(
                    "FORMAT.JWS-MALFORMED",
                    "O campo signatureData esta ausente ou vazio."
            );
        }

        try {
            String json = new String(Base64.getDecoder().decode(data), StandardCharsets.UTF_8);
            if (!json.contains("\"payload\"") || !json.contains("\"signatures\"")) {
                return AssinadorResponse.error(
                        "FORMAT.JWS-MALFORMED",
                        "Estrutura JWS invalida: campos obrigatorios 'payload' ou 'signatures' ausentes."
                );
            }
        } catch (IllegalArgumentException error) {
            return AssinadorResponse.error(
                    "FORMAT.BASE64-INVALID",
                    "O campo signatureData nao e um base64 valido."
            );
        }

        return AssinadorResponse.ok("VALIDATION.SUCCESS");
    }

    private String buildFakeJws(Long referenceTimestamp) {
        long issuedAt = referenceTimestamp != null
                ? referenceTimestamp
                : currentUnixTimestamp();

        String protectedHeaderJson = """
                {"alg":"RS256","x5c":["FAKE_CERT_BASE64"],"iat":%d,"sigPId":{"id":"https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2"}}
                """.formatted(issuedAt).strip();

        String protectedHeader = Base64.getUrlEncoder().withoutPadding()
                .encodeToString(protectedHeaderJson.getBytes(StandardCharsets.UTF_8));

        return """
                {"payload":"%s","signatures":[{"protected":"%s","header":{"rRefs":{"ocspRefs":[],"crlRefs":[]}},"signature":"%s"}]}
                """.formatted(FAKE_PAYLOAD, protectedHeader, FAKE_SIGNATURE).strip();
    }

    private long currentUnixTimestamp() {
        return System.currentTimeMillis() / 1000L;
    }
}
