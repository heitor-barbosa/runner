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
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;

public class AssinadorHttpServer implements AutoCloseable {

    public static final int DEFAULT_PORT = 8080;
    private static final String CONTENT_TYPE_JSON = "application/json; charset=utf-8";

    private final HttpServer server;
    private final ExecutorService executor;
    private final ScheduledExecutorService timeoutExecutor;
    private final AtomicLong lastInteractionMillis;
    private final AtomicBoolean stopped;
    private final long timeoutMillis;
    private final Runnable onStop;

    private AssinadorHttpServer(
            HttpServer server,
            ExecutorService executor,
            ScheduledExecutorService timeoutExecutor,
            long timeoutMillis,
            Runnable onStop
    ) {
        this.server = server;
        this.executor = executor;
        this.timeoutExecutor = timeoutExecutor;
        this.timeoutMillis = timeoutMillis;
        this.onStop = onStop;
        this.lastInteractionMillis = new AtomicLong(System.currentTimeMillis());
        this.stopped = new AtomicBoolean(false);
    }

    public static AssinadorHttpServer create(int port) throws IOException {
        return create(port, 0, null);
    }

    public static AssinadorHttpServer create(int port, int timeoutMinutes, Runnable onStop) throws IOException {
        HttpServer server = HttpServer.create(new InetSocketAddress(port), 0);
        ExecutorService executor = Executors.newCachedThreadPool();
        ScheduledExecutorService timeoutExecutor = Executors.newSingleThreadScheduledExecutor();
        long timeoutMillis = timeoutMinutes > 0 ? TimeUnit.MINUTES.toMillis(timeoutMinutes) : 0L;
        AssinadorHttpServer assinadorServer = new AssinadorHttpServer(
                server,
                executor,
                timeoutExecutor,
                timeoutMillis,
                onStop
        );
        SignatureService service = new FakeSignatureService();
        SignatureController controller = new SignatureController(
                service,
                new SignRequestValidator(),
                new ValidateRequestValidator()
        );

        server.createContext("/sign", exchange -> {
            assinadorServer.touch();
            controller.handleSign(exchange);
        });
        server.createContext("/validate", exchange -> {
            assinadorServer.touch();
            controller.handleValidate(exchange);
        });
        server.createContext("/health", exchange -> {
            assinadorServer.touch();
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

        return assinadorServer;
    }

    public void start() {
        server.start();
        if (timeoutMillis > 0) {
            timeoutExecutor.scheduleAtFixedRate(this::stopIfIdle, 1, 1, TimeUnit.MINUTES);
        }
    }

    public int getPort() {
        return server.getAddress().getPort();
    }

    public void stop() {
        if (!stopped.compareAndSet(false, true)) {
            return;
        }
        server.stop(0);
        executor.shutdownNow();
        timeoutExecutor.shutdownNow();
        if (onStop != null) {
            onStop.run();
        }
    }

    @Override
    public void close() {
        stop();
    }

    private void touch() {
        lastInteractionMillis.set(System.currentTimeMillis());
    }

    private void stopIfIdle() {
        long idleMillis = System.currentTimeMillis() - lastInteractionMillis.get();
        if (idleMillis >= timeoutMillis) {
            stop();
        }
    }
}
