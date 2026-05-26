package br.gov.go.ses.assinador.http;

import br.gov.go.ses.assinador.service.FakeSignatureService;
import br.gov.go.ses.assinador.service.Pkcs11ProviderLoader;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.Provider;
import java.time.Duration;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class Pkcs11IntegrationTest {

    private static final String POLICY_URI =
            "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2";

    private final ObjectMapper mapper = new ObjectMapper();
    private final HttpClient client = HttpClient.newHttpClient();
    private final SimulatedPkcs11ProviderLoader pkcs11Simulator = new SimulatedPkcs11ProviderLoader();

    private AssinadorHttpServer server;
    private String baseUrl;

    @BeforeEach
    void setUp() throws Exception {
        server = AssinadorHttpServer.create(0, new FakeSignatureService(pkcs11Simulator));
        server.start();
        baseUrl = "http://localhost:" + server.getPort();
    }

    @AfterEach
    void tearDown() {
        server.stop();
    }

    @Test
    void signShouldUseSimulatedPkcs11ProviderWhenCredentialTypeIsToken() throws Exception {
        HttpResponse<String> response = post("/sign", tokenSignRequest());
        JsonNode body = mapper.readTree(response.body());

        assertEquals(200, response.statusCode());
        assertTrue(body.get("success").asBoolean());
        assertNotNull(body.get("data").asText());
        assertTrue(pkcs11Simulator.wasLoaded());
        assertEquals("memory://softhsm2-simulator.cfg", pkcs11Simulator.getLoadedConfigPath());
    }

    private HttpResponse<String> post(String path, Object payload) throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .timeout(Duration.ofSeconds(5))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(mapper.writeValueAsString(payload)))
                .build();
        return client.send(request, HttpResponse.BodyHandlers.ofString());
    }

    private Map<String, Object> tokenSignRequest() {
        return Map.ofEntries(
                Map.entry("bundle", "{\"resourceType\":\"Bundle\",\"entry\":[{}]}"),
                Map.entry("provenance", "{\"resourceType\":\"Provenance\",\"target\":[{\"reference\":\"urn:uuid:abc\"}]}"),
                Map.entry("credentialType", "TOKEN"),
                Map.entry("credentialContent", "simulated-token"),
                Map.entry("credentialAlias", "assinatura"),
                Map.entry("pkcs11ConfigPath", "memory://softhsm2-simulator.cfg"),
                Map.entry("tokenLabel", "token-a"),
                Map.entry("certificateChain", "[\"CERT1_BASE64\",\"CERT2_BASE64\"]"),
                Map.entry("referenceTimestamp", currentUnixTimestamp()),
                Map.entry("strategy", "iat"),
                Map.entry("policyUri", POLICY_URI)
        );
    }

    private long currentUnixTimestamp() {
        return System.currentTimeMillis() / 1000L;
    }

    private static class SimulatedPkcs11ProviderLoader extends Pkcs11ProviderLoader {
        private String loadedConfigPath;

        @Override
        public Provider load(String configPath) {
            loadedConfigPath = configPath;
            return new Provider("SoftHSM2Simulator", "1.0", "PKCS#11 simulator for integration tests") {
            };
        }

        boolean wasLoaded() {
            return loadedConfigPath != null;
        }

        String getLoadedConfigPath() {
            return loadedConfigPath;
        }
    }
}
