package br.gov.go.ses.assinador.validation;

import br.gov.go.ses.assinador.model.SignRequest;

import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;

public class SignRequestValidator {

    private static final String POLICY_URI_PREFIX =
            "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|";
    private static final Pattern SEMVER = Pattern.compile("^\\d+\\.\\d+\\.\\d+$");
    private static final long TIMESTAMP_MIN = 1_751_328_000L;
    private static final long TIMESTAMP_MAX = 4_102_444_800L;
    private static final long TOLERANCE_SECONDS = 300L;
    private static final List<String> VALID_CREDENTIAL_TYPES =
            List.of("PEM", "PKCS12", "SMARTCARD", "TOKEN");

    public List<ValidationError> validate(SignRequest request) {
        List<ValidationError> errors = new ArrayList<>();

        if (request == null) {
            errors.add(new ValidationError("CONFIG.MISSING-PARAMETER", "Request nao pode ser nulo."));
            return errors;
        }

        validatePolicyUri(request.getPolicyUri(), errors);
        validateReferenceTimestamp(request.getReferenceTimestamp(), errors);
        validateStrategy(request.getStrategy(), errors);
        validateBundle(request.getBundle(), errors);
        validateProvenance(request.getProvenance(), errors);
        validateCredential(request, errors);
        validateCertificateChain(request.getCertificateChain(), errors);

        return errors;
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
        long separators = policyUri.chars().filter(value -> value == '|').count();
        if (separators != 1) {
            errors.add(new ValidationError(
                    "POLICY.URI-INVALID",
                    "A URI da politica deve conter exatamente um caractere '|'."
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
            return;
        }
        long now = System.currentTimeMillis() / 1000L;
        long difference = Math.abs(timestamp - now);
        if (difference > TOLERANCE_SECONDS) {
            errors.add(new ValidationError(
                    "TIMESTAMP.OUT-OF-TOLERANCE-WINDOW",
                    "O timestamp de referencia difere mais de %d segundos do relogio do servidor. Diferenca: %d segundos."
                            .formatted(TOLERANCE_SECONDS, difference)
            ));
        }
    }

    private void validateStrategy(String strategy, List<ValidationError> errors) {
        if (strategy == null || strategy.isBlank()) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro strategy e obrigatorio (valores aceitos: iat, tsa)."
            ));
            return;
        }
        if (!strategy.equals("iat") && !strategy.equals("tsa")) {
            errors.add(new ValidationError(
                    "CONFIG.INVALID-STRATEGY",
                    "Estrategia invalida: '" + strategy + "'. Valores aceitos: iat, tsa."
            ));
        }
    }

    private void validateBundle(String bundle, List<ValidationError> errors) {
        if (bundle == null || bundle.isBlank()) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro bundle e obrigatorio (instancia Bundle FHIR R4 em JSON)."
            ));
            return;
        }
        String trimmed = bundle.strip();
        if (!trimmed.startsWith("{")) {
            errors.add(new ValidationError("FORMAT.BUNDLE-MALFORMED", "O bundle deve ser um objeto JSON valido."));
            return;
        }
        if (!trimmed.contains("\"resourceType\"")) {
            errors.add(new ValidationError("FORMAT.BUNDLE-MALFORMED", "O bundle deve conter o campo 'resourceType'."));
            return;
        }
        if (!trimmed.contains("\"Bundle\"")) {
            errors.add(new ValidationError(
                    "FORMAT.BUNDLE-MALFORMED",
                    "O campo 'resourceType' do bundle deve ser 'Bundle'."
            ));
        }
        if (!trimmed.contains("\"entry\"")) {
            errors.add(new ValidationError("FORMAT.BUNDLE-EMPTY", "O bundle deve conter ao menos uma entrada em 'entry'."));
        }
    }

    private void validateProvenance(String provenance, List<ValidationError> errors) {
        if (provenance == null || provenance.isBlank()) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro provenance e obrigatorio (instancia Provenance FHIR R4 em JSON)."
            ));
            return;
        }
        String trimmed = provenance.strip();
        if (!trimmed.startsWith("{")) {
            errors.add(new ValidationError("FORMAT.PROVENANCE-INVALID", "O provenance deve ser um objeto JSON valido."));
            return;
        }
        if (!trimmed.contains("\"Provenance\"")) {
            errors.add(new ValidationError(
                    "FORMAT.PROVENANCE-INVALID",
                    "O campo 'resourceType' do provenance deve ser 'Provenance'."
            ));
        }
        if (!trimmed.contains("\"target\"")) {
            errors.add(new ValidationError(
                    "FORMAT.PROVENANCE-INVALID",
                    "O provenance deve conter ao menos uma referencia em 'target'."
            ));
        }
    }

    private void validateCredential(SignRequest request, List<ValidationError> errors) {
        String type = request.getCredentialType();
        if (type == null || type.isBlank()) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro credentialType e obrigatorio. Valores aceitos: PEM, PKCS12, SMARTCARD, TOKEN."
            ));
            return;
        }
        if (!VALID_CREDENTIAL_TYPES.contains(type.toUpperCase())) {
            errors.add(new ValidationError(
                    "CONFIG.INVALID-PARAMETER",
                    "Tipo de credencial invalido: '" + type + "'. Valores aceitos: PEM, PKCS12, SMARTCARD, TOKEN."
            ));
            return;
        }

        String content = request.getCredentialContent();
        if (content == null || content.isBlank()) {
            errors.add(new ValidationError("CONFIG.MISSING-PARAMETER", "O parametro credentialContent e obrigatorio."));
        }

        String normalizedType = type.toUpperCase();
        if ("PKCS12".equals(normalizedType)) {
            if (request.getCredentialPassword() == null || request.getCredentialPassword().isBlank()) {
                errors.add(new ValidationError(
                        "CONFIG.MISSING-PARAMETER",
                        "A senha credentialPassword e obrigatoria para credencial PKCS12."
                ));
            }
            if (request.getCredentialAlias() == null || request.getCredentialAlias().isBlank()) {
                errors.add(new ValidationError(
                        "CONFIG.MISSING-PARAMETER",
                        "O alias credentialAlias e obrigatorio para credencial PKCS12."
                ));
            } else if (request.getCredentialAlias().length() > 64) {
                errors.add(new ValidationError(
                        "MIDDLEWARE.TOKEN-LABEL-INVALID",
                        "O alias credentialAlias deve ter no maximo 64 caracteres."
                ));
            }
        }

        if ("TOKEN".equals(normalizedType) || "SMARTCARD".equals(normalizedType)) {
            if (request.getPkcs11ConfigPath() == null || request.getPkcs11ConfigPath().isBlank()) {
                errors.add(new ValidationError(
                        "PKCS11.CONFIG-MISSING",
                        "O parametro pkcs11ConfigPath e obrigatorio para credenciais TOKEN ou SMARTCARD."
                ));
            }
            if (request.getCredentialAlias() == null || request.getCredentialAlias().isBlank()) {
                errors.add(new ValidationError(
                        "PKCS11.ALIAS-MISSING",
                        "O parametro credentialAlias e obrigatorio para selecionar a chave no dispositivo criptografico."
                ));
            } else if (request.getCredentialAlias().length() > 64) {
                errors.add(new ValidationError(
                        "MIDDLEWARE.TOKEN-LABEL-INVALID",
                        "O alias credentialAlias deve ter no maximo 64 caracteres."
                ));
            }
            if (request.getTokenLabel() != null && request.getTokenLabel().length() > 64) {
                errors.add(new ValidationError(
                        "MIDDLEWARE.TOKEN-LABEL-INVALID",
                        "O rotulo tokenLabel deve ter no maximo 64 caracteres."
                ));
            }
        }
    }

    private void validateCertificateChain(String certificateChain, List<ValidationError> errors) {
        if (certificateChain == null || certificateChain.isBlank()) {
            errors.add(new ValidationError(
                    "CONFIG.MISSING-PARAMETER",
                    "O parametro certificateChain e obrigatorio (array JSON de certificados X.509 DER em base64)."
            ));
            return;
        }
        String trimmed = certificateChain.strip();
        if (!trimmed.startsWith("[")) {
            errors.add(new ValidationError("CERT.INVALID-FORMAT", "certificateChain deve ser um array JSON."));
            return;
        }
        long certificateCount = trimmed.chars().filter(value -> value == '"').count() / 2;
        if (certificateCount < 2) {
            errors.add(new ValidationError(
                    "CERT.CHAIN-INCOMPLETE",
                    "A cadeia de certificados deve conter pelo menos 2 certificados."
            ));
        }
    }
}
