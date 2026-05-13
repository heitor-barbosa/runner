package br.gov.go.ses.assinador.http;

import br.gov.go.ses.assinador.service.FakeSignatureService;
import br.gov.go.ses.assinador.service.SignatureService;
import br.gov.go.ses.assinador.validation.SignRequestValidator;
import br.gov.go.ses.assinador.validation.ValidateRequestValidator;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class AssinadorHttpServer implements AutoCloseable {

    public static final int DEFAULT_PORT = 8080;

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
