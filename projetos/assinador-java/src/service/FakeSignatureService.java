package br.gov.go.ses.assinador.service;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.model.SignRequest;
import br.gov.go.ses.assinador.model.ValidateRequest;
import java.nio.charset.StandardCharsets;
import java.util.Base64;

public class FakeSignatureService implements SignatureService {
  
    private static final String FAKE_PROTECTED_HEADER = Base64.getUrlEncoder().withoutPadding()
            .encodeToString("""
                    {"alg":"RS256","x5c":["FAKE_CERT_BASE64"],"sigPId":{"id":"https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2"}}
                    """.strip().getBytes(StandardCharsets.UTF_8));

    /** Payload fixo: hash SHA-256 fictício em base64Url de 43 caracteres. */
    private static final String FAKE_PAYLOAD =
            "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";

    /** Assinatura criptográfica fictícia em base64Url. */
    private static final String FAKE_SIGNATURE =
            "FAKE_SIGNATURE_VALUE_BASE64URL_PLACEHOLDER_SPRINT2";

    @Override
    public AssinadorResponse sign(SignRequest request) {
        // Monta JWS JSON Serialization simulado conforme RFC 7515 seção 3.2
        String jws = buildFakeJws(request.getReferenceTimestamp());

        // Codifica em base64 padrão (não base64Url) conforme Signature.data FHIR
        String signatureData = Base64.getEncoder()
                .encodeToString(jws.getBytes(StandardCharsets.UTF_8));

        return AssinadorResponse.ok(signatureData);
    }

    @Override
    public AssinadorResponse validate(ValidateRequest request) {
        // Simulação simples: qualquer signatureData não nulo é considerado válido.
        // A Sprint 3 introduzirá validação criptográfica real via endpoints HTTP.
        String data = request.getSignatureData();

        if (data == null || data.isBlank()) {
            return AssinadorResponse.error(
                    "FORMAT.JWS-MALFORMED",
                    "O campo signatureData está ausente ou vazio."
            );
        }
      
        // Tenta decodificar o base64 para verificar integridade básica do envelope
        try {
            String json = new String(Base64.getDecoder().decode(data), StandardCharsets.UTF_8);
            if (!json.contains("\"payload\"") || !json.contains("\"signatures\"")) {
                return AssinadorResponse.error(
                        "FORMAT.JWS-MALFORMED",
                        "Estrutura JWS inválida: campos obrigatórios 'payload' ou 'signatures' ausentes."
                );
            }
        } catch (IllegalArgumentException e) {
            return AssinadorResponse.error(
                    "FORMAT.BASE64-INVALID",
                    "O campo signatureData não é um base64 válido."
            );
        }

        return AssinadorResponse.ok("VALIDATION.SUCCESS");
    }

    private String buildFakeJws(Long referenceTimestamp) {
        long iat = (referenceTimestamp != null) ? referenceTimestamp : currentUnixTimestamp();

        // Protected header com iat embutido
        String protectedWithIat = Base64.getUrlEncoder().withoutPadding()
                .encodeToString(("""
                        {"alg":"RS256","x5c":["FAKE_CERT_BASE64"],"iat":%d,"sigPId":{"id":"https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2"}}
                        """.formatted(iat)).strip().getBytes(StandardCharsets.UTF_8));

        return """
                {"payload":"%s","signatures":[{"protected":"%s","header":{"rRefs":{"ocspRefs":[],"crlRefs":[]}},"signature":"%s"}]}
                """.formatted(FAKE_PAYLOAD, protectedWithIat, FAKE_SIGNATURE).strip();
    }

    private long currentUnixTimestamp() {
        return System.currentTimeMillis() / 1000L;
    }
}
