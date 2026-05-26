package br.gov.go.ses.assinador.service;

import java.nio.file.Files;
import java.nio.file.Path;
import java.security.Provider;
import java.security.Security;

public class Pkcs11ProviderLoader {

    public Provider load(String configPath) {
        if (configPath == null || configPath.isBlank()) {
            throw new IllegalArgumentException("Caminho da configuracao PKCS#11 nao informado.");
        }

        Path path = Path.of(configPath);
        if (!Files.isRegularFile(path)) {
            throw new IllegalArgumentException("Arquivo de configuracao PKCS#11 nao encontrado: " + configPath);
        }

        Provider baseProvider = Security.getProvider("SunPKCS11");
        if (baseProvider == null) {
            throw new IllegalStateException("Provider SunPKCS11 nao esta disponivel neste JDK.");
        }

        Provider configuredProvider = baseProvider.configure(path.toAbsolutePath().toString());
        Security.addProvider(configuredProvider);
        return configuredProvider;
    }
}
