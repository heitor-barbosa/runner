package br.gov.go.ses.assinador.http;

import java.util.concurrent.CountDownLatch;

public class HttpServerMain {

    public void startAndBlock(int port) {
        startAndBlock(port, 0);
    }

    public void startAndBlock(int port, int timeoutMinutes) {
        CountDownLatch shutdownLatch = new CountDownLatch(1);
        try {
            AssinadorHttpServer server = AssinadorHttpServer.create(port, timeoutMinutes, shutdownLatch::countDown);
            Runtime.getRuntime().addShutdownHook(new Thread(server::stop));
            server.start();
            System.out.println("Assinador HTTP iniciado na porta " + server.getPort());
            if (timeoutMinutes > 0) {
                System.out.println("Timeout por inatividade configurado para " + timeoutMinutes + " minuto(s)");
            }
            shutdownLatch.await();
        } catch (InterruptedException error) {
            Thread.currentThread().interrupt();
        } catch (Exception error) {
            System.err.println("Erro ao iniciar servidor HTTP do Assinador: " + error.getMessage());
            System.exit(1);
        }
    }
}
