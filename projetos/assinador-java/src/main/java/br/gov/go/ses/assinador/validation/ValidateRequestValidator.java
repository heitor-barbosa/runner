package br.gov.go.ses.assinador.validation;

import br.gov.go.ses.assinador.model.ValidateRequest;

import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Base64;
import java.util.List;
import java.util.regex.Pattern;

public class ValidateRequestValidator {

    private static final String POLICY_URI_PREFIX =
            "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|";
    private static final Pattern SEMVER = Pattern.compile("^\\d+\\.\\d+\\.\\d+$");
    private static final long TIMESTAMP_MIN = 1_751_328_000L;
    private static final long TIMESTAMP_MAX = 4_102_444_800L;

    public List<ValidationError> validate(ValidateRequest request) {
        List<ValidationError> errors = new ArrayList<>();

        if (request == null) {
            errors.add(new ValidationError("CONFIG.MISSING-PARAMETER", "Request nao pode ser nulo."));
            return errors;
        }

        validateSignatureData(request.getSignatureData(), errors);
        validateReferenceTimestamp(request.getReferenceTimestamp(), errors);
        validatePolicyUri(request.getPolicyUri(), errors);

        return errors;
    }

    private void validateSignatureData(String signatureData, List<ValidationError> errors) {
        if (signatureData == null || signatureData.isBlank()) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro signatureData e obrigatorio (valor de Signature.data em base64)."
            ));
            return;
        }

        try {
            byte[] decoded = Base64.getDecoder().decode(signatureData.strip());
            String json = new String(decoded, StandardCharsets.UTF_8);
            if (!json.contains("\"payload\"")) {
                errors.add(new ValidationError(
                        "FORMAT.JWS-MALFORMED",
                        "Estrutura JWS invalida: campo obrigatorio 'payload' ausente."
                ));
            }
            if (!json.contains("\"signatures\"")) {
                errors.add(new ValidationError(
                        "FORMAT.JWS-MALFORMED",
                        "Estrutura JWS invalida: campo obrigatorio 'signatures' ausente."
                ));
            }
        } catch (IllegalArgumentException error) {
            errors.add(new ValidationError(
                    "FORMAT.BASE64-INVALID",
                    "O campo signatureData nao e um base64 valido: " + error.getMessage()
            ));
        }
    }

    private void validateReferenceTimestamp(Long timestamp, List<ValidationError> errors) {
        if (timestamp == null) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro referenceTimestamp e obrigatorio."
            ));
            return;
        }
        if (timestamp < TIMESTAMP_MIN || timestamp > TIMESTAMP_MAX) {
            errors.add(new ValidationError(
                    "CONFIG.TIMESTAMP-OUT-OF-RANGE",
                    "O timestamp de referencia deve estar no intervalo [%d, %d]. Recebido: %d."
                            .formatted(TIMESTAMP_MIN, TIMESTAMP_MAX, timestamp)
            ));
        }
    }

    private void validatePolicyUri(String policyUri, List<ValidationError> errors) {
        if (policyUri == null || policyUri.isBlank()) {
            errors.add(new ValidationError("POLICY.MISSING", "O parametro policyUri e obrigatorio."));
            return;
        }
        if (!policyUri.startsWith(POLICY_URI_PREFIX)) {
            errors.add(new ValidationError(
                    "POLICY.URI-INVALID",
                    "A URI da politica deve iniciar com: " + POLICY_URI_PREFIX
            ));
            return;
        }
        String version = policyUri.substring(POLICY_URI_PREFIX.length());
        if (!SEMVER.matcher(version).matches()) {
            errors.add(new ValidationError(
                    "POLICY.URI-INVALID",
                    "A versao da politica deve seguir SemVer (ex: 0.1.2). Recebido: " + version
            ));
        }
    }
}
