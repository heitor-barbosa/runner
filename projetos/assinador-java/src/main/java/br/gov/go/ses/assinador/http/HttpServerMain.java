package br.gov.go.ses.assinador.http;

import java.util.concurrent.CountDownLatch;

public class HttpServerMain {

    public void startAndBlock(int port) {
        try {
            AssinadorHttpServer server = AssinadorHttpServer.create(port);
            Runtime.getRuntime().addShutdownHook(new Thread(server::stop));
            server.start();
            System.out.println("Assinador HTTP iniciado na porta " + server.getPort());
            new CountDownLatch(1).await();
        } catch (InterruptedException error) {
            Thread.currentThread().interrupt();
        } catch (Exception error) {
            System.err.println("Erro ao iniciar servidor HTTP do Assinador: " + error.getMessage());
            System.exit(1);
        }
    }
}
