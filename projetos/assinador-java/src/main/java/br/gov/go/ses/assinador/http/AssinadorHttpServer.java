package br.gov.go.ses.assinador.http;

import br.gov.go.ses.assinador.model.AssinadorResponse;
import br.gov.go.ses.assinador.service.FakeSignatureService;
import br.gov.go.ses.assinador.service.SignatureService;
import br.gov.go.ses.assinador.validation.SignRequestValidator;
import br.gov.go.ses.assinador.validation.ValidateRequestValidator;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class AssinadorHttpServer implements AutoCloseable {

    public static final int DEFAULT_PORT = 8080;
    private static final String CONTENT_TYPE_JSON = "application/json; charset=utf-8";

    private final HttpServer server;
    private final ExecutorService executor;

    private AssinadorHttpServer(HttpServer server, ExecutorService executor) {
        this.server = server;
        this.executor = executor;
    }

    public static AssinadorHttpServer create(int port) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        ExecutorService executor = Executors.newCachedThreadPool();
        SignatureService service = new FakeSignatureService();
        SignatureController controller = new SignatureController(
                service,
                new SignRequestValidator(),
                new ValidateRequestValidator()
        );

        server.createContext("/sign", controller::handleSign);
        server.createContext("/validate", controller::handleValidate);
        server.createContext("/health", exchange -> {
            if (!"GET".equalsIgnoreCase(exchange.getRequestMethod())) {
                exchange.getResponseHeaders().add("Allow", "GET");
                byte[] error = new ObjectMapper().writeValueAsBytes(AssinadorResponse.error(
                        "HTTP.METHOD-NOT-ALLOWED",
                        "Metodo nao permitido. Use GET."
                ));
                exchange.getResponseHeaders().set("Content-Type", CONTENT_TYPE_JSON);
                exchange.sendResponseHeaders(405, error.length);
                try (OutputStream output = exchange.getResponseBody()) {
                    output.write(error);
                }
                return;
            }

            byte[] payload = new ObjectMapper().writeValueAsBytes(AssinadorResponse.ok("HEALTH.OK"));
            exchange.getResponseHeaders().set("Content-Type", CONTENT_TYPE_JSON);
            exchange.sendResponseHeaders(200, payload.length);
            try (OutputStream output = exchange.getResponseBody()) {
                output.write(payload);
            }
        });
        server.setExecutor(executor);

        return new AssinadorHttpServer(server, executor);
    }

    public void start() {
        server.start();
    }

    public int getPort() {
        return server.getAddress().getPort();
    }

    public void stop() {
        server.stop(0);
        executor.shutdownNow();
    }

    @Override
    public void close() {
        stop();
    }
}
