package br.gov.go.ses.assinador.http;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.util.Map;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class SignatureControllerTest {

    private static final String POLICY_URI =
            "https://fhir.saude.go.gov.br/r4/seguranca/ImplementationGuide/br.go.ses.seguranca|0.1.2";

    private final ObjectMapper mapper = new ObjectMapper();
    private final HttpClient client = HttpClient.newHttpClient();

    private AssinadorHttpServer server;
    private String baseUrl;

    @BeforeEach
    void setUp() throws Exception {
        server = AssinadorHttpServer.create(0);
        server.start();
        baseUrl = "http://localhost:" + server.getPort();
    }

    @AfterEach
    void tearDown() {
        server.stop();
    }

    @Test
    void signShouldReturnSuccessResponse() throws Exception {
        HttpResponse<String> response = post("/sign", validSignRequest());
        JsonNode body = mapper.readTree(response.body());

        assertEquals(200, response.statusCode());
        assertTrue(body.get("success").asBoolean());
        assertNotNull(body.get("data").asText());
        assertFalse(body.get("data").asText().isBlank());
    }

    @Test
    void signShouldReturnValidationErrorForInvalidRequest() throws Exception {
        HttpResponse<String> response = post("/sign", Map.of());
        JsonNode body = mapper.readTree(response.body());

        assertEquals(400, response.statusCode());
        assertFalse(body.get("success").asBoolean());
        assertEquals("POLICY.MISSING", body.get("errorCode").asText());
    }

    @Test
    void validateShouldReturnSuccessResponseForGeneratedSignature() throws Exception {
        HttpResponse<String> signResponse = post("/sign", validSignRequest());
        String signatureData = mapper.readTree(signResponse.body()).get("data").asText();

        HttpResponse<String> validateResponse = post("/validate", Map.of(
                "signatureData", signatureData,
                "referenceTimestamp", currentUnixTimestamp(),
                "policyUri", POLICY_URI
        ));
        JsonNode body = mapper.readTree(validateResponse.body());

        assertEquals(200, validateResponse.statusCode());
        assertTrue(body.get("success").asBoolean());
        assertEquals("VALIDATION.SUCCESS", body.get("data").asText());
    }

    @Test
    void validateShouldReturnValidationErrorForInvalidRequest() throws Exception {
        HttpResponse<String> response = post("/validate", Map.of(
                "signatureData", "!!!INVALIDO!!!",
                "referenceTimestamp", currentUnixTimestamp(),
                "policyUri", POLICY_URI
        ));
        JsonNode body = mapper.readTree(response.body());

        assertEquals(400, response.statusCode());
        assertFalse(body.get("success").asBoolean());
        assertEquals("FORMAT.BASE64-INVALID", body.get("errorCode").asText());
    }

    @Test
    void signShouldRejectNonPostMethod() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/sign"))
                .timeout(Duration.ofSeconds(5))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        JsonNode body = mapper.readTree(response.body());

        assertEquals(405, response.statusCode());
        assertEquals("POST", response.headers().firstValue("Allow").orElse(""));
        assertFalse(body.get("success").asBoolean());
        assertEquals("HTTP.METHOD-NOT-ALLOWED", body.get("errorCode").asText());
    }

    @Test
    void healthShouldReturnSuccessResponse() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/health"))
                .timeout(Duration.ofSeconds(5))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        JsonNode body = mapper.readTree(response.body());

        assertEquals(200, response.statusCode());
        assertTrue(body.get("success").asBoolean());
        assertEquals("HEALTH.OK", body.get("data").asText());
    }

    @Test
    void shutdownShouldTerminateServerWithDelete() throws Exception {
        HttpResponse<String> response = delete("/shutdown");
        JsonNode body = mapper.readTree(response.body());

        assertEquals(200, response.statusCode());
        assertTrue(body.get("success").asBoolean());
        assertEquals("SHUTDOWN.OK", body.get("data").asText());

        Thread.sleep(100);

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/health"))
                .timeout(Duration.ofSeconds(5))
                .GET()
                .build();

        try {
            client.send(request, HttpResponse.BodyHandlers.ofString());
            assertTrue(false, "Server should not respond after shutdown");
        } catch (Exception e) {
            assertTrue(true, "Server is down as expected");
        }
    }

    @Test
    void shutdownShouldRejectNonDeleteMethod() throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/shutdown"))
                .timeout(Duration.ofSeconds(5))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        JsonNode body = mapper.readTree(response.body());

        assertEquals(405, response.statusCode());
        assertEquals("DELETE", response.headers().firstValue("Allow").orElse(""));
        assertFalse(body.get("success").asBoolean());
        assertEquals("HTTP.METHOD-NOT-ALLOWED", body.get("errorCode").asText());
    }

    @Test
    void idleTimeoutShouldRestartAfterEachRequest() throws Exception {
        server.stop();

        CountDownLatch stopped = new CountDownLatch(1);
        server = AssinadorHttpServer.createForTests(0, 250, 50, stopped::countDown);
        server.start();
        baseUrl = "http://localhost:" + server.getPort();

        for (int index = 0; index < 4; index++) {
            healthShouldReturnSuccessResponse();
            Thread.sleep(100);
        }

        assertFalse(stopped.await(50, TimeUnit.MILLISECONDS));
        assertTrue(stopped.await(2, TimeUnit.SECONDS));
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

    private HttpResponse<String> delete(String path) throws Exception {
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + path))
                .timeout(Duration.ofSeconds(5))
                .DELETE()
                .build();
        return client.send(request, HttpResponse.BodyHandlers.ofString());
    }

    private Map<String, Object> validSignRequest() {
        return Map.of(
                "bundle", "{\"resourceType\":\"Bundle\",\"entry\":[{}]}",
                "provenance", "{\"resourceType\":\"Provenance\",\"target\":[{\"reference\":\"urn:uuid:abc\"}]}",
                "credentialType", "PEM",
                "credentialContent", "-----BEGIN PRIVATE KEY-----\nFAKE\n-----END PRIVATE KEY-----",
                "certificateChain", "[\"CERT1_BASE64\",\"CERT2_BASE64\"]",
                "referenceTimestamp", currentUnixTimestamp(),
                "strategy", "iat",
                "policyUri", POLICY_URI
        );
    }

    private long currentUnixTimestamp() {
        return System.currentTimeMillis() / 1000L;
    }
}
